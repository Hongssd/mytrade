package mytrade

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

type GateTradeAccount struct {
	ExchangeBase

	gateConverter GateEnumConverter
	apiKey        string
	secretKey     string
	passphrase    string
}

func (a GateTradeAccount) GetAccountMode() (AccountMode, error) {
	res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).
		PrivateRestClient().NewPrivateRestUnifiedUnifiedModeGet().Do()
	if err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	return a.gateConverter.FromGateAccountMode(res.Data.Mode), nil
}

func (a GateTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	positions, err := a.GetPositions(accountType, symbol)
	if err != nil {
		return MARGIN_MODE_UNKNOWN, err
	}
	marginMode := MARGIN_MODE_UNKNOWN
	for _, p := range positions {
		if p.PositionSide == positionSide {
			marginMode = p.MarginMode
		}
	}

	return marginMode, nil
}

func (a GateTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {

	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		return POSITION_MODE_ONEWAY, nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		split := strings.Split(symbol, "_")
		if len(split) != 2 {
			return POSITION_MODE_UNKNOWN, errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
			NewPrivateRestFuturesSettleAccounts().Settle(settle).Do()
		if err != nil {
			return POSITION_MODE_UNKNOWN, err
		}
		return a.gateConverter.FromGatePositionMode(res.Data.InDualMode), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		split := strings.Split(symbol, "_")
		if len(split) < 3 {
			return POSITION_MODE_UNKNOWN, errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
			NewPrivateRestDeliverySettleAccounts().Settle(strings.ToLower(settle)).Do()
		if err != nil {
			return POSITION_MODE_UNKNOWN, err
		}
		return a.gateConverter.FromGatePositionMode(res.Data.InDualMode), nil
	}

	return POSITION_MODE_UNKNOWN, nil
}

func (a GateTradeAccount) GetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	if symbol == "" {
		return decimal.Zero, errors.New("symbol error")
	}
	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		return decimal.Zero, ErrorNotSupport
	case GATE_ACCOUNT_TYPE_MARGIN:
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PublicRestClient().NewPublicRestMarginUniCurrencyPairsCurrencyPair().
			CurrencyPair(symbol).Do()
		if err != nil {
			return decimal.Zero, err
		}
		return decimal.RequireFromString(res.Data.Leverage), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		split := strings.Split(symbol, "_")
		if len(split) < 2 {
			return decimal.Zero, errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
			NewPrivateRestFuturesSettlePositions().Settle(settle).Holding(false).
			Do()
		if err != nil {
			return decimal.Zero, err
		}
		for _, p := range res.Data {
			if p.Contract == symbol {
				if p.Leverage != "0" && marginMode == MARGIN_MODE_ISOLATED {
					return decimal.RequireFromString(p.Leverage), nil
				} else if p.Leverage == "0" && marginMode == MARGIN_MODE_CROSSED {
					return decimal.RequireFromString(p.CrossLeverageLimit), nil
				}
			}
		}
		return decimal.NewFromFloat(10), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		split := strings.Split(symbol, "_")
		if len(split) < 3 {
			return decimal.Zero, errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
			NewPrivateRestDeliverySettlePositions().Settle(settle).
			Do()
		if err != nil {
			return decimal.Zero, err
		}

		for _, p := range res.Data {
			if p.Contract == symbol {
				if p.Leverage != "0" && marginMode == MARGIN_MODE_ISOLATED {
					return decimal.RequireFromString(p.Leverage), nil
				} else if p.Leverage == "0" && marginMode == MARGIN_MODE_CROSSED {
					return decimal.RequireFromString(p.CrossLeverageLimit), nil
				}
			}
		}
		return decimal.NewFromFloat(10), nil
	}
	return decimal.Zero, ErrorNotSupport
}

func (a GateTradeAccount) SetAccountMode(mode AccountMode) error {
	currentAccountMode, err := a.GetAccountMode()
	if err != nil {
		return err
	}
	if currentAccountMode == mode {
		return nil
	}
	_, err = mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
		NewPrivateRestUnifiedUnifiedModePut().
		Mode(a.gateConverter.ToGateAccountMode(mode)).Do()
	if err != nil {
		return err
	}
	return nil
}

func (a GateTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {

	//获取仓位模式，设置仓位方向
	positionMode, err := a.GetPositionMode(accountType, symbol)
	if err != nil {
		return err
	}
	positionSide := POSITION_SIDE_BOTH
	switch positionMode {
	case POSITION_MODE_ONEWAY:
		positionSide = POSITION_SIDE_BOTH
	case POSITION_MODE_HEDGE:
		positionSide = POSITION_SIDE_LONG
	}

	//获取当前仓位方向的保证金模式
	currentMarginMode, err := a.GetMarginMode(accountType, symbol, positionSide)
	if err != nil {
		return err
	}

	//新旧保证金模式相同，无需修改
	if currentMarginMode == mode {
		return nil
	}

	//获取当前仓位方向的杠杆倍数，设置新保证金模式下的杠杆倍数
	leverage, err := a.GetLeverage(accountType, symbol, currentMarginMode, positionSide)
	if err != nil {
		return err
	}
	return a.SetLeverage(accountType, symbol, mode, positionSide, leverage)
}

func (a GateTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	// 仅支持永续合约
	log.Info(accountType)
	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_FUTURES:
		split := strings.Split(symbol, "_")
		if len(split) < 2 {
			return errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		_, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
			NewPrivateRestFuturesSettleDualMode().Settle(settle).
			DualMode(a.gateConverter.ToGatePositionMode(mode)).Do()
		if err != nil {
			return err
		}

		return nil
	default:
		return ErrorNotSupport
	}
}

func (a GateTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {
	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_FUTURES:
		split := strings.Split(symbol, "_")
		if len(split) < 2 {
			return errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		positionMode, err := a.GetPositionMode(accountType, symbol)
		if err != nil {
			return err
		}
		switch positionMode {
		case POSITION_MODE_ONEWAY:
			api := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
				NewPrivateRestFuturesSettlePositionsContractLeverage().Settle(settle).Contract(symbol)

			if marginMode == MARGIN_MODE_CROSSED {
				api.Leverage("0")
				api.CrossLeverageLimit(leverage.String())
			} else {
				api.Leverage(leverage.String())
			}
			_, err := api.Do()
			if err != nil {
				return err
			}
			return nil

		case POSITION_MODE_HEDGE:
			api := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
				NewPrivateRestFuturesSettleDualCompPositionsContractLeverage().Settle(settle).Contract(symbol)

			if marginMode == MARGIN_MODE_CROSSED {
				api.Leverage("0")
				api.CrossLeverageLimit(leverage.String())
			} else {
				api.Leverage(leverage.String())
			}
			_, err := api.Do()
			if err != nil {
				return err
			}
			return nil
		}

	}
	return ErrorNotSupport
}

func (a GateTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	var feeRate FeeRate
	api := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().NewPrivateRestWalletFee().CurrencyPair(symbol)
	if GateAccountType(accountType) == GATE_ACCOUNT_TYPE_FUTURES || GateAccountType(accountType) == GATE_ACCOUNT_TYPE_DELIVERY {
		split := strings.Split(symbol, "_")
		settle := strings.ToLower(split[1])
		api.Settle(settle)
	}
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_SPOT, GATE_ACCOUNT_TYPE_MARGIN:
		feeRate.Maker = decimal.RequireFromString(res.Data.MakerFee)
		feeRate.Taker = decimal.RequireFromString(res.Data.TakerFee)
	case GATE_ACCOUNT_TYPE_FUTURES:
		feeRate.Maker = decimal.RequireFromString(res.Data.FuturesMakerFee)
		feeRate.Taker = decimal.RequireFromString(res.Data.FuturesTakerFee)
	case GATE_ACCOUNT_TYPE_DELIVERY:
		feeRate.Maker = decimal.RequireFromString(res.Data.DeliveryMakerFee)
		feeRate.Taker = decimal.RequireFromString(res.Data.DeliveryTakerFee)
	}
	return &feeRate, nil
}

func (a GateTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	var positions []*Position
	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_FUTURES:
		var errG errgroup.Group
		settles := []string{"usdt", "btc"}
		for _, settle := range settles {
			settle := settle
			errG.Go(func() error {
				res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().NewPrivateRestFuturesSettlePositions().Settle(settle).Do()
				if err != nil {
					return nil
				}
				for _, p := range res.Data {
					var marginMode MarginMode
					var leverage string
					if p.Leverage == "0" {
						marginMode = MARGIN_MODE_CROSSED
						leverage = p.CrossLeverageLimit
					} else {
						marginMode = MARGIN_MODE_ISOLATED
						leverage = p.Leverage
					}

					positionSide := a.gateConverter.FromGatePositionSide(p.Mode)

					positions = append(positions, &Position{
						Exchange:               a.ExchangeType().String(),
						AccountType:            accountType,
						Symbol:                 p.Contract,
						MarginCcy:              settle,
						InitialMargin:          p.InitialMargin,
						MaintMargin:            p.MaintenanceMargin,
						UnrealizedProfit:       p.UnrealisedPnl,
						PositionInitialMargin:  p.InitialMargin,
						OpenOrderInitialMargin: p.InitialMargin,
						Leverage:               leverage,
						MarginMode:             marginMode,
						EntryPrice:             p.EntryPrice,
						MaxNotional:            "0",
						PositionSide:           positionSide,
						PositionAmt:            decimal.NewFromInt(p.Size).String(),
						MarkPrice:              p.MarkPrice,
						LiquidationPrice:       p.LiqPrice,
						MarginRatio:            p.MaintenanceRate,
						UpdateTime:             p.UpdateTime * 1000,
					})
				}
				return nil
			})
		}

		if err := errG.Wait(); err != nil {
			return nil, err
		}
	case GATE_ACCOUNT_TYPE_DELIVERY:
		settles := []string{"usdt"}
		for _, settle := range settles {
			res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
				NewPrivateRestDeliverySettlePositions().Settle(settle).
				Do()
			if err != nil {
				return nil, err
			}

			for _, p := range res.Data {

				var marginMode MarginMode
				var leverage string
				if p.Leverage == "0" {
					marginMode = MARGIN_MODE_CROSSED
					leverage = p.CrossLeverageLimit
				} else {
					marginMode = MARGIN_MODE_ISOLATED
					leverage = p.Leverage
				}

				positionSide := a.gateConverter.FromGatePositionSide(p.Mode)

				positions = append(positions, &Position{
					Exchange:               a.ExchangeType().String(),
					AccountType:            accountType,
					Symbol:                 p.Contract,
					MarginCcy:              settle,
					InitialMargin:          p.InitialMargin,
					MaintMargin:            p.MaintenanceMargin,
					UnrealizedProfit:       p.UnrealisedPnl,
					PositionInitialMargin:  p.InitialMargin,
					OpenOrderInitialMargin: p.InitialMargin,
					Leverage:               leverage,
					MarginMode:             marginMode,
					EntryPrice:             p.EntryPrice,
					MaxNotional:            "0",
					PositionSide:           positionSide,
					PositionAmt:            decimal.NewFromInt(p.Size).String(),
					MarkPrice:              p.MarkPrice,
					LiquidationPrice:       p.LiqPrice,
					MarginRatio:            p.MaintenanceRate,
					UpdateTime:             p.UpdateTime * 1000,
				})

			}
		}
	default:
		return nil, ErrorPositionNotFound
	}

	filterPositions := []*Position{}
	if len(symbols) > 0 {
		for _, p := range positions {
			// log.Warn(p.Symbol, symbols)
			isExist := false
			for _, symbol := range symbols {
				if p.Symbol == symbol {
					isExist = true
					break
				}
			}
			if isExist {
				filterPositions = append(filterPositions, p)
			}
		}
		return filterPositions, nil
	}
	return positions, nil
}

func (a GateTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	var assets []*Asset

	// 现货资产
	switch accountType {
	case GATE_ASSET_TYPE_SPOT:
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().NewPrivateRestSpotInstruments().Do()
		if err != nil {
			return nil, err
		}
		for _, d := range res.Data {
			assets = append(assets, &Asset{
				Exchange:               a.ExchangeType().String(),
				AccountType:            accountType,
				Asset:                  d.Currency,
				Borrowed:               "0",
				Interest:               "0",
				Free:                   d.Available,
				Locked:                 d.Locked,
				WalletBalance:          d.Available,
				UnrealizedProfit:       "0",
				MarginBalance:          "0",
				MaintMargin:            "0",
				InitialMargin:          "0",
				PositionInitialMargin:  "0",
				OpenOrderInitialMargin: "0",
				CrossWalletBalance:     "0",
				CrossUnPnl:             "0",
				AvailableBalance:       d.Available,
				MaxWithdrawAmount:      "0",
				MarginAvailable:        false,
				UpdateTime:             time.Now().UnixMilli(),
			})
		}
	case GATE_ASSET_TYPE_UNFIED:
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).
			PrivateRestClient().NewPrivateRestUnifiedAccounts().Do()
		if err != nil {
			return nil, err
		}
		for asset, b := range res.Data.Balances {
			free, _ := decimal.NewFromString(b.Available)
			locked, _ := decimal.NewFromString(b.Freeze)
			walletBalance := free.Add(locked)
			assets = append(assets, &Asset{
				Exchange:      a.ExchangeType().String(), //交易所
				AccountType:   accountType,               //账户类型
				Asset:         asset,                     //资产
				WalletBalance: walletBalance.String(),    //钱包余额
				Free:          b.Available,               //可用余额
				Locked:        b.Freeze,                  //冻结余额
				Borrowed:      b.TotalLiab,               //已借
				UpdateTime:    time.Now().UnixMilli(),
			})
		}
	case GATE_ASSET_TYPE_ISOLATED_MARGIN:
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).
			PrivateRestClient().NewPrivateRestMarginAccounts().Do()
		if err != nil {
			return nil, err
		}
		for _, d := range res.Data {
			assetBase := d.CurrencyPair + "/" + d.Base.Currency
			assetQuote := d.CurrencyPair + "/" + d.Quote.Currency

			baseBorrowed, _ := decimal.NewFromString(d.Base.Borrowed)
			baseInterest, _ := decimal.NewFromString(d.Base.Interest)

			quoteBorrowed, _ := decimal.NewFromString(d.Quote.Borrowed)
			quoteInterest, _ := decimal.NewFromString(d.Quote.Interest)

			baseFree, _ := decimal.NewFromString(d.Base.Available)
			quoteFree, _ := decimal.NewFromString(d.Quote.Available)

			baseLocked, _ := decimal.NewFromString(d.Base.Locked)
			quoteLocked, _ := decimal.NewFromString(d.Quote.Locked)

			//钱包余额=free+locked
			baseWalletBalance := baseFree.Add(baseLocked)
			quoteWalletBalance := quoteFree.Add(quoteLocked)

			//最大可转=free-borrowed-interest
			baseMaxWithdrawAmount := baseFree.Sub(baseBorrowed).Sub(baseInterest)
			quoteMaxWithdrawAmount := quoteFree.Sub(quoteBorrowed).Sub(quoteInterest)

			assets = append(assets, &Asset{
				Exchange:          a.ExchangeType().String(),      //交易所
				AccountType:       accountType,                    //账户类型
				Asset:             assetBase,                      //资产
				Borrowed:          d.Base.Borrowed,                //已借
				Interest:          d.Base.Interest,                //利息
				Free:              d.Base.Available,               //可用余额
				Locked:            d.Base.Locked,                  //冻结余额
				WalletBalance:     baseWalletBalance.String(),     //钱包余额
				MaxWithdrawAmount: baseMaxWithdrawAmount.String(), //最大可转出余额
				UpdateTime:        time.Now().UnixMilli(),
			})

			assets = append(assets, &Asset{
				Exchange:          a.ExchangeType().String(),       //交易所
				AccountType:       accountType,                     //账户类型
				Asset:             assetQuote,                      //资产
				Borrowed:          d.Quote.Borrowed,                //已借
				Interest:          d.Quote.Interest,                //利息
				Free:              d.Quote.Available,               //可用余额
				Locked:            d.Quote.Locked,                  //冻结余额
				WalletBalance:     quoteWalletBalance.String(),     //钱包余额
				MaxWithdrawAmount: quoteMaxWithdrawAmount.String(), //最大可转出余额
				UpdateTime:        time.Now().UnixMilli(),
			})
		}
	case GATE_ASSET_TYPE_FUTURES:
		settles := []string{"usdt", "btc"}
		var errG errgroup.Group
		for _, settle := range settles {
			settle := settle
			errG.Go(func() error {
				res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
					NewPrivateRestFuturesSettleAccounts().Settle(settle).Do()
				if err != nil {
					// log.Error(err)
					return nil
				}

				d := res.Data
				marginAvailable := true
				if d.MarginMode != 2 {
					marginAvailable = false
				}
				assets = append(assets, &Asset{
					Exchange:               a.ExchangeType().String(), //交易所
					AccountType:            accountType,               //账户类型
					Asset:                  d.Currency,                //资产
					WalletBalance:          d.Total,                   //余额
					UnrealizedProfit:       d.UnrealisedPnl,           //未实现盈亏
					MarginBalance:          d.Available,               //保证金余额
					MaintMargin:            d.MaintenanceMargin,       //维持保证金
					InitialMargin:          d.PositionInitialMargin,   //当前所需起始保证金
					PositionInitialMargin:  d.PositionInitialMargin,   //持仓所需起始保证金(基于最新标记价格)
					OpenOrderInitialMargin: d.PositionInitialMargin,   //当前挂单所需起始保证金(基于最新标记价格)
					MaxWithdrawAmount:      d.Available,               //最大可转出余额
					CrossWalletBalance:     d.CrossAvailable,          //全仓账户余额
					CrossUnPnl:             d.CrossUnrealisedPnl,      //全仓持仓未实现盈亏
					AvailableBalance:       d.Available,               //可用余额
					MarginAvailable:        marginAvailable,           //否可用作联合保证金
					UpdateTime:             time.Now().UnixMilli(),
				})

				return nil
			})
		}

		if err := errG.Wait(); err != nil {
			return nil, err
		}
	case GATE_ASSET_TYPE_DELIVERY:
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleAccounts().
			Settle("usdt").Do()
		if err != nil {
			return nil, err
		}

		marginAvailable := true
		if res.Data.MarginMode != 2 {
			marginAvailable = false
		}
		d := res.Data
		assets = append(assets, &Asset{
			Exchange:               a.ExchangeType().String(), //交易所
			AccountType:            accountType,               //账户类型
			Asset:                  d.Currency,                //资产
			WalletBalance:          d.Total,                   //余额
			UnrealizedProfit:       d.UnrealisedPnl,           //未实现盈亏
			MarginBalance:          d.Available,               //保证金余额
			MaintMargin:            d.MaintenanceMargin,       //维持保证金
			InitialMargin:          d.PositionInitialMargin,   //当前所需起始保证金
			PositionInitialMargin:  d.PositionInitialMargin,   //持仓所需起始保证金(基于最新标记价格)
			OpenOrderInitialMargin: d.PositionInitialMargin,   //当前挂单所需起始保证金(基于最新标记价格)
			MaxWithdrawAmount:      d.Available,               //最大可转出余额
			CrossWalletBalance:     d.CrossAvailable,          //全仓账户余额
			CrossUnPnl:             d.CrossUnrealisedPnl,      //全仓持仓未实现盈亏
			AvailableBalance:       d.Available,               //可用余额
			MarginAvailable:        marginAvailable,           //否可用作联合保证金
			UpdateTime:             time.Now().UnixMilli(),
		})
	}

	return assets, nil
}

