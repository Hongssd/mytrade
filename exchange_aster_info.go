package mytrade

import (
	"github.com/Hongssd/myasterapi"
)

type AsterExchangeInfo struct {
	ExchangeBase
	isLoaded              bool
	spotSymbolMap         *MySyncMap[string, *myasterapi.SpotExchangeInfoResSymbol]
	futureSymbolMap       *MySyncMap[string, *myasterapi.FutureExchangeInfoResSymbol]
	spotExchangeInfoMap   *MySyncMap[string, TradeSymbolInfo]
	futureExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
}

type asterSymbolInfo struct {
	symbolInfo
}

func (b *AsterExchangeInfo) loadExchangeInfo() error {
	b.spotSymbolMap = GetPointer(NewMySyncMap[string, *myasterapi.SpotExchangeInfoResSymbol]())
	b.futureSymbolMap = GetPointer(NewMySyncMap[string, *myasterapi.FutureExchangeInfoResSymbol]())

	var err error
	spotRes, err := aster.NewSpotRestClient("", "").NewExchangeInfo().Do()
	if err != nil {
		return err
	}
	futureRes, err := aster.NewFutureRestClient("", "").NewExchangeInfo().Do()
	if err != nil {
		return err
	}

	for _, v := range spotRes.Symbols {
		newSymbol := v
		b.spotSymbolMap.Store(v.Symbol, &newSymbol)
	}
	for _, v := range futureRes.Symbols {
		newSymbol := v
		b.futureSymbolMap.Store(v.Symbol, &newSymbol)
	}

	b.isLoaded = true
	b.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	b.futureExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())

	return nil
}

func (b *AsterExchangeInfo) Refresh() error {
	b.isLoaded = false
	return b.loadExchangeInfo()
}

func (b *AsterExchangeInfo) GetSymbolInfo(accountType string, symbol string) (TradeSymbolInfo, error) {
	if !b.isLoaded {
		err := b.Refresh()
		if err != nil {
			return nil, err
		}
	}

	var pricePrecision, amtPrecision, maxOrderNum int
	var tickSize, minPrice, maxPrice = "0", "0", "0"
	var lotSize, minAmt, maxLmtAmt, maxMktAmt = "0", "0", "0", "0"
	var minNotional = "0"
	var baseCoin, quoteCoin string
	var isTrading, isContract, isContractAmt bool
	var contractSize, contractCoin, contractType = "0", "", ""
	var minLeverage, maxLeverage, stepLeverage = "0", "0", "0"
	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		v, ok := b.spotSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseAsset, v.QuoteAsset
		isTrading, isContract, isContractAmt = v.Status == "TRADING", false, false
		for _, m := range v.Filters {
			switch m["filterType"].(string) {
			case "PRICE_FILTER":
				pricePrecision = countDecimalPlaces(m["tickSize"].(string))
				tickSize = m["tickSize"].(string)
				minPrice = m["minPrice"].(string)
				maxPrice = m["maxPrice"].(string)
			case "LOT_SIZE":
				amtPrecision = countDecimalPlaces(m["stepSize"].(string))
				lotSize = m["stepSize"].(string)
				minAmt = m["minQty"].(string)
				maxLmtAmt = m["maxQty"].(string)
			case "MARKET_LOT_SIZE":
				maxMktAmt = m["maxQty"].(string)
			case "MAX_NUM_ORDERS":
				maxOrderNum = int(m["maxNumOrders"].(float64))
			case "NOTIONAL":
				minNotional = m["minNotional"].(string)
			default:
				continue
			}
		}
	case ASTER_AC_FUTURE:
		v, ok := b.futureSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseAsset, v.QuoteAsset
		isTrading, isContract, isContractAmt = v.Status == "TRADING", true, false
		contractType = v.ContractType
		for _, m := range v.Filters {
			switch m["filterType"].(string) {
			case "PRICE_FILTER":
				pricePrecision = countDecimalPlaces(m["tickSize"].(string))
				tickSize = m["tickSize"].(string)
				minPrice = m["minPrice"].(string)
				maxPrice = m["maxPrice"].(string)
			case "LOT_SIZE":
				amtPrecision = countDecimalPlaces(m["stepSize"].(string))
				lotSize = m["stepSize"].(string)
				minAmt = m["minQty"].(string)
				maxLmtAmt = m["maxQty"].(string)
			case "MARKET_LOT_SIZE":
				maxMktAmt = m["maxQty"].(string)
			case "MAX_NUM_ORDERS":
				maxOrderNum = int(m["limit"].(float64))
			case "MIN_NOTIONAL":
				minNotional = m["notional"].(string)
			default:
				continue
			}
		}

	default:
		return nil, ErrorAccountType
	}
	return &asterSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:      ASTER_NAME.String(),
			AccountType:   accountType,
			Symbol:        symbol,
			BaseCoin:      baseCoin,
			QuoteCoin:     quoteCoin,
			IsTrading:     isTrading,
			IsContract:    isContract,
			IsContractAmt: isContractAmt,
			ContractSize:  contractSize,
			ContractCoin:  contractCoin,
			ContractType:  contractType,

			PricePrecision: pricePrecision,
			AmtPrecision:   amtPrecision,

			TickSize: tickSize,
			MinPrice: minPrice,
			MaxPrice: maxPrice,

			LotSize:   lotSize,
			MinAmt:    minAmt,
			MaxLmtAmt: maxLmtAmt,
			MaxMktAmt: maxMktAmt,

			MaxLeverage:  maxLeverage,
			MinLeverage:  minLeverage,
			StepLeverage: stepLeverage,
			MaxOrderNum:  maxOrderNum,
			MinNotional:  minNotional,
		},
	}}, nil
}

func (b *AsterExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {

	if !b.isLoaded {
		err := b.Refresh()
		if err != nil {
			return nil, err
		}
	}

	var symbolInfoList []TradeSymbolInfo

	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT:
		b.spotSymbolMap.Range(func(key string, value *myasterapi.SpotExchangeInfoResSymbol) bool {
			symbolInfo, err := b.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case ASTER_AC_FUTURE:
		b.futureSymbolMap.Range(func(key string, value *myasterapi.FutureExchangeInfoResSymbol) bool {
			symbolInfo, err := b.GetSymbolInfo(accountType, key)
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
