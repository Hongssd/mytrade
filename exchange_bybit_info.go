package mytrade

import (
	"github.com/Hongssd/mybybitapi"
)

type BybitExchangeInfo struct {
	exchangeBase
	isLoaded               bool
	spotSymbolMap          *MySyncMap[string, *mybybitapi.InstrumentsInfoResRow]
	linearSymbolMap        *MySyncMap[string, *mybybitapi.InstrumentsInfoResRow]
	inverseSymbolMap       *MySyncMap[string, *mybybitapi.InstrumentsInfoResRow]
	spotExchangeInfoMap    *MySyncMap[string, TradeSymbolInfo]
	linearExchangeInfoMap  *MySyncMap[string, TradeSymbolInfo]
	inverseExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
}

type bybitSymbolInfo struct {
	symbolInfo
}

func (b *BybitExchangeInfo) loadExchangeInfo() error {
	b.spotSymbolMap = GetPointer(NewMySyncMap[string, *mybybitapi.InstrumentsInfoResRow]())
	b.linearSymbolMap = GetPointer(NewMySyncMap[string, *mybybitapi.InstrumentsInfoResRow]())
	b.inverseSymbolMap = GetPointer(NewMySyncMap[string, *mybybitapi.InstrumentsInfoResRow]())
	var err error
	spotRes, err := mybybitapi.NewRestClient("", "").PublicRestClient().
		NewMarketInstrumentsInfo().Category(BYBIT_AC_SPOT.String()).Limit(1000).Do()
	if err != nil {
		return err
	}
	for _, v := range spotRes.Result.List {
		newSymbol := v
		b.spotSymbolMap.Store(v.Symbol, &newSymbol)
	}
	//bybit独有检测翻页
	for spotRes.Result.NextPageCursor != "" {
		spotRes, err = mybybitapi.NewRestClient("", "").PublicRestClient().
			NewMarketInstrumentsInfo().Category(BYBIT_AC_SPOT.String()).Limit(1000).Cursor(spotRes.Result.NextPageCursor).Do()
		if err != nil {
			return err
		}
		for _, v := range spotRes.Result.List {
			newSymbol := v
			b.spotSymbolMap.Store(v.Symbol, &newSymbol)
		}
	}

	linearRes, err := mybybitapi.NewRestClient("", "").PublicRestClient().
		NewMarketInstrumentsInfo().Category(BYBIT_AC_LINEAR.String()).Limit(1000).Do()
	if err != nil {
		return err
	}
	for _, v := range linearRes.Result.List {
		newSymbol := v
		b.linearSymbolMap.Store(v.Symbol, &newSymbol)
	}

	//bybit独有检测翻页
	for linearRes.Result.NextPageCursor != "" {
		linearRes, err = mybybitapi.NewRestClient("", "").PublicRestClient().
			NewMarketInstrumentsInfo().Category(BYBIT_AC_LINEAR.String()).Limit(1000).Cursor(linearRes.Result.NextPageCursor).Do()
		if err != nil {
			return err
		}
		for _, v := range linearRes.Result.List {
			newSymbol := v
			b.linearSymbolMap.Store(v.Symbol, &newSymbol)
		}
	}

	inverseRes, err := mybybitapi.NewRestClient("", "").PublicRestClient().
		NewMarketInstrumentsInfo().Category(BYBIT_AC_INVERSE.String()).Limit(1000).Do()
	if err != nil {
		return err
	}
	for _, v := range inverseRes.Result.List {
		newSymbol := v
		b.inverseSymbolMap.Store(v.Symbol, &newSymbol)
	}

	//bybit独有检测翻页
	for inverseRes.Result.NextPageCursor != "" {
		inverseRes, err = mybybitapi.NewRestClient("", "").PublicRestClient().
			NewMarketInstrumentsInfo().Category(BYBIT_AC_INVERSE.String()).Limit(1000).Cursor(inverseRes.Result.NextPageCursor).Do()
		if err != nil {
			return err
		}
		for _, v := range inverseRes.Result.List {
			newSymbol := v
			b.inverseSymbolMap.Store(v.Symbol, &newSymbol)
		}
	}
	b.isLoaded = true
	b.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	b.linearExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	b.inverseExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())

	return nil
}

