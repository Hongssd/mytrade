package mytrade

import (
	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

type GateExchangeInfo struct {
	ExchangeBase
	isLoaded          bool
	spotSymbolMap     *MySyncMap[string, *mygateapi.PublicRestSpotCurrencyPairCommon]
	futuresSymbolMap  *MySyncMap[string, *mygateapi.ContractCommon]
	deliverySymbolMap *MySyncMap[string, *mygateapi.DeliveryContractCommon]

	spotExchangeInfoMap     *MySyncMap[string, TradeSymbolInfo]
	futuresExchangeInfoMap  *MySyncMap[string, TradeSymbolInfo]
	deliveryExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
}

type gateSymbolInfo struct {
	symbolInfo
}

func (e *GateExchangeInfo) loadExchangeInfo() error {
	e.spotSymbolMap = GetPointer(NewMySyncMap[string, *mygateapi.PublicRestSpotCurrencyPairCommon]())
	e.futuresSymbolMap = GetPointer(NewMySyncMap[string, *mygateapi.ContractCommon]())
	e.deliverySymbolMap = GetPointer(NewMySyncMap[string, *mygateapi.DeliveryContractCommon]())

	var err error
	spotRes, err := mygateapi.NewRestClient("", "").PublicRestClient().
		NewPublicRestCurrencyPairsAll().Do()
	if err != nil {
		return err
	}

	futuresRes, err := mygateapi.NewRestClient("", "").PublicRestClient().
		NewPublicRestFuturesSettleContracts().Settle("usdt").Do()
	if err != nil {
		return err
	}

	deliveryRes, err := mygateapi.NewRestClient("", "").PublicRestClient().
		NewPublicRestDeliverySettleContracts().Settle("usdt").Do()
	if err != nil {
		return err
	}

	for _, v := range spotRes.Data {
		newSymbol := v
		e.spotSymbolMap.Store(v.ID, &newSymbol)
	}
	for _, v := range deliveryRes.Data {
		newSymbol := v
		e.deliverySymbolMap.Store(v.Name, &newSymbol)
	}
	for _, v := range futuresRes.Data {
		newSymbol := v
		e.futuresSymbolMap.Store(v.Name, &newSymbol)
	}
	e.isLoaded = true
	e.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.deliveryExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.futuresExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	return nil
}

func (e *GateExchangeInfo) Refresh() error {
	e.isLoaded = false
	return e.loadExchangeInfo()
}

func (e *GateExchangeInfo) GetSymbolInfo(accountType string, symbol string) (TradeSymbolInfo, error) {
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

	var isTrading bool

	switch GateAccountType(accountType) {
	case GATE_AC_SPOT:
		v, ok := e.spotSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.Base, v.Quote
		isContract, isContractAmt = false, false
		contractSize, contractCoin = "0", ""

		pricePrecision = v.Precision
		tickSize = getSizeFromPrecision(v.Precision)
		amtPrecision = v.AmountPrecision
		lotSize = getSizeFromPrecision(v.AmountPrecision)
		if v.MinBaseAmount != "null" {
			minAmt = v.MinBaseAmount
		}
		if v.MinQuoteAmount != "null" {
			minNotional = v.MinQuoteAmount
		}

		if v.MaxBaseAmount != "null" {
			maxMktAmt = v.MaxBaseAmount
			maxLmtAmt = v.MaxBaseAmount
		}
		if v.TradeStatus == "tradable" {
			isTrading = true
		}

	case GATE_AC_FUTURES:
		v, ok := e.futuresSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		sp := strings.Split(v.Name, "_")
		if len(sp) != 2 {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = sp[0], sp[1]
		isContract, isContractAmt = true, true
		contractSize = v.QuantoMultiplier
		if v.Type == "direct" {
			//正向合约
			contractCoin = baseCoin
		} else {
			//反向合约
			contractCoin = quoteCoin
		}
		minLeverage, maxLeverage, stepLeverage = v.LeverageMin, v.LeverageMax, "1"
		contractType = v.Type

		markPrice, _ := decimal.NewFromString(v.MarkPrice)
		priceDeviate, _ := decimal.NewFromString(v.OrderPriceDeviate)

		minPrice = markPrice.Sub(markPrice.Mul(priceDeviate)).String()
		maxPrice = markPrice.Add(markPrice.Mul(priceDeviate)).String()

		pricePrecision = countDecimalPlaces(v.OrderPriceRound)
		tickSize = v.OrderPriceRound
		amtPrecision = countDecimalPlaces(strconv.FormatInt(v.OrderSizeMin, 10))
		lotSize = strconv.FormatInt(v.OrderSizeMin, 10)

		minAmt = strconv.FormatInt(v.OrderSizeMin, 10)

		maxMktAmt = strconv.FormatInt(v.OrderSizeMax, 10)
		maxLmtAmt = strconv.FormatInt(v.OrderSizeMax, 10)

		isTrading = true
	case GATE_AC_DELIVERY:
		v, ok := e.deliverySymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		sp := strings.Split(v.Underlying, "_")
		if len(sp) != 2 {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = sp[0], sp[1]
		isContract, isContractAmt = true, true
		contractSize = v.QuantoMultiplier
		if v.Type == "direct" {
			//正向合约
			contractCoin = baseCoin
		} else {
			//反向合约
			contractCoin = quoteCoin
		}
		minLeverage, maxLeverage, stepLeverage = v.LeverageMin, v.LeverageMax, "1"
		contractType = v.Type

		markPrice, _ := decimal.NewFromString(v.MarkPrice)
		priceDeviate, _ := decimal.NewFromString(v.OrderPriceDeviate)

		minPrice = markPrice.Sub(markPrice.Mul(priceDeviate)).String()
		maxPrice = markPrice.Add(markPrice.Mul(priceDeviate)).String()

		pricePrecision = countDecimalPlaces(v.OrderPriceRound)
		tickSize = v.OrderPriceRound
		amtPrecision = countDecimalPlaces(strconv.FormatInt(v.OrderSizeMin, 10))
		lotSize = strconv.FormatInt(v.OrderSizeMin, 10)

		minAmt = strconv.FormatInt(v.OrderSizeMin, 10)

		maxMktAmt = strconv.FormatInt(v.OrderSizeMax, 10)
		maxLmtAmt = strconv.FormatInt(v.OrderSizeMax, 10)

		isTrading = true
	default:
		return nil, ErrorAccountType
	}

	return &gateSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:      GATE_NAME.String(),
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

func (e *GateExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {

	var symbolInfoList []TradeSymbolInfo
	switch GateAccountType(accountType) {
	case GATE_AC_SPOT:
		e.spotSymbolMap.Range(func(key string, value *mygateapi.PublicRestSpotCurrencyPairCommon) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case GATE_AC_FUTURES:
		e.futuresSymbolMap.Range(func(key string, value *mygateapi.ContractCommon) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case GATE_AC_DELIVERY:
		e.deliverySymbolMap.Range(func(key string, value *mygateapi.DeliveryContractCommon) bool {
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
