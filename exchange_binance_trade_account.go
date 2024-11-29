package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

type BinanceTradeAccount struct {
	ExchangeBase

	bnConverter BinanceEnumConverter
	apiKey      string
	secretKey   string
}

func (b BinanceTradeAccount) GetAccountMode() (AccountMode, error) {
	res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureMultiAssetsMarginGet().Do()
	if err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	return b.bnConverter.FromBNAccountMode(res.MultiAssetsMargin), nil
}

func (b BinanceTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		return MARGIN_MODE_CROSSED, nil
	case BN_AC_FUTURE:
		res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
		if err != nil {
			return MARGIN_MODE_UNKNOWN, err
		}
		for _, p := range res.Positions {
			if p.Symbol == symbol && p.PositionSide == b.bnConverter.ToBNPositionSide(positionSide) {
				return b.bnConverter.FromBNMarginMode(p.Isolated), nil
			}
		}
	case BN_AC_SWAP:
		res, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapAccount().Do()
		if err != nil {
			return MARGIN_MODE_UNKNOWN, err
		}
		for _, p := range res.Positions {
			if p.Symbol == symbol && p.PositionSide == b.bnConverter.ToBNPositionSide(positionSide) {
				return b.bnConverter.FromBNMarginMode(p.Isolated), nil
			}
		}
	default:
		return MARGIN_MODE_UNKNOWN, ErrorAccountType
	}
	return MARGIN_MODE_UNKNOWN, ErrorSymbolNotFound
}

func (b BinanceTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		return POSITION_MODE_ONEWAY, nil
	case BN_AC_FUTURE:
		res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFuturePositionSideDualGet().Do()
		if err != nil {
			return POSITION_MODE_UNKNOWN, err
		}
		return b.bnConverter.FromBNPositionMode(res.DualSidePosition), nil
	case BN_AC_SWAP:
		res, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapPositionSideDualGet().Do()
		if err != nil {
			return POSITION_MODE_UNKNOWN, err
		}
		return b.bnConverter.FromBNPositionMode(res.DualSidePosition), nil
	default:
		return POSITION_MODE_UNKNOWN, ErrorNotSupport
	}
}

func (b BinanceTradeAccount) GetLeverage(accountType, symbol string,
	marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	leverage := decimal.NewFromInt(0)

	if accountType == BN_AC_SPOT.String() {
		if marginMode == MARGIN_MODE_CROSSED {
			// tips: maxleverage only 3x， 5x，10x are supported
			res, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginTradeCoeff().Do()
			if err != nil {
				log.Error(err)
				return leverage, err
			}
			// 通过初始风险率反推杠杆倍数
			normalBar := res.NormalBar
			switch normalBar {
			case "1.5":
				leverage = decimal.NewFromInt(3)
			case "1.25":
				leverage = decimal.NewFromInt(5)
			case "2":
				leverage = decimal.NewFromInt(10)
			}
			return leverage, nil
		} else {
			// get spot isolated leverage is not supported （现货逐仓杠杆没有设置和查询接口）
			return leverage, ErrorNotSupport
		}
	}

	positionMode, err := b.GetPositionMode(accountType, symbol)
	if err != nil {
		return leverage, err
	}

	check := func(s string, i bool, ps string, l string) (decimal.Decimal, bool) {
		var result decimal.Decimal
		if s == symbol && i == b.bnConverter.ToBNMarginMode(marginMode) {
			//根据保证金模式，仓位模式，仓位方向获取指定杠杆
			switch positionMode {
			case POSITION_MODE_ONEWAY:
				if ps == b.bnConverter.ToBNPositionSide(POSITION_SIDE_BOTH) {
					result, _ = decimal.NewFromString(l)
					return result, true
				}
			case POSITION_MODE_HEDGE:
				if ps == b.bnConverter.ToBNPositionSide(positionSide) {
					result, _ = decimal.NewFromString(l)
					return result, true
				}
			default:
				return result, false
			}
		}
		return result, false
	}

	switch BinanceAccountType(accountType) {
	case BN_AC_FUTURE:
		res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
		if err != nil {
			return leverage, err
		}
		for _, p := range res.Positions {
			if r, ok := check(p.Symbol, p.Isolated, p.PositionSide, p.Leverage); ok {
				leverage = r
				break
			}
		}
	case BN_AC_SWAP:
		res, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapAccount().Do()
		if err != nil {
			return leverage, err
		}
		for _, p := range res.Positions {
			if r, ok := check(p.Symbol, p.Isolated, p.PositionSide, p.Leverage); ok {
				leverage = r
				break
			}
		}
	default:
		return leverage, ErrorNotSupport
	}
	if leverage.IsZero() {
		return leverage, ErrorPositionNotFound
	}
	return leverage, nil
}