func (b *BybitExchangeInfo) Refresh() error {
	b.isLoaded = false
	return b.loadExchangeInfo()
}

func (b *BybitExchangeInfo) GetSymbolInfo(accountType string, symbol string) (TradeSymbolInfo, error) {
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
	switch BybitAccountType(accountType) {
	case BYBIT_AC_SPOT:
		v, ok := b.spotSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseCoin, v.QuoteCoin
		isTrading, isContract, isContractAmt = v.Status == "Trading", false, false
		tickSize = v.PriceFilter.TickSize
		pricePrecision = countDecimalPlaces(tickSize)
		if v.PriceFilter.MinPrice != "" {
			minPrice = v.PriceFilter.MinPrice
		}
		if v.PriceFilter.MaxPrice != "" {
			maxPrice = v.PriceFilter.MaxPrice
		}
		lotSize = v.LotSizeFilter.BasePrecision
		amtPrecision = countDecimalPlaces(lotSize)
		minAmt = v.LotSizeFilter.MinOrderQty
		maxLmtAmt = v.LotSizeFilter.MaxOrderQty
		if v.LotSizeFilter.MaxMktOrderQty != "" {
			maxMktAmt = v.LotSizeFilter.MaxMktOrderQty
		}
		if v.LotSizeFilter.MinNotionalValue != "" {
			minNotional = v.LotSizeFilter.MinNotionalValue
		}
	case BYBIT_AC_LINEAR:
		v, ok := b.linearSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.BaseCoin, v.QuoteCoin
		isTrading, isContract, isContractAmt = v.Status == "Trading", true, false
		tickSize = v.PriceFilter.TickSize
		pricePrecision = countDecimalPlaces(tickSize)
		if v.PriceFilter.MinPrice != "" {
			minPrice = v.PriceFilter.MinPrice
		}
		if v.PriceFilter.MaxPrice != "" {
			maxPrice = v.PriceFilter.MaxPrice
		}
		lotSize = v.LotSizeFilter.QtyStep
		amtPrecision = countDecimalPlaces(lotSize)
		minAmt = v.LotSizeFilter.MinOrderQty
		maxLmtAmt = v.LotSizeFilter.MaxOrderQty
		if v.LotSizeFilter.MaxMktOrderQty != "" {
			maxMktAmt = v.LotSizeFilter.MaxMktOrderQty
		}
		if v.LotSizeFilter.MinNotionalValue != "" {
			minNotional = v.LotSizeFilter.MinNotionalValue
		}
		maxLeverage = v.LinearDetails.LeverageFilter.MaxLeverage
		minLeverage = v.LinearDetails.LeverageFilter.MinLeverage
		stepLeverage = v.LinearDetails.LeverageFilter.LeverageStep
		contractType = v.ContractType
	case BYBIT_AC_INVERSE:
		v, ok := b.inverseSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.QuoteCoin, v.BaseCoin
		isTrading, isContract, isContractAmt = v.Status == "Trading", true, true
		contractSize, contractCoin = "1", v.QuoteCoin
		tickSize = v.PriceFilter.TickSize
		pricePrecision = countDecimalPlaces(tickSize)
		if v.PriceFilter.MinPrice != "" {
			minPrice = v.PriceFilter.MinPrice
		}
		if v.PriceFilter.MaxPrice != "" {
			maxPrice = v.PriceFilter.MaxPrice
		}
		lotSize = v.LotSizeFilter.QtyStep
		amtPrecision = countDecimalPlaces(lotSize)
		minAmt = v.LotSizeFilter.MinOrderQty
		maxLmtAmt = v.LotSizeFilter.MaxOrderQty
		if v.LotSizeFilter.MaxMktOrderQty != "" {
			maxMktAmt = v.LotSizeFilter.MaxMktOrderQty
		}
		if v.LotSizeFilter.MinNotionalValue != "" {
			minNotional = v.LotSizeFilter.MinNotionalValue
		}
		maxLeverage = v.LinearDetails.LeverageFilter.MaxLeverage
		minLeverage = v.LinearDetails.LeverageFilter.MinLeverage
		stepLeverage = v.LinearDetails.LeverageFilter.LeverageStep
		contractType = v.ContractType
	default:
		return nil, ErrorAccountType
	}
	return &bybitSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:      BYBIT_NAME.String(),
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
