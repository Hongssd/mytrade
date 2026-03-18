package mytrade

import (
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type XcoinTradeAccount struct {
	ExchangeBase

	xcoinConverter XcoinEnumConverter
	apiKey         string
	apiSecret      string
}

func (x XcoinTradeAccount) GetAccountMode() (AccountMode, error) {
	return ACCOUNT_MODE_UNIFIED, nil
}

func (x XcoinTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	return MARGIN_MODE_CROSSED, nil
}

func (x XcoinTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	return POSITION_MODE_ONEWAY, nil
}

func (x XcoinTradeAccount) GetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	if accountType == XCOIN_ACCOUNT_TYPE_SPOT.String() {
		return decimal.NewFromInt(1), nil
	}
	res, err := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().
		NewPrivateRestTradeLeverGet().BusinessType(accountType).Symbol(symbol).Do()
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.RequireFromString(res.Data.Lever), nil
}

func (x XcoinTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	return nil, ErrorNotSupport
}

func (x XcoinTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	switch XcoinAccountType(accountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT:
		return []*Position{}, nil
	case XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL, XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
	default:
		return nil, ErrorAccountType
	}

	api := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().
		NewPrivateRestTradePosition().BusinessType(accountType)
	if len(symbols) == 1 {
		api.Symbol(symbols[0])
	}

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	var positions []*Position
	for _, v := range res.Data {
		if len(symbols) > 1 && !stringInSlice(v.Symbol, symbols) {
			continue
		}
		liquidationPrice := v.LiquidationPrice
		if liquidationPrice == "" {
			liquidationPrice = "0"
		}
		positionSide := POSITION_SIDE_BOTH
		if stringToDecimal(v.PositionQty).GreaterThan(decimal.Zero) {
			positionSide = POSITION_SIDE_LONG
		} else if stringToDecimal(v.PositionQty).LessThan(decimal.Zero) {
			positionSide = POSITION_SIDE_SHORT
		}
		positions = append(positions, &Position{
			Exchange:               x.ExchangeType().String(),
			AccountType:            accountType,
			Symbol:                 v.Symbol,
			MarginCcy:              "",
			InitialMargin:          v.Im,
			MaintMargin:            "",
			UnrealizedProfit:       v.Upl,
			PositionInitialMargin:  v.Im,
			OpenOrderInitialMargin: "0",
			Leverage:               v.Lever,
			MarginMode:             MARGIN_MODE_CROSSED,
			EntryPrice:             v.AvgPrice,
			MaxNotional:            "0",
			PositionSide:           positionSide,
			PositionAmt:            v.PositionQty,
			MarkPrice:              v.MarkPrice,
			LiquidationPrice:       liquidationPrice,
			MarginRatio:            "0",
			UpdateTime:             stringToInt64(v.UpdateTime),
		})
	}
	return positions, nil
}