func (b BinanceTradeAccount) SetAccountMode(mode AccountMode) error {
	nowAccountMode, err := b.GetAccountMode()
	if err != nil {
		return err
	}
	if nowAccountMode == mode {
		return nil
	}
	_, err = binance.NewFutureRestClient(b.apiKey, b.secretKey).
		NewFutureMultiAssetsMarginPost().
		MultiAssetsMargin(b.bnConverter.ToBNAccountMode(mode)).Do()
	if err != nil {
		return err
	}
	return nil
}

func (b BinanceTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	if accountType == BN_AC_SPOT.String() {
		return nil
	}

	positionMode, err := b.GetPositionMode(accountType, symbol)
	if err != nil {
		return err
	}
	positionSide := POSITION_SIDE_LONG
	if positionMode == POSITION_MODE_ONEWAY {
		positionSide = POSITION_SIDE_BOTH
	}
	nowMarginMode, err := b.GetMarginMode(accountType, symbol, positionSide)
	if err != nil {
		return err
	}

	if nowMarginMode == mode {
		return nil
	}

	switch BinanceAccountType(accountType) {
	case BN_AC_FUTURE:
		_, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFutureMarginType().Symbol(symbol).MarginType(b.bnConverter.ToBNMarginModeStr(mode)).Do()
		if err != nil {
			return err
		}
	case BN_AC_SWAP:
		_, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).
			NewSwapMarginType().Symbol(symbol).MarginType(b.bnConverter.ToBNMarginModeStr(mode)).Do()
		if err != nil {
			return err
		}
	default:
		return ErrorAccountType
	}
	return nil
}

func (b BinanceTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {

	nowPositionMode, err := b.GetPositionMode(accountType, symbol)
	if err != nil {
		return err
	}

	if nowPositionMode == mode {
		return nil
	}

	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		return ErrorNotSupport
	case BN_AC_FUTURE:
		_, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFuturePositionSideDualPost().DualSidePosition(b.bnConverter.ToBNPositionMode(mode)).Do()
		if err != nil {
			return err
		}
	case BN_AC_SWAP:
		_, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).
			NewSwapPositionSideDualPost().DualSidePosition(b.bnConverter.ToBNPositionMode(mode)).Do()
		if err != nil {
			return err
		}
	default:
		return ErrorAccountType
	}
	return nil
}

func (b BinanceTradeAccount) SetLeverage(accountType, symbol string,
	marginMode MarginMode, positionSide PositionSide,
	leverage decimal.Decimal) error {
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		if marginMode == MARGIN_MODE_CROSSED {
			// tips: maxleverage only 3x， 5x， 10x are supported
			_, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).
				NewSpotMarginMaxLeverage().MaxLeverage(leverage.IntPart()).Do()
			if err != nil {
				return err
			}
		} else {
			// set spot isolated leverage is not supported
			return ErrorNotSupport
		}
	case BN_AC_FUTURE:
		_, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFutureLeverage().Symbol(symbol).Leverage(leverage.IntPart()).Do()
		if err != nil {
			return err
		}
	case BN_AC_SWAP:
		_, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).
			NewSwapLeverage().Symbol(symbol).Leverage(leverage.IntPart()).Do()
		if err != nil {
			return err
		}
	default:
		return ErrorAccountType
	}
	return nil
}

