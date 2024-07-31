package mytrade

import (
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

type OkxTradeAccount struct {
	exchangeBase

	okxConverter OkxEnumConverter
	apiKey       string
	secretKey    string
	passphrase   string
}

func (o OkxTradeAccount) GetAccountMode() (AccountMode, error) {
	res, err := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).
		PrivateRestClient().NewPrivateRestAccountConfig().Do()
	if err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	return o.okxConverter.FromOKXAccountMode(res.Data[0].AcctLv), nil
}

func (o OkxTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	return MARGIN_MODE_CROSSED, nil
}

func (o OkxTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	res, err := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).
		PrivateRestClient().NewPrivateRestAccountConfig().Do()
	if err != nil {
		return POSITION_MODE_UNKNOWN, err
	}
	return o.okxConverter.FromOKXPositionMode(res.Data[0].PosMode), nil
}

func (o OkxTradeAccount) GetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	res, err := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).
		PrivateRestClient().NewPrivateRestAccountLeverageInfo().
		InstId(symbol).MgnMode(o.okxConverter.ToOKXMarginMode(marginMode)).Do()
	if err != nil {
		return decimal.Zero, err
	}
	var leverage decimal.Decimal
	for _, d := range res.Data {
		if d.InstId == symbol && o.okxConverter.FromOKXPositionSide(d.PosSide) == positionSide {
			leverage, _ = decimal.NewFromString(d.Lever)
			break
		}
	}
	if leverage.IsZero() {
		return decimal.Zero, ErrorPositionNotFound
	}
	return leverage, nil
}

func (o OkxTradeAccount) SetAccountMode(mode AccountMode) error {
	nowAccountMode, err := o.GetAccountMode()
	if err != nil {
		return err
	}
	if nowAccountMode == mode {
		return nil
	}
	_, err = okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountSetAccountLevel().AcctLv(o.okxConverter.ToOKXAccountMode(mode)).Do()
	if err != nil {
		return err
	}
	return nil
}

// okx不需要设置逐仓/全仓，可同时支持两种模式开仓
func (o OkxTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	return nil
}

func (o OkxTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	nowPositionMode, err := o.GetPositionMode(accountType, symbol)
	if err != nil {
		return err
	}
	if nowPositionMode == mode {
		return nil
	}
	_, err = okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountSetPositionMode().PosMode(o.okxConverter.ToOKXPositionMode(mode)).Do()
	if err != nil {
		return err
	}
	return nil
}

