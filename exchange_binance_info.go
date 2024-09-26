package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"strconv"
)

type BinanceExchangeInfo struct {
	ExchangeBase
	isLoaded              bool
	spotSymbolMap         *MySyncMap[string, *mybinanceapi.SpotExchangeInfoResSymbol]
	futureSymbolMap       *MySyncMap[string, *mybinanceapi.FutureExchangeInfoResSymbol]
	swapSymbolMap         *MySyncMap[string, *mybinanceapi.SwapExchangeInfoResSymbol]
	spotExchangeInfoMap   *MySyncMap[string, TradeSymbolInfo]
	futureExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
	swapExchangeInfoMap   *MySyncMap[string, TradeSymbolInfo]
}

type binanceSymbolInfo struct {
	symbolInfo
}

func (b *BinanceExchangeInfo) loadExchangeInfo() error {
	b.spotSymbolMap = GetPointer(NewMySyncMap[string, *mybinanceapi.SpotExchangeInfoResSymbol]())
	b.futureSymbolMap = GetPointer(NewMySyncMap[string, *mybinanceapi.FutureExchangeInfoResSymbol]())
	b.swapSymbolMap = GetPointer(NewMySyncMap[string, *mybinanceapi.SwapExchangeInfoResSymbol]())
	var err error
	spotRes, err := binance.NewSpotRestClient("", "").NewExchangeInfo().Do()
	if err != nil {
		return err
	}
	futureRes, err := binance.NewFutureRestClient("", "").NewExchangeInfo().Do()
	if err != nil {
		return err
	}
	swapRes, err := binance.NewSwapRestClient("", "").NewExchangeInfo().Do()
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
	for _, v := range swapRes.Symbols {
		newSymbol := v
		b.swapSymbolMap.Store(v.Symbol, &newSymbol)
	}
	b.isLoaded = true
	b.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	b.futureExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	b.swapExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())

	return nil
}

func (b *BinanceExchangeInfo) Refresh() error {
	b.isLoaded = false
	return b.loadExchangeInfo()
}

func (b *BinanceExchangeInfo) GetSymbolInfo(accountType string, symbol string) (TradeSymbolInfo, error) {
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
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
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
	case BN_AC_FUTURE:
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
	case BN_AC_SWAP:
		v, ok := b.swapSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.QuoteAsset, v.BaseAsset
		isTrading, isContract, isContractAmt = v.ContractStatus == "TRADING", true, true
		contractSize = strconv.FormatInt(v.ContractSize, 10)
		contractCoin = v.QuoteAsset
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
	return &binanceSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:      BINANCE_NAME.String(),
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

func (b *BinanceExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {

	var symbolInfoList []TradeSymbolInfo

	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		b.spotSymbolMap.Range(func(key string, value *mybinanceapi.SpotExchangeInfoResSymbol) bool {
			symbolInfo, err := b.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case BN_AC_FUTURE:
		b.futureSymbolMap.Range(func(key string, value *mybinanceapi.FutureExchangeInfoResSymbol) bool {
			symbolInfo, err := b.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case BN_AC_SWAP:
		b.swapSymbolMap.Range(func(key string, value *mybinanceapi.SwapExchangeInfoResSymbol) bool {
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
