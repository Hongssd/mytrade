package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type BinanceTradeAccount struct {
	exchangeBase

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
				locked, _ := decimal.NewFromString(a.Locked)           //锁定
				freeze, _ := decimal.NewFromString(a.Freeze)           //冻结
				withdrawing, _ := decimal.NewFromString(a.Withdrawing) //提币中
				assetList = append(assetList, &Asset{
					Exchange:          b.ExchangeType().String(),                    //交易所
					AccountType:       accountType,                                  //账户类型
					Asset:             a.Asset,                                      //资产
					Free:              a.Free,                                       //可用余额
					Locked:            locked.Add(freeze).Add(withdrawing).String(), //locked=锁定+冻结+提币中
					MaxWithdrawAmount: a.Free,
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
					Exchange:      b.ExchangeType().String(), //交易所
					AccountType:   accountType,               //账户类型
					Asset:         a.Asset,                   //资产
					Free:          a.Free,                    //可用余额
					Locked:        a.Locked,                  //冻结余额
					WalletBalance: walletBalance.String(),    //钱包余额
					UpdateTime:    time.Now().UnixMilli(),
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