func (b BinanceTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	var feeRate FeeRate
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		res, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAccount().Do()
		if err != nil {
			return nil, err
		}
		feeRate.Maker, _ = decimal.NewFromString(res.CommissionRates.Maker)
		feeRate.Taker, _ = decimal.NewFromString(res.CommissionRates.Taker)
	case BN_AC_FUTURE:
		res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFutureCommissionRate().Symbol(symbol).Do()
		if err != nil {
			return nil, err
		}
		feeRate.Maker, _ = decimal.NewFromString(res.MakerCommissionRate)
		feeRate.Taker, _ = decimal.NewFromString(res.TakerCommissionRate)
	case BN_AC_SWAP:
		res, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).
			NewSwapCommissionRate().Symbol(symbol).Do()
		if err != nil {
			return nil, err
		}
		feeRate.Maker, _ = decimal.NewFromString(res.MakerCommissionRate)
		feeRate.Taker, _ = decimal.NewFromString(res.TakerCommissionRate)
	default:
		return nil, ErrorAccountType
	}
	return &feeRate, nil
}

func (b BinanceTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	var positionList []*Position
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		return positionList, nil
	case BN_AC_FUTURE:
		res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
		if err != nil {
			return nil, err
		}
		risk, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFuturePositionRisk().Do()
		if err != nil {
			return nil, err
		}
		riskMap := map[string]mybinanceapi.FuturePositionRiskRow{}
		for _, r := range *risk {
			if r.PositionAmt != "0" {
				key := r.Symbol + r.PositionSide
				riskMap[key] = r
			}
		}
		//保证金率=维持保证金总额/账户保证金总额
		//保证金总额
		totalMarginBalance, _ := decimal.NewFromString(res.TotalMarginBalance)
		//维持保证金总额
		totalMaintMargin, _ := decimal.NewFromString(res.TotalMaintMargin)

		var marginRatio decimal.Decimal
		if !totalMarginBalance.IsZero() {
			marginRatio = totalMaintMargin.Div(totalMarginBalance)
		}

		for _, p := range res.Positions {
			if len(symbols) == 0 || stringInSlice(p.Symbol, symbols) {
				position := &Position{
					Exchange:               b.ExchangeType().String(),
					AccountType:            accountType,
					Symbol:                 p.Symbol,
					InitialMargin:          p.InitialMargin,
					MaintMargin:            p.MaintMargin,
					UnrealizedProfit:       p.UnrealizedProfit,
					PositionInitialMargin:  p.PositionInitialMargin,
					OpenOrderInitialMargin: p.OpenOrderInitialMargin,
					Leverage:               p.Leverage,
					MarginMode:             b.bnConverter.FromBNMarginMode(p.Isolated),
					EntryPrice:             p.EntryPrice,
					MaxNotional:            p.MaxNotional,
					PositionSide:           b.bnConverter.FromBNPositionSide(p.PositionSide),
					PositionAmt:            p.PositionAmt,
					UpdateTime:             p.UpdateTime,
				}
				//对于持仓量大于0的持仓,进行持仓风险查询
				if p.PositionAmt != "0" {
					key := p.Symbol + p.PositionSide
					if r, ok := riskMap[key]; ok {
						position.MarkPrice = r.MarkPrice
						position.LiquidationPrice = r.LiquidationPrice
						position.MarginRatio = marginRatio.String()
					}
				}
				positionList = append(positionList, position)
			}
		}
	case BN_AC_SWAP:
		res, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapAccount().Do()
		if err != nil {
			return nil, err
		}
		risk, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapPositionRisk().Do()
		if err != nil {
			return nil, err
		}

		riskMap := map[string]mybinanceapi.SwapPositionRiskResRow{}
		for _, r := range *risk {
			if r.PositionAmt != "0" {
				key := r.Symbol + r.PositionSide
				riskMap[key] = r
			}
		}

		//计算每个币种的保证金比率
		ratioMap := map[string]decimal.Decimal{}
		for _, a := range res.Assets {
			//维持保证金
			maintMargin, _ := decimal.NewFromString(a.MaintMargin)
			//保证金余额
			marginBalance, _ := decimal.NewFromString(a.MarginBalance)
			var ratio decimal.Decimal
			if !marginBalance.IsZero() {
				ratio = maintMargin.Div(marginBalance)
			}
			ratioMap[a.Asset] = ratio
		}

		for _, p := range res.Positions {
			if len(symbols) == 0 || stringInSlice(p.Symbol, symbols) {
				position := &Position{
					Exchange:               b.ExchangeType().String(),
					AccountType:            accountType,
					Symbol:                 p.Symbol,
					InitialMargin:          p.InitialMargin,
					MaintMargin:            p.MaintMargin,
					UnrealizedProfit:       p.UnrealizedProfit,
					PositionInitialMargin:  p.PositionInitialMargin,
					OpenOrderInitialMargin: p.OpenOrderInitialMargin,
					Leverage:               p.Leverage,
					MarginMode:             b.bnConverter.FromBNMarginMode(p.Isolated),
					EntryPrice:             p.EntryPrice,
					MaxNotional:            p.MaxQty,
					PositionSide:           b.bnConverter.FromBNPositionSide(p.PositionSide),
					PositionAmt:            p.PositionAmt,
					UpdateTime:             p.UpdateTime,
				}
				if p.PositionAmt != "0" {
					key := p.Symbol + p.PositionSide
					if r, ok := riskMap[key]; ok {
						position.MarkPrice = r.MarkPrice
						position.LiquidationPrice = r.LiquidationPrice
						//如果交易对名称包含保证金比率的key
						for key, ratio := range ratioMap {
							if strings.Contains(p.Symbol, key) {
								position.MarginRatio = ratio.String()
								break
							}
						}
					}
				}
				positionList = append(positionList, position)
			}
		}
	default:
		return nil, ErrorAccountType
	}

	return positionList, nil
}

