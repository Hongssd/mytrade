package mytrade

import (
	"github.com/Hongssd/myxcoinapi"
)

type XcoinExchangeInfo struct {
	ExchangeBase
	isLoaded                       bool
	spotSymbolMap                  *MySyncMap[string, *myxcoinapi.PublicRestPublicSymbolsResRow]
	linearPerpetualSymbolMap       *MySyncMap[string, *myxcoinapi.PublicRestPublicSymbolsResRow]
	linearFuturesSymbolMap         *MySyncMap[string, *myxcoinapi.PublicRestPublicSymbolsResRow]
	spotExchangeInfoMap            *MySyncMap[string, TradeSymbolInfo]
	linearPerpetualExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
	linearFuturesExchangeInfoMap   *MySyncMap[string, TradeSymbolInfo]
}

type xcoinSymbolInfo struct {
	symbolInfo
}

func (e *XcoinExchangeInfo) loadExchangeInfo() error {
	e.spotSymbolMap = GetPointer(NewMySyncMap[string, *myxcoinapi.PublicRestPublicSymbolsResRow]())
	e.linearPerpetualSymbolMap = GetPointer(NewMySyncMap[string, *myxcoinapi.PublicRestPublicSymbolsResRow]())
	e.linearFuturesSymbolMap = GetPointer(NewMySyncMap[string, *myxcoinapi.PublicRestPublicSymbolsResRow]())

	var err error
	publicRestClient := xcoin.NewRestClient("", "").PublicRestClient()
	spotRes, err := publicRestClient.NewPublicRestPublicSymbols().BusinessType(XCOIN_ACCOUNT_TYPE_SPOT.String()).Do()
	if err != nil {
		return err
	}
	linearPerpetualRes, err := publicRestClient.NewPublicRestPublicSymbols().BusinessType(XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL.String()).Do()
	if err != nil {
		return err
	}
	linearFuturesRes, err := publicRestClient.NewPublicRestPublicSymbols().BusinessType(XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES.String()).Do()
	if err != nil {
		return err
	}
	for _, v := range spotRes.Data {
		newSymbol := v
		e.spotSymbolMap.Store(v.Symbol, &newSymbol)
	}
	for _, v := range linearPerpetualRes.Data {
		newSymbol := v
		e.linearPerpetualSymbolMap.Store(v.Symbol, &newSymbol)
	}
	for _, v := range linearFuturesRes.Data {
		newSymbol := v
		e.linearFuturesSymbolMap.Store(v.Symbol, &newSymbol)
	}

	e.isLoaded = true
	e.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.linearPerpetualExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.linearFuturesExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	return nil
}

func (e *XcoinExchangeInfo) Refresh() error {
	e.isLoaded = false
	return e.loadExchangeInfo()
}

func (e *XcoinExchangeInfo) getSymbolInfoFromMap(symbolMap *MySyncMap[string, *myxcoinapi.PublicRestPublicSymbolsResRow], symbol string) (*myxcoinapi.PublicRestPublicSymbolsResRow, bool) {
	if v, ok := symbolMap.Load(symbol); ok {
		return v, true
	}
	return nil, false
}

