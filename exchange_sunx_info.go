package mytrade

import (
	"strings"

	"github.com/Hongssd/mysunxapi"
	"github.com/shopspring/decimal"
)

type sunxContractInfo struct {
	mysunxapi.PublicRestPublicContractInfoResRow
	MinLeverage string `json:"min_leverage"`
	MaxLeverage string `json:"max_leverage"`
	MinVolume   string `json:"min_volume"`
	MaxVolume   string `json:"max_volume"`
}

type SunxExchangeInfo struct {
	ExchangeBase

	isLoaded            bool
	swapSymbolMap       *MySyncMap[string, *sunxContractInfo]
	swapExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
}

type sunxSymbolInfo struct {
	symbolInfo
}

func (i *SunxExchangeInfo) loadExchangeInfo() error {
	i.swapSymbolMap = GetPointer(NewMySyncMap[string, *sunxContractInfo]())

	var err error
	swapRes, err := sunx.NewPublicRestClient().NewPublicRestPublicContractInfo().Do()
	if err != nil {
		return err
	}

	for _, v := range swapRes.Data {
		riskLimitRes, err := sunx.NewPublicRestClient().NewPublicRestPublicRiskLimit().ContractCode(v.ContractCode).Do()
		if err != nil {
			log.Errorf("[%s] %s", v.ContractCode, err.Error())
			continue
		}
		if len(riskLimitRes.Data) == 0 {
			log.Errorf("[%s] %s", v.ContractCode, "data is empty")
			continue
		}
		maxLeverage := riskLimitRes.Data[0].MaxLever
		minVolume := riskLimitRes.Data[0].MinVolume
		maxVolume := riskLimitRes.Data[0].MaxVolume

		i.swapSymbolMap.Store(v.ContractCode, &sunxContractInfo{
			PublicRestPublicContractInfoResRow: v,
			MinLeverage:                        "1",
			MaxLeverage:                        maxLeverage,
			MinVolume:                          minVolume,
			MaxVolume:                          maxVolume,
		})
	}

	i.isLoaded = true
	i.swapExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())

	return nil
}

func (i *SunxExchangeInfo) Refresh() error {
	i.isLoaded = false
	return i.loadExchangeInfo()
}

func (e *SunxExchangeInfo) GetSymbolInfo(accountType string, symbol string) (TradeSymbolInfo, error) {
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

	switch SunxAccountType(accountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		v, ok := e.swapSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		splitContractCode := strings.Split(v.ContractCode, "-")
		if len(splitContractCode) == 2 {
			baseCoin, quoteCoin = splitContractCode[0], splitContractCode[1]
		}
		isContract, isContractAmt = true, true
		contractSize = decimal.NewFromFloat(v.ContractSize).String()

		// Sunx 只有正向合约
		contractCoin = baseCoin

		minLeverage, maxLeverage, stepLeverage = v.MinLeverage, v.MaxLeverage, "1"
		contractType = v.ContractType

		// 价格精度
		tickSize = decimal.NewFromFloat(v.PriceTick).String()
		pricePrecision = countDecimalPlaces(tickSize)

		// 下单数量精度
		lotSize = decimal.NewFromFloat(v.ContractSize).String()
		amtPrecision = countDecimalPlaces(lotSize)

		// 最小下单数量
		minAmt = v.MinVolume

		maxMktAmt = v.MaxVolume
		maxLmtAmt = v.MaxVolume

		isTrading = v.ContractStatus == 1
	default:
		return nil, ErrorAccountType
	}

	return &sunxSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:      SUNX_NAME.String(),
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

func (e *SunxExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {
	if !e.isLoaded {
		err := e.Refresh()
		if err != nil {
			return nil, err
		}
	}

	var symbolInfoList []TradeSymbolInfo
	switch SunxAccountType(accountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		e.swapSymbolMap.Range(func(key string, value *sunxContractInfo) bool {
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