func (b BinanceTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	var assetList []*Asset

	switch BinanceAccountType(accountType) {
	case BN_AC_FUNDING:
		res, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAssetGetFundingAsset().Do()
		if err != nil {
			return nil, err
		}
		for _, a := range *res {
			if len(currencies) == 0 || stringInSlice(a.Asset, currencies) {
				free, _ := decimal.NewFromString(a.Free)
				locked, _ := decimal.NewFromString(a.Locked)           //锁定
				freeze, _ := decimal.NewFromString(a.Freeze)           //冻结
				withdrawing, _ := decimal.NewFromString(a.Withdrawing) //提币中
				walletBalance := free.Add(locked).Add(freeze).Add(withdrawing)
				assetList = append(assetList, &Asset{
					Exchange:          b.ExchangeType().String(),                    //交易所
					AccountType:       accountType,                                  //账户类型
					Asset:             a.Asset,                                      //资产
					Free:              a.Free,                                       //可用余额
					Locked:            locked.Add(freeze).Add(withdrawing).String(), //locked=锁定+冻结+提币中
					WalletBalance:     walletBalance.String(),                       //钱包余额
					MaxWithdrawAmount: a.Free,                                       //最大可转
					UpdateTime:        time.Now().UnixMilli(),
				})
			}
		}
	case BN_AC_SPOT:
		res, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAccount().Do()
		if err != nil {
			return nil, err
		}
		for _, a := range res.Balance {
			if len(currencies) == 0 || stringInSlice(a.Asset, currencies) {
				free, _ := decimal.NewFromString(a.Free)
				lock, _ := decimal.NewFromString(a.Locked)
				walletBalance := free.Add(lock)
				assetList = append(assetList, &Asset{
					Exchange:          b.ExchangeType().String(), //交易所
					AccountType:       accountType,               //账户类型
					Asset:             a.Asset,                   //资产
					Free:              a.Free,                    //可用余额
					Locked:            a.Locked,                  //冻结余额
					WalletBalance:     walletBalance.String(),    //钱包余额
					MaxWithdrawAmount: a.Free,
					UpdateTime:        time.Now().UnixMilli(),
				})
			}
		}
	case BN_AC_MARGIN_CROSSED:
		res, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginAccount().Do()
		if err != nil {
			return nil, err
		}

		for _, a := range res.UserAssets {
			if len(currencies) == 0 || stringInSlice(a.Asset, currencies) {
				free, _ := decimal.NewFromString(a.Free)
				lock, _ := decimal.NewFromString(a.Locked)
				borrowed, _ := decimal.NewFromString(a.Borrowed)
				totalNetAssetOfBtc, _ := decimal.NewFromString(res.TotalNetAssetOfBtc)
				totalLiabilityOfBtc, _ := decimal.NewFromString(res.TotalLiabilityOfBtc)

				risk := decimal.NewFromFloat(999)
				if !totalLiabilityOfBtc.IsZero() {
					risk = totalNetAssetOfBtc.Div(totalLiabilityOfBtc).Add(decimal.NewFromFloat(1))
				}

				var maxTransferable decimal.Decimal
				if risk.GreaterThanOrEqual(decimal.NewFromFloat(2)) {
					maxTransferable = decimal.RequireFromString(a.NetAsset).Sub(decimal.RequireFromString(a.Borrowed).Add(decimal.RequireFromString(a.Interest)))
					if maxTransferable.LessThan(decimal.Zero) {
						maxTransferable = decimal.Zero
					}
				}

				walletBalance := free.Add(lock)
				assetList = append(assetList, &Asset{
					Exchange:          b.ExchangeType().String(), //交易所
					AccountType:       accountType,               //账户类型
					Asset:             a.Asset,                   //资产
					Borrowed:          borrowed.String(),         //已借
					Interest:          a.Interest,                //利息
					Free:              a.Free,                    //可用余额
					Locked:            a.Locked,                  //冻结余额
					WalletBalance:     walletBalance.String(),    //钱包余额
					MaxWithdrawAmount: maxTransferable.String(),  //最大可转
					UpdateTime:        time.Now().UnixMilli(),
				})
			}
		}
	case BN_AC_MARGIN_ISOLATED:
		var symbolList []string
		for _, c := range currencies {
			// example: c = btcusdt-btc
			split := strings.Split(c, "/")
			symbol := split[0]
			symbolList = append(symbolList, symbol)
		}
		symbols := strings.Join(symbolList, ",")
		res, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginIsolatedAccount().
			Symbols(symbols).
			Do()
		if err != nil {
			return nil, err
		}
		for _, a := range res.Assets {
			asset := a.Symbol + "/" + a.BaseAsset.Asset
			indexPrice, _ := decimal.NewFromString(a.IndexPrice)
			// BaseAsset
			if len(currencies) == 0 || stringInSlice(asset, currencies) {
				free, _ := decimal.NewFromString(a.BaseAsset.Free)         // 最大可用资产
				net, _ := decimal.NewFromString(a.BaseAsset.NetAsset)      // 净资产
				borrowed, _ := decimal.NewFromString(a.BaseAsset.Borrowed) // 已借
				interest, _ := decimal.NewFromString(a.BaseAsset.Interest) // 利息
				lock := net.Sub(free)                                      // 冻结

				// 质押率 = 抵押价值 / (负债 + 未归还利息)
				// 质押率 >= 2时能转出，因此 最大可转出 = 可用余额 - 质押率小于2的部分

				quoteBorrowed, quoteInterest := decimal.RequireFromString(a.QuoteAsset.Borrowed), decimal.RequireFromString(a.QuoteAsset.Interest)
				totalQuoteBorrowedAndInterest := quoteBorrowed.Add(quoteInterest).Mul(indexPrice)
				minCollateralValue := decimal.NewFromFloat(2).Mul(borrowed.Add(interest).Add(totalQuoteBorrowedAndInterest)) // 最小抵押价值(质押率<2的部分)
				QuoteAssetFree, _ := decimal.NewFromString(a.QuoteAsset.Free)                                                // 计算抵押价值
				QuoteCollateralValue := indexPrice.Mul(QuoteAssetFree)                                                       // 根据市价计算抵押价值
				var maxTransferable decimal.Decimal
				if QuoteCollateralValue.GreaterThan(minCollateralValue) {
					maxTransferable = free
				} else {
					maxTransferable = free.Sub(minCollateralValue)
				}
				assetList = append(assetList, &Asset{
					Exchange:          b.ExchangeType().String(), //交易所
					AccountType:       accountType,               //账户类型
					Asset:             asset,                     //资产
					Borrowed:          borrowed.String(),         //已借
					Interest:          interest.String(),         //利息
					Free:              free.String(),             //可用余额
					Locked:            lock.String(),             //冻结余额
					WalletBalance:     net.String(),              //钱包余额
					MaxWithdrawAmount: maxTransferable.String(),  //最大可转出余额
					UpdateTime:        time.Now().UnixMilli(),
				})
			}

			// QuoteAsset
			asset = a.Symbol + "/" + a.QuoteAsset.Asset
			if len(currencies) == 0 || stringInSlice(asset, currencies) {
				free, _ := decimal.NewFromString(a.QuoteAsset.Free)         // 最大可用资产
				net, _ := decimal.NewFromString(a.QuoteAsset.NetAsset)      // 净资产
				borrowed, _ := decimal.NewFromString(a.QuoteAsset.Borrowed) // 已借
				interest, _ := decimal.NewFromString(a.QuoteAsset.Interest) // 利息
				lock := net.Sub(free)                                       // 冻结

				// 质押率 = 抵押价值 / (负债 + 未归还利息)
				// 质押率 >= 2时能转出，因此 最大可转出 = 可用余额 - 质押率小于2的部分
				baseBorrowed, baseInterest := decimal.RequireFromString(a.BaseAsset.Borrowed), decimal.RequireFromString(a.BaseAsset.Interest)
				totalBaseBorrowedAndInterest := baseBorrowed.Add(baseInterest).Mul(indexPrice)
				minCollateralValue := decimal.NewFromFloat(2).Mul(borrowed.Add(interest).Add(totalBaseBorrowedAndInterest)) // 最小抵押价值(质押率<2的部分)
				QuoteCollateralValue, _ := decimal.NewFromString(a.BaseAsset.Free)                                          // 根据实时市价计算抵押价值

				var maxTransferable decimal.Decimal
				if QuoteCollateralValue.Div(indexPrice).GreaterThanOrEqual(minCollateralValue) {
					maxTransferable = free
				} else {
					maxTransferable = free.Sub(minCollateralValue)
				}
				assetList = append(assetList, &Asset{
					Exchange:          b.ExchangeType().String(), //交易所
					AccountType:       accountType,               //账户类型
					Asset:             asset,                     //资产
					Borrowed:          borrowed.String(),         //已借
					Interest:          interest.String(),         //利息
					Free:              free.String(),             //可用余额
					Locked:            lock.String(),             //冻结余额
					WalletBalance:     net.String(),              //钱包余额
					MaxWithdrawAmount: maxTransferable.String(),  //最大可转出余额
					UpdateTime:        time.Now().UnixMilli(),
				})
			}
		}
	case BN_AC_FUTURE:
		res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
		if err != nil {
			return nil, err
		}
		for _, a := range res.Assets {
			if len(currencies) == 0 || stringInSlice(a.Asset, currencies) {
				assetList = append(assetList, &Asset{
					Exchange:               b.ExchangeType().String(), //交易所
					AccountType:            accountType,               //账户类型
					Asset:                  a.Asset,                   //资产
					WalletBalance:          a.WalletBalance,           //余额
					UnrealizedProfit:       a.UnrealizedProfit,        //未实现盈亏
					MarginBalance:          a.MarginBalance,           //保证金余额
					MaintMargin:            a.MaintMargin,             //维持保证金
					InitialMargin:          a.InitialMargin,           //当前所需起始保证金
					PositionInitialMargin:  a.PositionInitialMargin,   //持仓所需起始保证金(基于最新标记价格)
					OpenOrderInitialMargin: a.OpenOrderInitialMargin,  //当前挂单所需起始保证金(基于最新标记价格)
					MaxWithdrawAmount:      a.MaxWithdrawAmount,       //最大可转出余额
					CrossWalletBalance:     a.CrossWalletBalance,      //全仓账户余额
					CrossUnPnl:             a.CrossUnPnl,              //全仓持仓未实现盈亏
					AvailableBalance:       a.AvailableBalance,        //可用余额
					UpdateTime:             time.Now().UnixMilli(),
				})
			}
		}
	case BN_AC_SWAP:
		res, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapAccount().Do()
		if err != nil {
			return nil, err
		}
		for _, a := range res.Assets {
			if len(currencies) == 0 || stringInSlice(a.Asset, currencies) {
				assetList = append(assetList, &Asset{
					Exchange:               b.ExchangeType().String(), //交易所
					AccountType:            accountType,               //账户类型
					Asset:                  a.Asset,                   //资产
					WalletBalance:          a.WalletBalance,           //余额
					UnrealizedProfit:       a.UnrealizedProfit,        //未实现盈亏
					MarginBalance:          a.MarginBalance,           //保证金余额
					MaintMargin:            a.MaintMargin,             //维持保证金
					InitialMargin:          a.InitialMargin,           //当前所需起始保证金
					PositionInitialMargin:  a.PositionInitialMargin,   //持仓所需起始保证金(基于最新标记价格)
					OpenOrderInitialMargin: a.OpenOrderInitialMargin,  //当前挂单所需起始保证金(基于最新标记价格)
					MaxWithdrawAmount:      a.MaxWithdrawAmount,       //最大可转出余额
					CrossWalletBalance:     a.CrossWalletBalance,      //全仓账户余额
					CrossUnPnl:             a.CrossUnPnl,              //全仓持仓未实现盈亏
					AvailableBalance:       a.AvailableBalance,        //可用余额
					UpdateTime:             time.Now().UnixMilli(),
				})
			}
		}
	default:
		return nil, ErrorAccountType
	}
	return assetList, nil
}

