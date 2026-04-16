package mytrade

import (
	"strconv"
	"strings"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
	"github.com/shopspring/decimal"
)

type BitgetExchangeInfo struct {
	ExchangeBase
	isLoaded bool

	spotSymbolMap        *MySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]
	marginSymbolMap      *MySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]
	usdtFuturesSymbolMap *MySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]
	coinFuturesSymbolMap *MySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]
	usdcFuturesSymbolMap *MySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]

	spotExchangeInfoMap        *MySyncMap[string, TradeSymbolInfo]
	marginExchangeInfoMap      *MySyncMap[string, TradeSymbolInfo]
	usdtFuturesExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
	coinFuturesExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
	usdcFuturesExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
}

type bitgetSymbolInfo struct {
	symbolInfo
}

func (e *BitgetExchangeInfo) preCheckAccountType(accountType string) bool {
	switch accountType {
	case BITGET_AC_SPOT, BITGET_AC_MARGIN, BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
		return true
	default:
		return false
	}
}

func (e *BitgetExchangeInfo) loadExchangeInfo() error {
	pub := mybitgetapi.NewRestClient("", "", "").PublicRestClient()
	e.spotSymbolMap = GetPointer(NewMySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]())
	e.marginSymbolMap = GetPointer(NewMySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]())
	e.usdtFuturesSymbolMap = GetPointer(NewMySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]())
	e.coinFuturesSymbolMap = GetPointer(NewMySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]())
	e.usdcFuturesSymbolMap = GetPointer(NewMySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]())

	loadByCategory := func(category string, target *MySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow]) error {
		res, err := pub.NewPublicRestUtaMarketInstruments().Category(category).Do()
		if err != nil {
			return err
		}
		for _, row := range res.Data {
			r := row
			target.Store(row.Symbol, &r)
		}
		return nil
	}
	if err := loadByCategory(BITGET_AC_SPOT, e.spotSymbolMap); err != nil {
		return err
	}
	if err := loadByCategory(BITGET_AC_MARGIN, e.marginSymbolMap); err != nil {
		return err
	}
	if err := loadByCategory(BITGET_AC_USDT_FUTURES, e.usdtFuturesSymbolMap); err != nil {
		return err
	}
	if err := loadByCategory(BITGET_AC_COIN_FUTURES, e.coinFuturesSymbolMap); err != nil {
		return err
	}
	if err := loadByCategory(BITGET_AC_USDC_FUTURES, e.usdcFuturesSymbolMap); err != nil {
		return err
	}

	e.isLoaded = true
	e.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.marginExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.usdtFuturesExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.coinFuturesExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.usdcFuturesExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	return nil
}

func (e *BitgetExchangeInfo) Refresh() error {
	e.isLoaded = false
	return e.loadExchangeInfo()
}

func (e *BitgetExchangeInfo) getSymbolMapByAccountType(accountType string) (*MySyncMap[string, *mybitgetapi.PublicRestUtaMarketInstrumentsResRow], error) {
	switch accountType {
	case BITGET_AC_SPOT:
		return e.spotSymbolMap, nil
	case BITGET_AC_MARGIN:
		return e.marginSymbolMap, nil
	case BITGET_AC_USDT_FUTURES:
		return e.usdtFuturesSymbolMap, nil
	case BITGET_AC_COIN_FUTURES:
		return e.coinFuturesSymbolMap, nil
	case BITGET_AC_USDC_FUTURES:
		return e.usdcFuturesSymbolMap, nil
	default:
		return nil, ErrorAccountType
	}
}

func (e *BitgetExchangeInfo) getExchangeInfoMapByAccountType(accountType string) (*MySyncMap[string, TradeSymbolInfo], error) {
	switch accountType {
	case BITGET_AC_SPOT:
		return e.spotExchangeInfoMap, nil
	case BITGET_AC_MARGIN:
		return e.marginExchangeInfoMap, nil
	case BITGET_AC_USDT_FUTURES:
		return e.usdtFuturesExchangeInfoMap, nil
	case BITGET_AC_COIN_FUTURES:
		return e.coinFuturesExchangeInfoMap, nil
	case BITGET_AC_USDC_FUTURES:
		return e.usdcFuturesExchangeInfoMap, nil
	default:
		return nil, ErrorAccountType
	}
}

func (e *BitgetExchangeInfo) GetSymbolInfo(accountType, symbol string) (TradeSymbolInfo, error) {
	if !e.isLoaded {
		if err := e.loadExchangeInfo(); err != nil {
			return nil, err
		}
	}

	if !e.preCheckAccountType(accountType) {
		return nil, ErrorAccountType
	}
	infoMap, err := e.getExchangeInfoMapByAccountType(accountType)
	if err != nil {
		return nil, err
	}
	if v, ok := infoMap.Load(symbol); ok {
		return v, nil
	}

	symbolMap, err := e.getSymbolMapByAccountType(accountType)
	if err != nil {
		return nil, err
	}
	row, ok := symbolMap.Load(symbol)
	if !ok {
		return nil, ErrorSymbolNotFound
	}
	info := bitgetRowToSymbolInfo(row, accountType)
	infoMap.Store(symbol, info)
	return info, nil
}

