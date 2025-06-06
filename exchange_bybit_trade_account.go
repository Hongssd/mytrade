package mytrade

import (
	"strconv"
	"strings"

	"github.com/Hongssd/mybybitapi"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BybitTradeAccount struct {
	ExchangeBase

	bybitConverter BybitEnumConverter
	apiKey         string
	secretKey      string
}

func (b BybitTradeAccount) GetAccountMode() (AccountMode, error) {
	res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().NewAccountInfo().Do()
	if err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	return b.bybitConverter.FromBYBITAccountMode(res.Result.MarginMode), nil
}

func (b BybitTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().NewAccountInfo().Do()
	if err != nil {
		return MARGIN_MODE_UNKNOWN, err
	}
	if res.Result.UnifiedMarginStatus == 1 || accountType == BYBIT_AC_INVERSE.String() {
		res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewPositionList().Category(accountType).Symbol(symbol).Do()
		if err != nil {
			return MARGIN_MODE_UNKNOWN, err
		}
		for _, p := range res.Result.List {
			if p.Symbol == symbol && b.bybitConverter.FromBYBITPositionSide(p.PositionIdx) == positionSide {
				return b.bybitConverter.FromBYBITMarginMode(p.TradeMode), nil
			}
		}
		return MARGIN_MODE_UNKNOWN, ErrorPositionNotFound
	} else {
		accountMode, err := b.GetAccountMode()
		if err != nil {
			return MARGIN_MODE_UNKNOWN, err
		}
		if accountMode == ACCOUNT_MODE_SINGLE_MARGIN {
			return MARGIN_MODE_ISOLATED, nil
		} else {
			return MARGIN_MODE_CROSSED, nil
		}
	}
}

func (b BybitTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	if accountType == BYBIT_AC_SPOT.String() {
		return POSITION_MODE_ONEWAY, nil
	}

	res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
		NewPositionList().Category(accountType).Symbol(symbol).Do()
	if err != nil {
		return POSITION_MODE_UNKNOWN, err
	}
	if len(res.Result.List) == 1 {
		//单向
		return POSITION_MODE_ONEWAY, nil
	} else if len(res.Result.List) == 2 {
		//双向
		return POSITION_MODE_HEDGE, nil
	}
	return POSITION_MODE_UNKNOWN, ErrorPositionNotFound
}

func (b BybitTradeAccount) GetLeverage(accountType, symbol string,
	marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	if accountType == BYBIT_AC_SPOT.String() {
		res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewSpotMarginTradeState().Do()
		if err != nil {
			return decimal.Zero, err
		}
		leverage, _ := decimal.NewFromString(res.Result.SpotLeverage)
		return leverage, nil
	} else {
		res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewPositionList().Category(accountType).Symbol(symbol).Do()
		if err != nil {
			return decimal.Zero, err
		}
		for _, p := range res.Result.List {
			// 如果为统一账户的反向合约
			if accountType == BYBIT_AC_INVERSE.String() {
				if p.Symbol == symbol &&
					b.bybitConverter.FromBYBITPositionSide(p.PositionIdx) == positionSide &&
					b.bybitConverter.FromBYBITMarginMode(p.TradeMode) == marginMode {
					leverage, _ := decimal.NewFromString(p.Leverage)
					return leverage, nil
				}
			} else {
				if p.Symbol == symbol &&
					b.bybitConverter.FromBYBITPositionSide(p.PositionIdx) == positionSide {
					leverage, _ := decimal.NewFromString(p.Leverage)
					return leverage, nil
				}
			}
		}

	}
	return decimal.Zero, ErrorPositionNotFound
}

func (b BybitTradeAccount) SetAccountMode(mode AccountMode) error {
	_, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
		NewAccountSetMarginMode().SetMarginMode(b.bybitConverter.ToBYBITAccountMode(mode)).Do()
	if err != nil {
		return err
	}
	return nil
}

func (b BybitTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	//if accountType == BYBIT_AC_SPOT.String() {
	//	return ErrorNotSupport
	//}

	res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().NewAccountInfo().Do()
	if err != nil {
		return err
	}
	if res.Result.UnifiedMarginStatus == 1 || accountType == BYBIT_AC_INVERSE.String() {
		//经典账号
		_, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewPositionSwitchIsolated().Category(accountType).
			Symbol(symbol).TradeMode(b.bybitConverter.ToBYBITMarginMode(mode)).BuyLeverage("10").SellLeverage("10").Do()
		if err != nil {
			return err
		}
		return nil
	} else {
		//统一账号，设置账户模式
		if mode == MARGIN_MODE_ISOLATED {
			//统一逐仓保证金模式
			return b.SetAccountMode(ACCOUNT_MODE_SINGLE_MARGIN)
		} else if mode == MARGIN_MODE_CROSSED {
			//统一全仓保证金模式
			return b.SetAccountMode(ACCOUNT_MODE_MULTI_MARGIN)
		}
	}

	return nil
}

func (b BybitTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	if accountType == BYBIT_AC_SPOT.String() {
		return ErrorNotSupport
	}
	currentMode, err := b.GetPositionMode(accountType, symbol)
	if err != nil {
		return err
	}
	if currentMode == mode {
		return nil
	}

	_, err = mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
		NewPositionSwitchMode().Category(accountType).Symbol(symbol).
		Mode(b.bybitConverter.ToBYBITPositionMode(mode)).Do()
	if err != nil {
		return err
	}
	return nil
}

func (b BybitTradeAccount) SetLeverage(accountType, symbol string,
	marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {

	// spot cross leverage
	if accountType == BYBIT_AC_SPOT.String() {
		if marginMode != MARGIN_MODE_CROSSED {
			return ErrorNotSupport
		}

		api := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewSpotMarginTradeSetLeverage().Leverage(leverage.String())
		_, err := api.Do()
		if err != nil {
			return err
		}
	} else {
		positionMode, err := b.GetPositionMode(accountType, symbol)
		if err != nil {
			return err
		}

		nowLeverage, err := b.GetLeverage(accountType, symbol, marginMode, positionSide)
		if err != nil {
			return err
		}

		if nowLeverage.Equal(leverage) {
			return nil
		}

		api := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewPositionSetLeverage().Category(accountType).Symbol(symbol)

		//單倉模式: 經典帳戶和統一帳戶的buyLeverage 必須等於sellLeverage
		//雙倉模式:
		//經典帳戶和統一帳戶(逐倉模式)buyLeverage可以與sellLeverage不想等;
		//統一帳戶(全倉模式)的buyLeverage 必須等於sellLeverage
		accountMode, err := b.GetAccountMode()
		if err != nil {
			return err
		}

		if positionMode == POSITION_MODE_ONEWAY || accountMode != ACCOUNT_MODE_SINGLE_MARGIN {
			api.BuyLeverage(leverage.String()).SellLeverage(leverage.String())
		} else {
			if positionSide == POSITION_SIDE_LONG {
				api.BuyLeverage(leverage.String())
			} else if positionSide == POSITION_SIDE_SHORT {
				api.SellLeverage(leverage.String())
			}
		}
		_, err = api.Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (b BybitTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	var feeRate FeeRate
	res, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
		NewAccountFeeRate().Category(accountType).Symbol(symbol).Do()
	if err != nil {
		return nil, err
	}
	if len(res.Result.List) != 1 {
		return nil, ErrorSymbolNotFound
	}
	feeRate.Maker, _ = decimal.NewFromString(res.Result.List[0].MakerFeeRate)
	feeRate.Taker, _ = decimal.NewFromString(res.Result.List[0].TakerFeeRate)
	return &feeRate, nil
}

func (b BybitTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	if BybitAccountType(accountType) == BYBIT_AC_SPOT {
		return []*Position{}, nil
	}

	api := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
		NewPositionList().Category(accountType)
	if len(symbols) == 1 {
		api.Symbol(symbols[0])
	} else if len(symbols) == 0 {
		switch BybitAccountType(accountType) {
		case BYBIT_AC_LINEAR:
			api.SettleCoin("USDT")
		case BYBIT_AC_INVERSE:
			api.SettleCoin("BTC")
		}
	}
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	//log.Warn(len(res.Result.List))
	var positions []*Position
	for _, p := range res.Result.List {
		if len(symbols) == 0 || stringInSlice(p.Symbol, symbols) {
			updateTime, _ := strconv.ParseInt(p.UpdatedTime, 10, 64)
			if p.LiqPrice == "" {
				p.LiqPrice = "0"
			}
			amt, _ := decimal.NewFromString(p.Size)

			if p.PositionIdx == BYBIT_POSITION_SIDE_BOTH && p.Side == "Sell" && amt.GreaterThan(decimal.Zero) {
				//单向持仓模式下持有空仓，数量设置为负数
				amt = amt.Mul(decimal.NewFromInt(-1))
			}
			position := &Position{
				Exchange:               b.ExchangeType().String(),
				AccountType:            accountType,
				Symbol:                 p.Symbol,
				InitialMargin:          p.PositionIM,
				MaintMargin:            p.PositionMM,
				UnrealizedProfit:       p.UnrealisedPnl,
				PositionInitialMargin:  p.PositionIM,
				OpenOrderInitialMargin: p.PositionMM,
				Leverage:               p.Leverage,
				MarginMode:             b.bybitConverter.FromBYBITMarginMode(p.TradeMode),
				EntryPrice:             p.AvgPrice,
				MaxNotional:            "0",
				PositionSide:           b.bybitConverter.FromBYBITPositionSide(p.PositionIdx),
				PositionAmt:            amt.String(),
				MarkPrice:              p.MarkPrice,
				LiquidationPrice:       p.LiqPrice,
				MarginRatio:            p.Leverage,
				UpdateTime:             updateTime,
			}
			positions = append(positions, position)
		}
	}
	return positions, nil
}

func (b BybitTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {

	var assets []*Asset
	if accountType == BYBIT_AC_FUNDING.String() {
		api := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewAssetTransferQueryAccountCoinsBalance()
		res, err := api.AccountType(accountType).Do()
		if err != nil {
			return nil, err
		}

		for _, a := range res.Result.Balance {
			walletBalance, err := decimal.NewFromString(a.WalletBalance)
			if err != nil {
				return nil, err
			}
			free, err := decimal.NewFromString(a.TransferBalance)
			if err != nil {
				return nil, err
			}
			locked := walletBalance.Sub(free)

			asset := &Asset{
				Exchange:               b.ExchangeType().String(),
				AccountType:            accountType,
				Asset:                  a.Coin,
				Free:                   a.TransferBalance,
				Locked:                 locked.String(),
				WalletBalance:          a.WalletBalance,
				UnrealizedProfit:       "0",
				MarginBalance:          "0",
				MaintMargin:            "0",
				InitialMargin:          "0",
				PositionInitialMargin:  "0",
				OpenOrderInitialMargin: "0",
				CrossWalletBalance:     "0",
				CrossUnPnl:             "0",
				AvailableBalance:       a.TransferBalance,
				MaxWithdrawAmount:      a.TransferBalance,
				MarginAvailable:        false,
				UpdateTime:             res.Time,
			}
			assets = append(assets, asset)
		}
	} else {
		acType := mybybitapi.ACCT_UNIFIED
		if accountType == BYBIT_AC_INVERSE.String() {
			acType = mybybitapi.ACCT_CONTRACT
		}

		api := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().
			NewAccountWalletBalance().AccountType(acType.String())
		if len(currencies) == 1 {
			api.Coin(currencies[0])
		} else if len(currencies) > 1 {
			coins := strings.Join(currencies, ",")
			api.Coin(coins)
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		availableWithdrawalMap := map[string]string{}

		CoinList := []string{}
		for _, a := range res.Result.List[0].Coin {
			CoinList = append(CoinList, a.Coin)
		}
		//每20个币种查询一次可划转余额，若不足20个币种，则查询所有币种
		for i := 0; i < len(CoinList); i += 20 {
			if i+20 > len(CoinList) {
				CoinList = CoinList[i:]
			} else {
				CoinList = CoinList[i : i+20]
			}
			r, err := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().NewAccountWithdrawal().CoinName(strings.Join(CoinList, ",")).Do()
			if err != nil {
				return nil, err
			}
			for coin, amount := range r.Result.AvailableWithdrawalMap {
				availableWithdrawalMap[coin] = amount
			}
		}

		for _, a := range res.Result.List[0].Coin {

			tpIm, _ := decimal.NewFromString(a.TotalPositionIM)
			toIm, _ := decimal.NewFromString(a.TotalOrderIM)
			im := tpIm.Add(toIm)

			eq, _ := decimal.NewFromString(a.Equity)
			lock, _ := decimal.NewFromString(a.Locked)
			avb := eq.Sub(lock)

			availableWithdrawal, ok := availableWithdrawalMap[a.Coin]
			if !ok {
				availableWithdrawal = "0"
			}

			asset := &Asset{
				Exchange:               b.ExchangeType().String(),
				AccountType:            accountType,
				Asset:                  a.Coin,
				Free:                   avb.String(),
				Locked:                 a.Locked,
				Borrowed:               a.BorrowAmount,
				Interest:               a.AccruedInterest,
				WalletBalance:          a.WalletBalance,
				UnrealizedProfit:       a.UnrealisedPnl,
				MarginBalance:          a.Equity,
				MaintMargin:            a.TotalPositionMM,
				InitialMargin:          im.String(),
				PositionInitialMargin:  a.TotalPositionIM,
				OpenOrderInitialMargin: a.TotalOrderIM,
				CrossWalletBalance:     a.WalletBalance,
				CrossUnPnl:             a.UnrealisedPnl,
				AvailableBalance:       avb.String(),
				MaxWithdrawAmount:      availableWithdrawal,
				MarginAvailable:        false,
				UpdateTime:             res.Time,
			}
			assets = append(assets, asset)
		}
	}

	return assets, nil
}

// 资金划转
func (b BybitTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	api := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().NewAssetTransferInterTransfer()

	transferId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	api.TransferId(transferId.String())

	From := b.bybitConverter.ToBYBITAssetType(req.From)
	To := b.bybitConverter.ToBYBITAssetType(req.To)
	if From == "" || To == "" {
		return nil, err
	}
	api.FromAccountType(From).ToAccountType(To)

	api.Coin(req.Asset).Amount(req.Amount.String())

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	var assetTransfers []*AssetTransfer
	d := res.Result
	assetTransfers = append(assetTransfers, &AssetTransfer{
		Exchange: b.ExchangeType().String(),
		TranId:   d.TransferId,
		Asset:    req.Asset,
		From:     req.From,
		To:       req.To,
		Amount:   req.Amount.String(),
		Status:   d.Status,
		ClientId: "",
	})

	return assetTransfers, nil
}

func (b BybitTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	api := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient().NewAssetTransferQueryInterTransferList()
	api.StartTime(req.StartTime).EndTime(req.EndTime)
	api.Limit(1000)

	// api.Coin 接口有问题 弃用
	//if req.Asset != "" {
	//	fmt.Println(req.Asset)
	//	api.Coin(req.Asset)
	//}

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	var assetTransfers []*QueryAssetTransfer
	var previousCursor string
	for res.Result.NextPageCursor != previousCursor {
		for _, d := range res.Result.List {
			if b.bybitConverter.FromBYBITAssetType(d.FromAccountType) != req.From ||
				b.bybitConverter.FromBYBITAssetType(d.ToAccountType) != req.To {
				continue
			}
			if req.Asset != "" && d.Coin != req.Asset {
				continue
			}
			assetTransfers = append(assetTransfers, &QueryAssetTransfer{
				TranId: d.TransferId,
				Asset:  d.Coin,
				Amount: stringToDecimal(d.Amount),
				From:   b.bybitConverter.FromBYBITAssetType(d.FromAccountType),
				To:     b.bybitConverter.FromBYBITAssetType(d.ToAccountType),
				Status: b.bybitConverter.FromBYBITTransferStatus(d.Status),
			})
		}
		res, err = api.Cursor(res.Result.NextPageCursor).Do()
		if err != nil {
			return nil, err
		}

		previousCursor = res.Result.NextPageCursor
	}

	return assetTransfers, nil
}