func (x XcoinTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	switch XcoinAccountType(accountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT, XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL, XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
	default:
		return nil, ErrorAccountType
	}

	currencyList := strings.Join(currencies, ",")

	// 从钱包获取资产
	balanceApi := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().
		NewPrivateRestAccountBalance()
	if currencyList != "" {
		balanceApi.CurrencyList(currencyList)
	}
	balanceRes, err := balanceApi.Do()
	if err != nil {
		return nil, err
	}

	// 获取 MaxTransfer
	transferApi := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().
		NewPrivateRestAccountTransferBalance()
	if currencyList != "" {
		transferApi.CurrencyList(currencyList)
	}
	transferRes, err := transferApi.Do()
	if err != nil {
		return nil, err
	}

	currencyFilter := make(map[string]struct{}, len(currencies))
	for _, c := range currencies {
		currencyFilter[c] = struct{}{}
	}

	maxTransferMap := make(map[string]string, len(transferRes.Data))
	for _, v := range transferRes.Data {
		maxTransferMap[v.Currency] = v.MaxTransfer
	}

	var assets []*Asset
	updateTime := stringToInt64(balanceRes.Ts)
	if updateTime == 0 {
		updateTime = time.Now().UnixMilli()
	}
	marginAvailable := XcoinAccountType(accountType) != XCOIN_ACCOUNT_TYPE_SPOT
	for _, d := range balanceRes.Data.Details {
		if len(currencyFilter) > 0 {
			if _, ok := currencyFilter[d.Currency]; !ok {
				continue
			}
		}

		walletBalance := stringToDecimal(d.TotalBalance)
		borrowed := stringToDecimal(d.RealLiability)
		interest := stringToDecimal(d.AccruedInterest)
		maxWithdrawAmount := d.TotalBalance
		if v, ok := maxTransferMap[d.Currency]; ok && v != "" {
			maxWithdrawAmount = v
		}
		maxTransfer := stringToDecimal(maxWithdrawAmount)

		free := maxTransfer
		locked := walletBalance.Sub(free)
		if free.LessThan(decimal.Zero) {
			free = decimal.Zero
		}
		if locked.LessThan(decimal.Zero) {
			locked = decimal.Zero
		}

		initialMargin := stringToDecimal(d.InitialMargin)
		orderInitialMargin := stringToDecimal(d.OrderInitialMargin)
		positionInitialMargin := stringToDecimal(d.PositionInitialMargin)
		maintMargin := initialMargin.Sub(orderInitialMargin)
		if maintMargin.LessThan(decimal.Zero) {
			maintMargin = decimal.Zero
		}

		assets = append(assets, &Asset{
			Exchange:               x.ExchangeType().String(),
			AccountType:            accountType,
			Asset:                  d.Currency,
			Borrowed:               borrowed.String(),
			Interest:               interest.String(),
			Free:                   free.String(),
			Locked:                 locked.String(),
			WalletBalance:          walletBalance.String(),
			UnrealizedProfit:       d.Upl,
			MarginBalance:          d.Equity,
			MaintMargin:            maintMargin.String(),
			InitialMargin:          initialMargin.String(),
			PositionInitialMargin:  positionInitialMargin.String(),
			OpenOrderInitialMargin: d.OrderInitialMargin,
			CrossWalletBalance:     d.CashBalance,
			CrossUnPnl:             d.Upl,
			AvailableBalance:       free.String(),
			MaxWithdrawAmount:      maxWithdrawAmount,
			MarginAvailable:        marginAvailable,
			UpdateTime:             updateTime,
		})
	}

	return assets, nil
}

func (x XcoinTradeAccount) SetAccountMode(mode AccountMode) error {
	return ErrorNotSupport
}

func (x XcoinTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	return ErrorNotSupport
}

func (x XcoinTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	return ErrorNotSupport
}

func (x XcoinTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {
	switch XcoinAccountType(accountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT:
		return ErrorNotSupport
	case XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL, XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
		_, err := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().
			NewPrivateRestTradeLeverPost().
			BusinessType(accountType).
			Symbol(symbol).
			Lever(leverage.String()).Do()
		if err != nil {
			return err
		}
		return nil
	}
	return ErrorNotSupport
}

func (x XcoinTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	api := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().NewPrivateRestAssetTransfer()

	api.Currency(req.Asset).Amount(req.Amount.String())
	api.FromAccountType(x.xcoinConverter.ToXcoinAssetType(req.From))
	api.ToAccountType(x.xcoinConverter.ToXcoinAssetType(req.To))
	api.LoanTrans(req.LoanTrans)

	clientId := GetInstanceId("XCOIN")
	api.ClientTransferId(clientId)

	_, err := api.Do()
	if err != nil {
		return nil, err
	}

	var assetTransfers []*AssetTransfer
	assetTransfers = append(assetTransfers, &AssetTransfer{
		Exchange: x.ExchangeType().String(),
		Asset:    req.Asset,
		From:     req.From,
		To:       req.To,
		Amount:   req.Amount.String(),
		Status:   "",
		ClientId: clientId,
	})

	return assetTransfers, nil
}

func (x XcoinTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	api := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().NewPrivateRestAssetTransferHistory()

	api.Currency(req.Asset)

	if req.StartTime != 0 {
		api.BeginTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	var QueryAssetTransfers []*QueryAssetTransfer
	for _, v := range res.Data {
		QueryAssetTransfers = append(QueryAssetTransfers, &QueryAssetTransfer{
			TranId: v.Id,
			Asset:  v.Currency,
			Amount: stringToDecimal(v.Amount).Abs(),
			From:   x.xcoinConverter.FromXcoinAssetType(v.FromAccountType),
			To:     x.xcoinConverter.FromXcoinAssetType(v.ToAccountType),
			Status: x.xcoinConverter.FromXcoinTransferStatus(v.Status),
		})
	}
	return QueryAssetTransfers, nil
}