// TODO IMPL
func (a GateTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	api := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().NewPrivateRestWalletTransfers().
		Currency(req.Asset).Amount(req.Amount.String())

	// required
	from := a.gateConverter.ToGateAssetType(req.From)
	to := a.gateConverter.ToGateAssetType(req.To)
	api.From(from).To(to)

	// 逐仓杠杆
	if from == GATE_ASSET_TYPE_ISOLATED_MARGIN {
		api.CurrencyPair(req.FromSymbol)
	} else if to == GATE_ASSET_TYPE_ISOLATED_MARGIN {
		api.CurrencyPair(req.ToSymbol)
	}

	// futures or delivery
	if from == GATE_ASSET_TYPE_FUTURES || to == GATE_ASSET_TYPE_FUTURES ||
		from == GATE_ASSET_TYPE_DELIVERY || to == GATE_ASSET_TYPE_DELIVERY {
		api.Settle(req.Settle)
	}

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	var assetTransfers []*AssetTransfer
	assetTransfers = append(assetTransfers, &AssetTransfer{
		Exchange:   a.ExchangeType().String(),
		TranId:     strconv.FormatInt(res.Data.TxId, 10),
		Asset:      req.Asset,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount.String(),
		Status:     "",
		ClientId:   "",
		FromSymbol: req.FromSymbol,
		ToSymbol:   req.ToSymbol,
	})
	return assetTransfers, nil
}

