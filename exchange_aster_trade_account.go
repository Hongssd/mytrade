package mytrade

import (
	"strconv"
	"strings"
	"time"

	"github.com/Hongssd/myasterapi"
	"github.com/shopspring/decimal"
)

type AsterTradeAccount struct {
	ExchangeBase

	asterConverter AsterEnumConverter
	apiKey         string
	secretKey      string
}

func (b AsterTradeAccount) GetAccountMode() (AccountMode, error) {
	res, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureMultiAssetsMarginGet().Do()
	if err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	return b.asterConverter.FromAsterAccountMode(res.MultiAssetsMargin), nil
}

func (b AsterTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		return MARGIN_MODE_CROSSED, nil
	case ASTER_AC_FUTURE:
		res, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
		if err != nil {
			return MARGIN_MODE_UNKNOWN, err
		}
		for _, p := range res.Positions {
			if p.Symbol == symbol && p.PositionSide == b.asterConverter.ToAsterPositionSide(positionSide) {
				return b.asterConverter.FromAsterMarginMode(p.Isolated), nil
			}
		}
	default:
		return MARGIN_MODE_UNKNOWN, ErrorAccountType
	}
	return MARGIN_MODE_UNKNOWN, ErrorSymbolNotFound
}

func (b AsterTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {

	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		return POSITION_MODE_ONEWAY, nil
	case ASTER_AC_FUTURE:
		res, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFuturePositionSideDualGet().Do()
		if err != nil {
			return POSITION_MODE_UNKNOWN, err
		}
		return b.asterConverter.FromAsterPositionMode(res.DualSidePosition), nil
	default:
		return POSITION_MODE_UNKNOWN, ErrorNotSupport
	}

}