func (b BinanceTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAssetTransferPost()

	FromAsset := b.bnConverter.ToBNAssetType(req.From)
	ToAsset := b.bnConverter.ToBNAssetType(req.To)
	BNTransferType := FromAsset + "_" + ToAsset
	api.Type(mybinanceapi.AssetTransferType(BNTransferType))

	api.Asset(req.Asset).Amount(req.Amount)

	res, err := api.Do()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var assetTransfers []*AssetTransfer
	tranId := strconv.FormatInt(res.TranId, 10)
	assetTransfers = append(assetTransfers, &AssetTransfer{
		Exchange: b.ExchangeType().String(),
		TranId:   tranId,
		Asset:    req.Asset,
		From:     req.From,
		To:       req.To,
		Amount:   req.Amount.String(),
		Status:   "",
		ClientId: "",
	})

	return assetTransfers, nil
}

func (b BinanceTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAssetTransferGet()
	Type := b.bnConverter.ToBNAssetType(req.From) + "_" + b.bnConverter.ToBNAssetType(req.To)
	api.Type(mybinanceapi.AssetTransferType(Type))
	api.StartTime(req.StartTime).EndTime(req.EndTime)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	var QueryAssetTransfers []*QueryAssetTransfer
	for _, r := range res.Rows {
		if req.Asset != "" && r.Asset != req.Asset {
			continue
		}
		QueryAssetTransfers = append(QueryAssetTransfers, &QueryAssetTransfer{
			TranId: strconv.FormatInt(r.TranId, 10),
			Asset:  r.Asset,
			Amount: stringToDecimal(r.Amount),
			From:   b.bnConverter.FromBNAssetType(strings.Split(r.Type, "_")[0]),
			To:     b.bnConverter.FromBNAssetType(strings.Split(r.Type, "_")[1]),
			Status: b.bnConverter.FromBinanceTransferStatus(r.Status),
		})
	}

	return QueryAssetTransfers, nil
}
