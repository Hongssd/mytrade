package mytrade

import (
	"github.com/Hongssd/myokxapi"
)

type OkxExchangeInfo struct {
	exchangeBase
	isLoaded               bool
	spotSymbolMap          *MySyncMap[string, *myokxapi.PublicRestPublicInstrumentsResRow]
	swapSymbolMap          *MySyncMap[string, *myokxapi.PublicRestPublicInstrumentsResRow]
	futuresSymbolMap       *MySyncMap[string, *myokxapi.PublicRestPublicInstrumentsResRow]
	spotExchangeInfoMap    *MySyncMap[string, TradeSymbolInfo]
	swapExchangeInfoMap    *MySyncMap[string, TradeSymbolInfo]
	futuresExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
}

type okxSymbolInfo struct {
	symbolInfo
}

func (e *OkxExchangeInfo) loadExchangeInfo() error {
	e.spotSymbolMap = GetPointer(NewMySyncMap[string, *myokxapi.PublicRestPublicInstrumentsResRow]())
	e.swapSymbolMap = GetPointer(NewMySyncMap[string, *myokxapi.PublicRestPublicInstrumentsResRow]())
	e.futuresSymbolMap = GetPointer(NewMySyncMap[string, *myokxapi.PublicRestPublicInstrumentsResRow]())

	var err error
	spotRes, err := okx.NewRestClient("", "", "").PublicRestClient().
		NewPublicRestPublicInstruments().InstType(OKX_AC_SPOT.String()).Do()
	if err != nil {
		return err
	}
	swapRes, err := okx.NewRestClient("", "", "").PublicRestClient().
		NewPublicRestPublicInstruments().InstType(OKX_AC_SWAP.String()).Do()
	if err != nil {
		return err
	}
	futuresRes, err := okx.NewRestClient("", "", "").PublicRestClient().
		NewPublicRestPublicInstruments().InstType(OKX_AC_FUTURES.String()).Do()
	if err != nil {
		return err
	}
	for _, v := range spotRes.Data {
		newSymbol := v
		e.spotSymbolMap.Store(v.InstId, &newSymbol)
	}
	for _, v := range swapRes.Data {
		newSymbol := v
		e.swapSymbolMap.Store(v.InstId, &newSymbol)
	}
	for _, v := range futuresRes.Data {
		newSymbol := v
		e.futuresSymbolMap.Store(v.InstId, &newSymbol)
	}
	e.isLoaded = true
	e.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.swapExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.futuresExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	return nil
}

func (e *OkxExchangeInfo) Refresh() error {
	e.isLoaded = false
	return e.loadExchangeInfo()
}
func (e *OkxExchangeInfo) getSymbolInfoFromMap(symbolMap *MySyncMap[string, *myokxapi.PublicRestPublicInstrumentsResRow], symbol string) (*myokxapi.PublicRestPublicInstrumentsResRow, bool) {
	if v, ok := symbolMap.Load(symbol); ok {
		return v, true
	}
	return nil, false
}

func (e *OkxExchangeInfo) GetSymbolInfo(accountType string, symbol string) (TradeSymbolInfo, error) {
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
	var v *myokxapi.PublicRestPublicInstrumentsResRow
	var ok bool
	switch OkxAccountType(accountType) {
	case OKX_AC_SPOT:
		v, ok = e.getSymbolInfoFromMap(e.spotSymbolMap, symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseCcy, v.QuoteCcy
		isContract, isContractAmt = false, false
		contractSize, contractCoin = "0", ""
	case OKX_AC_SWAP:
		v, ok = e.getSymbolInfoFromMap(e.swapSymbolMap, symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.CtValCcy, v.SettleCcy
		isContract, isContractAmt = true, true
		contractSize, contractCoin = v.CtVal, v.CtValCcy
		minLeverage, maxLeverage, stepLeverage = "1", v.Lever, "0.01"
		contractType = v.CtType
	case OKX_AC_FUTURES:
		v, ok = e.getSymbolInfoFromMap(e.futuresSymbolMap, symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.CtValCcy, v.SettleCcy
		isContract, isContractAmt = true, true
		contractSize, contractCoin = v.CtVal, v.CtValCcy
		minLeverage, maxLeverage, stepLeverage = "1", v.Lever, "0.01"
		contractType = v.CtType
	default:
		return nil, ErrorAccountType
	}

	pricePrecision = countDecimalPlaces(v.TickSz)
	tickSize = v.TickSz
	amtPrecision = countDecimalPlaces(v.LotSz)
	lotSize = v.LotSz
	minAmt = v.MinSz
	maxLmtAmt = v.MaxLmtSz
	if v.MaxMktSz != "" {
		maxMktAmt = v.MaxMktSz
	}

	return &okxSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:      OKX_NAME.String(),
			AccountType:   v.InstType,
			Symbol:        v.InstId,
			BaseCoin:      baseCoin,
			QuoteCoin:     quoteCoin,
			IsTrading:     v.State == "live",
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

func (e *OkxExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {

	var symbolInfoList []TradeSymbolInfo
	switch OkxAccountType(accountType) {
	case OKX_AC_SPOT:
		e.spotSymbolMap.Range(func(key string, value *myokxapi.PublicRestPublicInstrumentsResRow) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case OKX_AC_SWAP:
		e.swapSymbolMap.Range(func(key string, value *myokxapi.PublicRestPublicInstrumentsResRow) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case OKX_AC_FUTURES:
		e.futuresSymbolMap.Range(func(key string, value *myokxapi.PublicRestPublicInstrumentsResRow) bool {
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