func (o OkxTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {
	_, err := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountSetLeverage().InstId(symbol).
		Lever(leverage.String()).MgnMode(o.okxConverter.ToOKXMarginMode(marginMode)).
		PosSide(o.okxConverter.ToOKXPositionSide(positionSide)).Do()
	if err != nil {
		return err
	}
	return nil
}

func (o OkxTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	var feeRate FeeRate
	api := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountTradeFee().InstType(accountType)
	if OkxAccountType(accountType) == OKX_AC_SPOT {
		api = api.InstId(symbol)
	}
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	feeRate.Maker, _ = decimal.NewFromString(res.Data[0].Maker)
	feeRate.Taker, _ = decimal.NewFromString(res.Data[0].Taker)

	//okx正常费率为负数,需要乘以-1
	feeRate.Maker = feeRate.Maker.Mul(decimal.NewFromInt(-1))
	feeRate.Taker = feeRate.Taker.Mul(decimal.NewFromInt(-1))

	return &feeRate, nil
}

func (o OkxTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {

	var positions []*Position

	if OkxAccountType(accountType) == OKX_AC_SPOT {
		return positions, nil
	}

	api := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountPosition().InstType(accountType)
	if len(symbols) == 1 {
		api = api.InstId(symbols[1])
	} else if len(symbols) > 1 && len(symbols) < 10 {
		instIds := strings.Join(symbols, ",")
		api = api.InstId(instIds)
	}
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	for _, d := range res.Data {
		if len(symbols) == 0 || stringInSlice(d.InstId, symbols) {
			updateTime, _ := strconv.ParseInt(d.UTime, 10, 64)
			if d.LiqPx == "" {
				d.LiqPx = "0"
			}
			positions = append(positions, &Position{
				Exchange:               o.ExchangeType().String(),
				AccountType:            d.InstType,
				Symbol:                 d.InstId,
				InitialMargin:          d.Imr,
				MaintMargin:            d.Mmr,
				UnrealizedProfit:       d.Upl,
				PositionInitialMargin:  d.Imr,
				OpenOrderInitialMargin: "0",
				Leverage:               "0",
				MarginMode:             o.okxConverter.FromOKXMarginMode(d.MgnMode),
				EntryPrice:             d.AvgPx,
				MaxNotional:            "0",
				PositionSide:           o.okxConverter.FromOKXPositionSide(d.PosSide),
				PositionAmt:            d.Pos,
				MarkPrice:              d.IdxPx,
				LiquidationPrice:       d.LiqPx,
				MarginRatio:            d.MgnRatio,
				UpdateTime:             updateTime,
			})
		}
	}
	return positions, nil
}

func (o OkxTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	api := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountBalance()
	if len(currencies) == 1 {
		api = api.Ccy(currencies[0])
	} else if len(currencies) > 1 && len(currencies) < 20 {
		ccys := strings.Join(currencies, ",")
		api = api.Ccy(ccys)
	}
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	var assets []*Asset
	for _, d := range res.Data[0].Details {
		if len(currencies) == 0 || stringInSlice(d.Ccy, currencies) {
			updateTime, _ := strconv.ParseInt(d.UTime, 10, 64)
			ordFronzen, _ := decimal.NewFromString(d.OrdFrozen) //挂单冻结保证金
			frozenBal, _ := decimal.NewFromString(d.FrozenBal)  //仓位占用保证金+挂单冻结保证金
			//仓位维持保证金= d.FrozenBal - d.OrdFrozen
			MaintMargin := frozenBal.Sub(ordFronzen)
			assets = append(assets, &Asset{
				Exchange:               o.ExchangeType().String(),
				AccountType:            accountType,
				Asset:                  d.Ccy,
				Free:                   d.AvailBal,
				Locked:                 d.OrdFrozen,
				WalletBalance:          d.CashBal,            //钱包余额=币种余额
				UnrealizedProfit:       d.Upl,                //未实现盈亏
				MarginBalance:          d.Eq,                 //保证金余额=钱包余额+未实现盈亏
				MaintMargin:            MaintMargin.String(), //维持保证金=仓位占用保证金
				InitialMargin:          d.FrozenBal,          //当前所需起始保证金=仓位占用保证金+挂单冻结保证金
				PositionInitialMargin:  MaintMargin.String(), //持仓所需起始保证金=仓位占用保证金
				OpenOrderInitialMargin: d.OrdFrozen,          //挂单所需起始保证金
				CrossWalletBalance:     d.CashBal,            //全仓账户余额
				CrossUnPnl:             d.Upl,                //全仓持仓未实现盈亏
				AvailableBalance:       d.AvailEq,            //可用余额=钱包余额
				MaxWithdrawAmount:      d.AvailBal,           //最大可转出余额=币种余额
				MarginAvailable:        true,
				UpdateTime:             updateTime,
			})
		}
	}

	return assets, nil
}

// 资金划转（账户内）
func (o OkxTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	api := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAssetTransfer()

	// required
	api.Type("0").Ccy(req.Asset).Amt(req.Amount.String())
	api.From(o.okxConverter.ToOKXAssetType(req.From))
	api.To(o.okxConverter.ToOKXAssetType(req.To))

	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	log.Info(res)
	var assetTransfers []*AssetTransfer
	for _, d := range res.Data {
		assetTransfers = append(assetTransfers, &AssetTransfer{
			Exchange: o.ExchangeType().String(),
			TranId:   d.TransId,
			Asset:    d.Ccy,
			From:     o.okxConverter.FromOKXAssetType(d.From),
			To:       o.okxConverter.FromOKXAssetType(d.To),
			Amount:   d.Amt,
			Status:   "",
		})
	}

	return assetTransfers, nil
}
