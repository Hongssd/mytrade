package mytrade

import (
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

type OkxTradeAccount struct {
	ExchangeBase

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
	if accountType == OKX_AC_SPOT.String() {
		return POSITION_MODE_ONEWAY, nil
	}
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
		if accountType != OKX_AC_SPOT.String() {
			if d.InstId == symbol && o.okxConverter.FromOKXPositionSide(d.PosSide) == positionSide {
				leverage, _ = decimal.NewFromString(d.Lever)
				break
			}
		} else {
			if marginMode == MARGIN_MODE_ISOLATED || len(res.Data) == 1 {
				if d.InstId == symbol {
					leverage, _ = decimal.NewFromString(d.Lever)
					break
				}
			} else {
				spilit := strings.Split(symbol, "-")
				if len(spilit) != 2 {
					return decimal.Zero, ErrorPositionNotFound
				}
				switch positionSide {
				case POSITION_SIDE_LONG, POSITION_SIDE_BOTH:
					if d.InstId == spilit[1] {
						leverage, _ = decimal.NewFromString(d.Lever)
						break
					}
				case POSITION_SIDE_SHORT:
					if d.InstId == spilit[0] {
						leverage, _ = decimal.NewFromString(d.Lever)
						break
					}
				default:
				}
			}
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

	if accountType == OKX_AC_SPOT.String() {
		return ErrorNotSupport
	}

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

func (o OkxTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode,  positionSide PositionSide, leverage decimal.Decimal) error {

	api := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountSetLeverage().InstId(symbol).
		Lever(leverage.String()).MgnMode(o.okxConverter.ToOKXMarginMode(marginMode))

	if accountType == OKX_AC_SPOT.String() && marginMode == MARGIN_MODE_CROSSED {
		spilit := strings.Split(symbol, "-")
		if len(spilit) != 2 {
			return ErrorPositionNotFound
		}
		switch positionSide {
		case POSITION_SIDE_LONG, POSITION_SIDE_BOTH:
			api.Ccy(spilit[1])
		case POSITION_SIDE_SHORT:
			api.Ccy(spilit[0])
		}
	} else {
		api.PosSide(o.okxConverter.ToOKXPositionSide(positionSide))
	}

	_, err := api.Do()
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

	api := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().
		NewPrivateRestAccountPosition()

	if accountType == OKX_AC_SPOT.String() {
		api.InstType(OKX_AC_MARGIN.String())
	} else {
		api.InstType(accountType)
	}

	if len(symbols) == 1 {
		api = api.InstId(symbols[0])
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
				MarginCcy:              d.Ccy,
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
	var assets []*Asset
	// 获取资金账户余额
	if accountType == OKX_AC_FUNDING.String() {
		assetBalanceApi := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient().NewPrivateRestAssetBalances()
		res, err := assetBalanceApi.Do()
		if err != nil {
			return nil, err
		}
		for _, d := range res.Data {
			updateTime := time.Now().Unix()
			assets = append(assets, &Asset{
				Exchange:               o.ExchangeType().String(),
				AccountType:            accountType,
				Asset:                  d.Ccy,
				Free:                   d.AvailBal,
				Locked:                 d.FrozenBal,
				WalletBalance:          d.Bal,      //钱包余额=币种余额
				UnrealizedProfit:       "0",        //未实现盈亏
				MarginBalance:          "0",        //保证金余额=钱包余额+未实现盈亏
				MaintMargin:            "0",        //维持保证金=仓位占用保证金
				InitialMargin:          "0",        //当前所需起始保证金=仓位占用保证金+挂单冻结保证金
				PositionInitialMargin:  "0",        //持仓所需起始保证金=仓位占用保证金
				OpenOrderInitialMargin: "0",        //挂单所需起始保证金
				CrossWalletBalance:     "0",        //全仓账户余额
				CrossUnPnl:             "0",        //全仓持仓未实现盈亏
				AvailableBalance:       d.AvailBal, //可用余额=钱包余额
				MaxWithdrawAmount:      d.AvailBal, //最大可转出余额=币种余额
				MarginAvailable:        true,
				UpdateTime:             updateTime,
			})
		}
	} else {
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
					Borrowed:               d.Liab,               //币种负债额 值为正数，如 21625.64 适用于跨币种保证金模式/组合保证金模式
					Interest:               d.Interest,           //计息，应扣未扣利息 值为正数，如 9.01 适用于跨币种保证金模式/组合保证金模式
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

	// 生成自定义id，以便查询账单（查询划转历史）时使用
	clientId := GetInstanceId("OKX")
	api.ClientId(clientId)

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

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
			ClientId: d.ClientId,
		})
	}

	return assetTransfers, nil
}

func (o OkxTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	api := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()

	// 从资金流水获取划转记录
	var bills = []*OKXAssetBill{}
	// 130：从交易账户转入账单
	if o.okxConverter.ToOKXAssetType(req.From) == OKX_ASSET_TYPE_UNIFIED && o.okxConverter.ToOKXAssetType(req.To) == OKX_ASSET_TYPE_FUND {
		res, err := api.NewPrivateRestAssetBills().Type("130").
			Before(strconv.FormatInt(req.StartTime, 10)).
			After(strconv.FormatInt(req.EndTime, 10)).
			Do()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, d := range res.Data {
			bills = append(bills, &OKXAssetBill{
				BillId: d.BillId,
				Ccy:    d.Ccy,
				BalChg: d.BalChg,
				Bal:    d.Bal,
			})
		}
	} else if o.okxConverter.ToOKXAssetType(req.From) == OKX_ASSET_TYPE_FUND && o.okxConverter.ToOKXAssetType(req.To) == OKX_ASSET_TYPE_UNIFIED {
		// 131：转出至交易账户账单
		res, err := api.NewPrivateRestAssetBills().Type("131").
			Before(strconv.FormatInt(req.StartTime, 10)).
			After(strconv.FormatInt(req.EndTime, 10)).
			Do()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, d := range res.Data {
			bills = append(bills, &OKXAssetBill{
				BillId: d.BillId,
				Ccy:    d.Ccy,
				BalChg: d.BalChg,
				Bal:    d.Bal,
			})
		}
	}

	var QueryAssetTransfers []*QueryAssetTransfer
	for _, bill := range bills {
		if req.Asset != "" && bill.Ccy != req.Asset {
			continue
		}
		QueryAssetTransfers = append(QueryAssetTransfers, &QueryAssetTransfer{
			TranId: bill.BillId,
			Asset:  bill.Ccy,
			Amount: stringToDecimal(bill.BalChg).Abs(),
			From:   req.From,
			To:     req.To,
			Status: o.okxConverter.FromOKXTransferStatus(OKX_TRANSFER_STATUS_TYPE_SUCCESS),
		})
	}
	return QueryAssetTransfers, nil
}
