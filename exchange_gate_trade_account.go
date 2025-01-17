package mytrade

import (
	"errors"
	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
	"time"
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
		PrivateRestClient().
		NewPrivateRestAccountDetail().Do()
	if err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	return a.gateConverter.FromGateAccountMode(res.Data.Key.Mode), nil
}

func (a GateTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	return MARGIN_MODE_CROSSED, nil
}

func (a GateTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	switch a.gateConverter.ToGateAssetType(AssetType(accountType)) {
	case GATE_ACCOUNT_TYPE_SPOT, GATE_ACCOUNT_TYPE_MARGIN:
		return POSITION_MODE_ONEWAY, nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		split := strings.Split(symbol, "_")
		if len(split) != 2 {
			return POSITION_MODE_UNKNOWN, errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
			NewPrivateRestFuturesSettleAccounts().Settle(strings.ToLower(settle)).Do()
		if err != nil {
			return POSITION_MODE_UNKNOWN, err
		}
		if res.Data.InDualMode {
			return POSITION_MODE_HEDGE, nil
		} else {
			return POSITION_MODE_ONEWAY, nil
		}
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
		if res.Data.InDualMode {
			return POSITION_MODE_HEDGE, nil
		} else {
			return POSITION_MODE_ONEWAY, nil
		}
	}

	return POSITION_MODE_UNKNOWN, nil
}

func (a GateTradeAccount) GetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	if symbol == "" {
		return decimal.Zero, errors.New("symbol error")
	}
	switch a.gateConverter.ToGateAssetType(AssetType(accountType)) {
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
	return ErrorNotSupport
}

func (a GateTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	return ErrorNotSupport
}

func (a GateTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	// 仅支持永续合约
	log.Info(a.gateConverter.ToGateAssetType(AssetType(accountType)))
	switch a.gateConverter.ToGateAssetType(AssetType(accountType)) {
	case GATE_ACCOUNT_TYPE_FUTURES:
		split := strings.Split(symbol, "_")
		if len(split) < 2 {
			return errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		_, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
			NewPrivateRestFuturesSettleDualMode().Settle(settle).
			DualMode(mode == POSITION_MODE_HEDGE).Do()
		if err != nil {
			return err
		}
		return nil
	default:
		return ErrorNotSupport
	}
}

func (a GateTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode, positionMode PositionMode, positionSide PositionSide, leverage decimal.Decimal) error {
	switch a.gateConverter.ToGateAssetType(AssetType(accountType)) {
	case GATE_ACCOUNT_TYPE_FUTURES:
		split := strings.Split(symbol, "_")
		if len(split) < 2 {
			return errors.New("symbol error")
		}
		settle := strings.ToLower(split[1])
		switch positionMode {
		case POSITION_MODE_ONEWAY:
			api := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
				NewPrivateRestFuturesSettlePositionsContractLeverage().Settle(settle).Contract(symbol).Leverage(leverage.String())

			// 全仓模式下的杠杆倍数（即 leverage 为 0 时）
			if marginMode == MARGIN_MODE_CROSSED {
				api.Leverage("0")
				api.CrossLeverageLimit(leverage.String())
			}
			_, err := api.Do()
			if err != nil {
				return err
			}
			return nil

		case POSITION_MODE_HEDGE:
			api := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
				NewPrivateRestFuturesSettleDualCompPositionsContractLeverage().Settle(settle).Contract(symbol).Leverage(leverage.String())

			// 全仓模式下的杠杆倍数（即 leverage 为 0 时）
			if marginMode == MARGIN_MODE_CROSSED {
				api.CrossLeverageLimit(leverage.String())
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
	if accountType == GATE_ACCOUNT_TYPE_FUTURES || accountType == GATE_ACCOUNT_TYPE_DELIVERY {
		split := strings.Split(symbol, "_")
		settle := strings.ToLower(split[1])
		api.Settle(settle)
	}
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	switch accountType {
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
	switch accountType {
	case GATE_ACCOUNT_TYPE_FUTURES:
		var errG errgroup.Group
		settles := []string{"usdt", "btc"}
		for _, settle := range settles {
			settle := settle
			errG.Go(func() error {
				res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().NewPrivateRestFuturesSettlePositions().Settle(settle).Do()
				if err != nil {
					return err
				}
				for _, p := range res.Data {
					for _, symbol := range symbols {
						if p.Contract == symbol {
							var marginMode MarginMode
							var leverage string
							if p.Leverage == "0" {
								marginMode = MARGIN_MODE_CROSSED
								leverage = p.CrossLeverageLimit
							} else {
								marginMode = MARGIN_MODE_ISOLATED
								leverage = p.Leverage
							}

							var positionSide PositionSide
							if p.EntryPrice > p.LiqPrice {
								positionSide = POSITION_SIDE_LONG
							} else {
								positionSide = POSITION_SIDE_SHORT
							}

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
								UpdateTime:             p.UpdateTime,
							})
						}
					}
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
				for _, symbol := range symbols {
					if p.Contract == symbol {
						var marginMode MarginMode
						var leverage string
						if p.Leverage == "0" {
							marginMode = MARGIN_MODE_CROSSED
							leverage = p.CrossLeverageLimit
						} else {
							marginMode = MARGIN_MODE_ISOLATED
							leverage = p.Leverage
						}

						var positionSide PositionSide
						if p.EntryPrice > p.LiqPrice {
							positionSide = POSITION_SIDE_LONG
						} else {
							positionSide = POSITION_SIDE_SHORT
						}

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
							UpdateTime:             p.UpdateTime,
						})
					}
				}
			}
		}
	default:
		return nil, ErrorPositionNotFound
	}
	return positions, nil
}

func (a GateTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	var assets []*Asset

	// 现货资产
	switch a.gateConverter.ToGateAssetType(AssetType(accountType)) {
	case GATE_ACCOUNT_TYPE_SPOT:
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
				UpdateTime:             time.Now().Unix(),
			})
		}
	case GATE_ACCOUNT_TYPE_MARGIN:
		//TODO
		return nil, ErrorNotSupport
	case GATE_ACCOUNT_TYPE_FUTURES:
		settles := []string{"usdt", "btc"}
		var errG errgroup.Group
		for _, settle := range settles {
			settle := settle
			errG.Go(func() error {
				res, err := mygateapi.NewRestClient(a.apiKey, a.secretKey).PrivateRestClient().
					NewPrivateRestFuturesSettleAccounts().Settle(settle).Do()
				if err != nil {
					log.Error(err)
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
	case GATE_ACCOUNT_TYPE_DELIVERY:
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

	// margin
	if from == GATE_ACCOUNT_TYPE_MARGIN || to == GATE_ACCOUNT_TYPE_MARGIN {
		api.CurrencyPair(req.CurrencyPair)
	}

	// futures or delivery
	if from == GATE_ACCOUNT_TYPE_FUTURES || to == GATE_ACCOUNT_TYPE_FUTURES || from == GATE_ACCOUNT_TYPE_DELIVERY || to == GATE_ACCOUNT_TYPE_DELIVERY {
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
		FromSymbol: "",
		ToSymbol:   "",
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