func (e *XcoinExchangeInfo) GetSymbolInfo(accountType, symbol string) (TradeSymbolInfo, error) {
	if !e.isLoaded {
		err := e.loadExchangeInfo()
		if err != nil {
			return nil, err
		}
	}

	var pricePrecision, amtPrecision, maxOrderNum int
	var tickSize, minPrice, maxPrice = "0", "0", "0"
	var lotSize, minAmt, maxLmtAmt, maxMktAmt = "0", "0", "0", "0"
	var minNotional = "0"
	var baseCoin, quoteCoin string
	var isContract, isContractAmt bool
	var contractSize, contractCoin, contractType = "0", "", ""
	var minLeverage, maxLeverage, stepLeverage = "0", "0", "0"
	var v *myxcoinapi.PublicRestPublicSymbolsResRow
	var ok bool

	switch XcoinAccountType(accountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT:
		v, ok = e.getSymbolInfoFromMap(e.spotSymbolMap, symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseCurrency, v.QuoteCurrency
		isContract, isContractAmt = false, false
	case XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL:
		v, ok = e.getSymbolInfoFromMap(e.linearPerpetualSymbolMap, symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseCurrency, v.QuoteCurrency
		isContract, isContractAmt = true, true
		contractSize, contractCoin = v.CtVal, v.BaseCurrency
		minLeverage, maxLeverage, stepLeverage = "1", v.MaxLeverage, "0.01"
		contractType = v.ContractType
	case XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
		v, ok = e.getSymbolInfoFromMap(e.linearFuturesSymbolMap, symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseCurrency, v.QuoteCurrency
		isContract, isContractAmt = true, true
		contractSize, contractCoin = v.CtVal, v.BaseCurrency
		minLeverage, maxLeverage, stepLeverage = "1", v.MaxLeverage, "0.01"
		contractType = v.ContractType
	default:
		return nil, ErrorAccountType
	}

	if v.TickSize != "" {
		tickSize = v.TickSize
		pricePrecision = countDecimalPlaces(v.TickSize)
	}
	if v.QuantityPrecision != "" {
		amtPrecision = int(stringToInt64(v.QuantityPrecision))
		lotSize = getSizeFromPrecision(amtPrecision)
	}
	if v.PricePrecision != "" && pricePrecision == 0 {
		pricePrecision = int(stringToInt64(v.PricePrecision))
		tickSize = getSizeFromPrecision(pricePrecision)
	}

	if v.OrderParameters.MinOrderQty != "" {
		minAmt = v.OrderParameters.MinOrderQty
	}
	if v.OrderParameters.MaxLmtOrderQty != "" {
		maxLmtAmt = v.OrderParameters.MaxLmtOrderQty
	}
	if v.OrderParameters.MaxMktOrderQty != "" {
		maxMktAmt = v.OrderParameters.MaxMktOrderQty
	}
	if maxLmtAmt == "0" && v.OrderParameters.MaxLmtOrderAmt != "" {
		maxLmtAmt = v.OrderParameters.MaxLmtOrderAmt
	}
	if maxMktAmt == "0" && v.OrderParameters.MaxMktOrderAmt != "" {
		maxMktAmt = v.OrderParameters.MaxMktOrderAmt
	}
	if v.OrderParameters.MaxOrderNum != "" {
		maxOrderNum = int(stringToInt64(v.OrderParameters.MaxOrderNum))
	}
	if v.OrderParameters.MinOrderAmt != "" {
		minNotional = v.OrderParameters.MinOrderAmt
	}
	if v.PriceParameters.MinLmtPriceDown != "" {
		minPrice = v.PriceParameters.MinLmtPriceDown
	}
	if v.PriceParameters.MaxLmtPriceUp != "" {
		maxPrice = v.PriceParameters.MaxLmtPriceUp
	}

	return &xcoinSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:       XCOIN_NAME.String(),
			AccountType:    v.BusinessType,
			Symbol:         v.Symbol,
			BaseCoin:       baseCoin,
			QuoteCoin:      quoteCoin,
			IsTrading:      v.Status == "trading",
			IsContract:     isContract,
			IsContractAmt:  isContractAmt,
			ContractSize:   contractSize,
			ContractCoin:   contractCoin,
			ContractType:   contractType,
			PricePrecision: pricePrecision,
			AmtPrecision:   amtPrecision,
			TickSize:       tickSize,
			MinPrice:       minPrice,
			MaxPrice:       maxPrice,
			LotSize:        lotSize,
			MinAmt:         minAmt,
			MaxLmtAmt:      maxLmtAmt,
			MaxMktAmt:      maxMktAmt,
			MaxLeverage:    maxLeverage,
			MinLeverage:    minLeverage,
			StepLeverage:   stepLeverage,
			MaxOrderNum:    maxOrderNum,
			MinNotional:    minNotional,
		},
	}}, nil
}

func (e *XcoinExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {
	if !e.isLoaded {
		err := e.Refresh()
		if err != nil {
			return nil, err
		}
	}

	var symbolInfoList []TradeSymbolInfo
	switch XcoinAccountType(accountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT:
		e.spotSymbolMap.Range(func(key string, value *myxcoinapi.PublicRestPublicSymbolsResRow) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL:
		e.linearPerpetualSymbolMap.Range(func(key string, value *myxcoinapi.PublicRestPublicSymbolsResRow) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
		e.linearFuturesSymbolMap.Range(func(key string, value *myxcoinapi.PublicRestPublicSymbolsResRow) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	default:
		return nil, ErrorAccountType
	}

	return symbolInfoList, nil
}