func (e *BitgetExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {
	if !e.isLoaded {
		if err := e.Refresh(); err != nil {
			return nil, err
		}
	}

	if !e.preCheckAccountType(accountType) {
		return nil, ErrorAccountType
	}
	symbolMap, err := e.getSymbolMapByAccountType(accountType)
	if err != nil {
		return nil, err
	}
	var symbolInfoList []TradeSymbolInfo
	symbolMap.Range(func(sym string, _ *mybitgetapi.PublicRestUtaMarketInstrumentsResRow) bool {
		info, err := e.GetSymbolInfo(accountType, sym)
		if err != nil {
			return false
		}
		symbolInfoList = append(symbolInfoList, info)
		return true
	})
	return symbolInfoList, nil
}

func bitgetDecimalStepFromPrecision(precStr string) string {
	n, err := strconv.Atoi(strings.TrimSpace(precStr))
	if err != nil || n <= 0 {
		return "1"
	}
	if n == 1 {
		return "0.1"
	}
	return "0." + strings.Repeat("0", n-1) + "1"
}

func bitgetRowToSymbolInfo(row *mybitgetapi.PublicRestUtaMarketInstrumentsResRow, requestedAccountType string) *bitgetSymbolInfo {
	pricePrec, _ := strconv.Atoi(strings.TrimSpace(row.PricePrecision))
	qtyPrec, _ := strconv.Atoi(strings.TrimSpace(row.QuantityPrecision))
	tickSize := strings.TrimSpace(row.PriceMultiplier)
	if tickSize == "" {
		tickSize = bitgetDecimalStepFromPrecision(row.PricePrecision)
	}
	lotSize := strings.TrimSpace(row.QuantityMultiplier)
	if lotSize == "" {
		lotSize = bitgetDecimalStepFromPrecision(row.QuantityPrecision)
	}

	isContract := strings.Contains(row.Category, "FUTURES")
	isTrading := row.Status == "online"

	minAmt := row.MinOrderQty
	minNotional := row.MinOrderAmount
	if minAmt == "" {
		minAmt = "0"
	}
	if minNotional == "" {
		minNotional = "0"
	}

	maxLmt := row.MaxOrderQty
	maxMkt := row.MaxMarketOrderQty
	if maxLmt == "" {
		maxLmt = "0"
	}
	if maxMkt == "" {
		maxMkt = "0"
	}

	maxOrderNum, _ := strconv.Atoi(strings.TrimSpace(row.MaxSymbolOrderNum))

	contractSize := "1"
	contractCoin := row.BaseCoin
	if !isContract {
		contractSize = "0"
		contractCoin = ""
	}

	var maxLeverage, minLeverage string
	if row.Category == BITGET_AC_MARGIN {
		crossedLeverage := decimal.RequireFromString(row.MaxCrossedLeverage)
		isolatedLeverage := decimal.RequireFromString(row.MaxIsolatedLeverage)
		if crossedLeverage.GreaterThan(isolatedLeverage) {
			maxLeverage = isolatedLeverage.String()
			minLeverage = "0"
		} else {
			maxLeverage = crossedLeverage.String()
			minLeverage = "0"
		}
	} else {
		maxLeverage = row.MaxLeverage
		minLeverage = row.MinLeverage
	}

	return &bitgetSymbolInfo{symbolInfo: symbolInfo{symbolInfoStruct: symbolInfoStruct{
		Exchange:       BITGET_NAME.String(),
		AccountType:    requestedAccountType,
		Symbol:         row.Symbol,
		BaseCoin:       row.BaseCoin,
		QuoteCoin:      row.QuoteCoin,
		IsTrading:      isTrading,
		IsContract:     isContract,
		IsContractAmt:  false,
		ContractSize:   contractSize,
		ContractCoin:   contractCoin,
		ContractType:   row.Type,
		PricePrecision: pricePrec,
		AmtPrecision:   qtyPrec,
		TickSize:       tickSize,
		MinPrice:       "0",
		MaxPrice:       "0",
		LotSize:        lotSize,
		MinAmt:         minAmt,
		MaxLmtAmt:      maxLmt,
		MaxMktAmt:      maxMkt,
		MaxLeverage:    maxLeverage,
		MinLeverage:    minLeverage,
		StepLeverage:   "1",
		MaxOrderNum:    maxOrderNum,
		MinNotional:    minNotional,
	}}}
}
