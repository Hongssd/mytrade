package mytrade

import (
	"github.com/shopspring/decimal"
)

type SunxTradeAccount struct {
	ExchangeBase

	sunxConverter SunxEnumConverter
	accessKey     string
	secretKey     string
}

// 默认统一账户模式
func (s SunxTradeAccount) GetAccountMode() (AccountMode, error) {
	return ACCOUNT_MODE_UNIFIED, nil
}

// 默认全仓保证金模式
func (s SunxTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	return MARGIN_MODE_CROSSED, nil
}

// 默认单向持仓模式
func (s SunxTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	res, err := sunx.NewPrivateRestClient(s.accessKey, s.secretKey).
		NewPrivateRestPositionModeGet().Do()
	if err != nil {
		return POSITION_MODE_UNKNOWN, err
	}

	return s.sunxConverter.FromSunxPositionMode(res.Data.PositionMode), nil
}

// 获取杠杆
func (s SunxTradeAccount) GetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	res, err := sunx.NewPrivateRestClient(s.accessKey, s.secretKey).
		NewPrivateRestPositionLeverGet().
		MarginMode(SUNX_MARGIN_MODE_CROSSED).
		ContractCode(symbol).
		Do()
	if err != nil {
		return decimal.Zero, err
	}
	if len(res.Data) == 0 {
		return decimal.Zero, ErrorPositionNotFound
	}
	return decimal.NewFromInt(res.Data[0].LeverRate), nil
}

// sunx不需要设置账户模式，仅支持统一账户模式
func (s SunxTradeAccount) SetAccountMode(mode AccountMode) error {
	return nil
}

// sunx不需要设置逐仓/全仓，仅支持全仓保证金模式
func (s SunxTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	return nil
}

func (s SunxTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	_, err := sunx.NewPrivateRestClient(s.accessKey, s.secretKey).
		NewPrivateRestPositionModePost().
		PositionMode(s.sunxConverter.ToSunxPositionMode(mode)).
		Do()
	if err != nil {
		return err
	}
	return nil
}

func (s SunxTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {
	_, err := sunx.NewPrivateRestClient(s.accessKey, s.secretKey).
		NewPrivateRestPositionLeverPost().ContractCode(symbol).
		MarginMode(s.sunxConverter.ToSunxMarginMode(marginMode)).
		LeverRate(leverage.String()).
		Do()

	if err != nil {
		return err
	}
	return nil
}

func (s SunxTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	var feeRate FeeRate
	res, err := sunx.NewPrivateRestClient(s.accessKey, s.secretKey).
		NewPrivateRestAccountFeeRate().ContractCode(symbol).
		Do()
	if err != nil {
		return nil, err
	}
	if len(res.Data) == 0 {
		return nil, ErrorSymbolNotFound
	}
	feeRate.Maker, _ = decimal.NewFromString(res.Data[0].OpenMakerFee)
	feeRate.Taker, _ = decimal.NewFromString(res.Data[0].OpenTakerFee)
	return &feeRate, nil
}

func (s SunxTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	symbolMap := make(map[string]string)
	for _, symbol := range symbols {
		symbolMap[symbol] = symbol
	}
	res, err := sunx.NewPrivateRestClient(s.accessKey, s.secretKey).
		NewPrivateRestTradePositionOpens().Do()
	if err != nil {
		return nil, err
	}
	var positions []*Position
	for _, p := range res.Data {
		if _, ok := symbolMap[p.ContractCode]; ok {
			positions = append(positions, &Position{
				Exchange:               s.ExchangeType().String(),
				AccountType:            accountType,
				Symbol:                 p.ContractCode,
				MarginCcy:              p.MarginCurrency,
				InitialMargin:          p.InitialMargin,
				MaintMargin:            p.MaintenanceMargin,
				UnrealizedProfit:       p.ProfitUnreal,
				PositionInitialMargin:  p.InitialMargin,
				OpenOrderInitialMargin: "0",
				Leverage:               decimal.NewFromInt(p.LeverRate).String(),
				MarginMode:             s.sunxConverter.FromSunxMarginMode(p.MarginMode),
				EntryPrice:             p.OpenAvgPrice,
				MaxNotional:            "0",
				PositionSide:           s.sunxConverter.FromSunxPositionSide(p.PositionSide),
				PositionAmt:            p.Volume,
				MarkPrice:              p.MarkPrice,
				LiquidationPrice:       p.LiquidationPrice,
				MarginRatio:            p.MarginRate,
				UpdateTime:             decimal.RequireFromString(p.UpdatedTime).IntPart(),
			})
		}
	}
	return positions, nil
}

func (s SunxTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	currencyMap := make(map[string]string)
	for _, currency := range currencies {
		currencyMap[currency] = currency
	}
	res, err := sunx.NewPrivateRestClient(s.accessKey, s.secretKey).
		NewPrivateRestAccountBalance().Do()
	if err != nil {
		return nil, err
	}
	var assets []*Asset
	for _, d := range res.Data.Details {
		if len(currencyMap) > 0 {
			if _, ok := currencyMap[d.Currency]; !ok {
				continue
			}
		}

		equity, _ := decimal.NewFromString(d.Equity)
		initialMargin, _ := decimal.NewFromString(d.InitialMargin)
		// profitUnreal, _ := decimal.NewFromString(d.ProfitUnreal)

		// 冻结 = 初始保证金
		locked := initialMargin
		// 可用 = 权益 - 初始保证金
		free := equity.Sub(locked)

		assets = append(assets, &Asset{
			Exchange:               s.ExchangeType().String(),
			AccountType:            SUNX_ACCOUNT_TYPE_SWAP.String(),
			Asset:                  d.Currency,
			Free:                   free.String(),
			Locked:                 locked.String(),
			WalletBalance:          d.Available,
			UnrealizedProfit:       d.ProfitUnreal,
			MarginBalance:          d.Equity,
			MaintMargin:            d.MaintenanceMargin,
			InitialMargin:          d.InitialMargin,
			PositionInitialMargin:  d.InitialMargin,
			OpenOrderInitialMargin: "0",
			CrossWalletBalance:     d.Available,
			CrossUnPnl:             d.ProfitUnreal,
			AvailableBalance:       free.String(),
			MaxWithdrawAmount:      d.Available,
			MarginAvailable:        true,
			UpdateTime:             res.Data.UpdatedTime,
		})
	}
	return assets, nil
}

func (s SunxTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	return nil, ErrorNotSupport
}

func (s SunxTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	return nil, ErrorNotSupport
}