// TODO IMPL
func (a GateTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	//switch req.AssetType {
	//case ASSET_TYPE_FUND, ASSET_TYPE_MARGIN:
	//	api := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().NewPrivateRestSpotAccountBook()
	//	if req.Asset != "" {
	//		api.Currency(req.Asset)
	//	}
	//	if req.StartTime != 0 {
	//		api.From(req.StartTime)
	//	}
	//	if req.EndTime != 0 {
	//		api.To(req.EndTime)
	//	}
	//
	//	res, err := api.Do()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	var assetTransfers []*QueryAssetTransfer
	//	for _, d := range res.Data {
	//		var from, to AssetType
	//		split := strings.Split(d.Type, "_")
	//		if len(split) < 2 {
	//			return nil, errors.New("type error")
	//		}
	//		if split[1] == "in" { // 转入
	//			from, to = GATE_ACCOUNT_TYPE_UNKNOWN, a.gateConverter.FromGateAssetType(split[0])
	//		} else {
	//			from, to = a.gateConverter.FromGateAssetType(split[0]), GATE_ACCOUNT_TYPE_UNKNOWN
	//		}
	//		assetTransfers = append(assetTransfers, &QueryAssetTransfer{
	//			TranId:   "",
	//			ChangeId: d.ID,
	//			Asset:    d.Currency,
	//			From:     from,
	//			To:       to,
	//			Amount:   decimal.RequireFromString(d.Change),
	//			Status:   "",
	//		})
	//	}
	//
	//	return assetTransfers, nil
	//}
	return nil, ErrorNotSupport
}