func (b AsterTradeAccount) GetLeverage(accountType, symbol string,
	marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	leverage := decimal.NewFromInt(0)

	if accountType == ASTER_AC_SPOT.String() {
		// get spot isolated leverage is not supported （现货逐仓杠杆没有设置和查询接口）
		return leverage, ErrorNotSupport
	}

	positionMode, err := b.GetPositionMode(accountType, symbol)
	if err != nil {
		return leverage, err
	}

	check := func(s string, i bool, ps string, l string) (decimal.Decimal, bool) {
		var result decimal.Decimal
		if s == symbol && i == b.asterConverter.ToAsterMarginMode(marginMode) {
			//根据保证金模式，仓位模式，仓位方向获取指定杠杆
			switch positionMode {
			case POSITION_MODE_ONEWAY:
				if ps == b.asterConverter.ToAsterPositionSide(POSITION_SIDE_BOTH) {
					result, _ = decimal.NewFromString(l)
					return result, true
				}
			case POSITION_MODE_HEDGE:
				if ps == b.asterConverter.ToAsterPositionSide(positionSide) {
					result, _ = decimal.NewFromString(l)
					return result, true
				}
			default:
				return result, false
			}
		}
		return result, false
	}

	switch AsterAccountType(accountType) {
	case ASTER_AC_FUTURE:
		res, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
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

func (b AsterTradeAccount) SetAccountMode(mode AccountMode) error {
	nowAccountMode, err := b.GetAccountMode()
	if err != nil {
		return err
	}
	if nowAccountMode == mode {
		return nil
	}
	_, err = aster.NewFutureRestClient(b.apiKey, b.secretKey).
		NewFutureMultiAssetsMarginPost().
		MultiAssetsMargin(b.asterConverter.ToAsterAccountMode(mode)).Do()
	if err != nil {
		return err
	}
	return nil
}

func (b AsterTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	if accountType == ASTER_AC_SPOT.String() {
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

	switch AsterAccountType(accountType) {
	case ASTER_AC_FUTURE:
		_, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFutureMarginType().Symbol(symbol).MarginType(b.asterConverter.ToAsterMarginModeStr(mode)).Do()
		if err != nil {
			return err
		}
	default:
		return ErrorAccountType
	}
	return nil
}

func (b AsterTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {

	nowPositionMode, err := b.GetPositionMode(accountType, symbol)
	if err != nil {
		return err
	}

	if nowPositionMode == mode {
		return nil
	}

	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		return ErrorNotSupport
	case ASTER_AC_FUTURE:
		_, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFuturePositionSideDualPost().DualSidePosition(b.asterConverter.ToAsterPositionMode(mode)).Do()
		if err != nil {
			return err
		}
	default:
		return ErrorAccountType
	}
	return nil
}

func (b AsterTradeAccount) SetLeverage(accountType, symbol string,
	marginMode MarginMode, positionSide PositionSide,
	leverage decimal.Decimal) error {
	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		// set spot isolated leverage is not supported
		return ErrorNotSupport
	case ASTER_AC_FUTURE:
		_, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFutureLeverage().Symbol(symbol).Leverage(leverage.IntPart()).Do()
		if err != nil {
			return err
		}
	default:
		return ErrorAccountType
	}
	return nil
}

func (b AsterTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	var feeRate FeeRate
	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		res, err := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAccount().Do()
		if err != nil {
			return nil, err
		}
		feeRate.Maker, _ = decimal.NewFromString(res.CommissionRates.Maker)
		feeRate.Taker, _ = decimal.NewFromString(res.CommissionRates.Taker)
	case ASTER_AC_FUTURE:
		res, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).
			NewFutureCommissionRate().Symbol(symbol).Do()
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

func (b AsterTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	var positionList []*Position
	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		return positionList, nil
	case ASTER_AC_FUTURE:

		var res *myasterapi.FutureAccountRes
		var risk *myasterapi.FuturePositionRiskRes
		err := ErrGroupWait(func() error {
			r, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
			if err != nil {
				return err
			}
			res = r
			return nil
		}, func() error {
			r, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).
				NewFuturePositionRisk().Do()
			if err != nil {
				return err
			}
			risk = r
			return nil
		})
		if err != nil {
			return nil, err
		}

		riskMap := map[string]myasterapi.FuturePositionRiskRow{}
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
					MarginMode:             b.asterConverter.FromAsterMarginMode(p.Isolated),
					EntryPrice:             p.EntryPrice,
					MaxNotional:            p.MaxNotional,
					PositionSide:           b.asterConverter.FromAsterPositionSide(p.PositionSide),
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
	default:
		return nil, ErrorAccountType
	}

	return positionList, nil
}

func (b AsterTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	var assetList []*Asset

	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		res, err := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAccount().Do()
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
	case ASTER_AC_FUTURE:
		res, err := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureAccount().Do()
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

func (b AsterTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAssetTransferPost()

	FromAsset := b.asterConverter.ToAsterAssetType(req.From)
	ToAsset := b.asterConverter.ToAsterAssetType(req.To)
	AsterTransferType := FromAsset + "_" + ToAsset
	api.Type(myasterapi.AssetTransferType(AsterTransferType))

	api.Asset(req.Asset).Amount(req.Amount)

	if req.FromSymbol != "" {
		api.FromSymbol(req.FromSymbol)
	}

	if req.ToSymbol != "" {
		api.ToSymbol(req.ToSymbol)
	}

	res, err := api.Do()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var assetTransfers []*AssetTransfer
	tranId := strconv.FormatInt(res.TranId, 10)
	assetTransfers = append(assetTransfers, &AssetTransfer{
		Exchange:   b.ExchangeType().String(),
		TranId:     tranId,
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

func (b AsterTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotAssetTransferGet()
	Type := b.asterConverter.ToAsterAssetType(req.From) + "_" + b.asterConverter.ToAsterAssetType(req.To)
	api.Type(myasterapi.AssetTransferType(Type))
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
			From:   b.asterConverter.FromAsterAssetType(strings.Split(r.Type, "_")[0]),
			To:     b.asterConverter.FromAsterAssetType(strings.Split(r.Type, "_")[1]),
			Status: b.asterConverter.FromAsterTransferStatus(r.Status),
		})
	}

	return QueryAssetTransfers, nil
}
