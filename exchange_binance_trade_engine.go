package mytrade

import (
	"errors"
	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
	"strconv"
	"sync"
)

type BinanceTradeEngine struct {
	exchangeBase

	bnConverter BinanceEnumConverter
	apiKey      string
	secretKey   string

	wsSpotAccount   *mybinanceapi.SpotWsStreamClient
	wsFutureAccount *mybinanceapi.FutureWsStreamClient
	wsSwapAccount   *mybinanceapi.SwapWsStreamClient

	wsSpotWsApi   *mybinanceapi.SpotWsStreamClient
	wsFutureWsApi *mybinanceapi.FutureWsStreamClient
	wsSwapWsApi   *mybinanceapi.SwapWsStreamClient
}

func (b BinanceTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}

func (b BinanceTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}

func (b BinanceTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (b BinanceTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	var orders []*Order
	binance := mybinanceapi.MyBinance{}
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		res, err := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewOpenOrders().Symbol(req.Symbol).Do()
		if err != nil {
			return nil, err
		}

		for _, order := range *res {
			orders = append(orders, &Order{
				Exchange:      BINANCE_NAME.String(),
				AccountType:   req.AccountType,
				Symbol:        req.Symbol,
				OrderId:       strconv.FormatInt(order.OrderId, 10),
				ClientOrderId: order.ClientOrderId,
				Price:         order.Price,
				Quantity:      order.OrigQty,
				ExecutedQty:   order.ExecutedQty,
				CumQuoteQty:   order.CummulativeQuoteQty,
				Status:        b.bnConverter.FromBNOrderStatus(order.Status),
				Type:          b.bnConverter.FromBNOrderType(order.Type),
				Side:          b.bnConverter.FromBNOrderSide(order.Side),
				TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
				CreateTime:    order.Time,
				UpdateTime:    order.UpdateTime,
			})
		}

	case BN_AC_FUTURE:
		res, err := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewOpenOrders().Symbol(req.Symbol).Do()
		if err != nil {
			return nil, err
		}

		for _, order := range *res {
			orders = append(orders, &Order{
				Exchange:      BINANCE_NAME.String(),
				AccountType:   req.AccountType,
				Symbol:        req.Symbol,
				OrderId:       strconv.FormatInt(order.OrderId, 10),
				ClientOrderId: order.ClientOrderId,
				Price:         order.Price,
				Quantity:      order.OrigQty,
				ExecutedQty:   order.ExecutedQty,
				CumQuoteQty:   order.CumQuote,
				Status:        b.bnConverter.FromBNOrderStatus(order.Status),
				Type:          b.bnConverter.FromBNOrderType(order.Type),
				Side:          b.bnConverter.FromBNOrderSide(order.Side),
				PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
				TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
				ReduceOnly:    order.ReduceOnly,
				CreateTime:    order.Time,
				UpdateTime:    order.UpdateTime,
			})
		}
	case BN_AC_SWAP:
		res, err := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewOpenOrders().Symbol(req.Symbol).Do()
		if err != nil {
			return nil, err
		}

		for _, order := range *res {
			orders = append(orders, &Order{
				Exchange:      BINANCE_NAME.String(),
				AccountType:   req.AccountType,
				Symbol:        req.Symbol,
				OrderId:       strconv.FormatInt(order.OrderId, 10),
				ClientOrderId: order.ClientOrderId,
				Price:         order.Price,
				Quantity:      order.OrigQty,
				ExecutedQty:   order.ExecutedQty,
				CumQuoteQty:   order.CumQuote,
				Status:        b.bnConverter.FromBNOrderStatus(order.Status),
				Type:          b.bnConverter.FromBNOrderType(order.Type),
				Side:          b.bnConverter.FromBNOrderSide(order.Side),
				PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
				TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
				ReduceOnly:    order.ReduceOnly,
				CreateTime:    order.Time,
				UpdateTime:    order.UpdateTime,
			})
		}
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b BinanceTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	var order *Order
	var err error
	binance := mybinanceapi.MyBinance{}
	var isOrderIdParam bool
	var orderId int64
	var clientOrderId string
	if req.OrderId != "" {
		isOrderIdParam = true
		orderId, err = strconv.ParseInt(req.OrderId, 10, 64)
		if err != nil {
			return nil, ErrorInvalid("order id")
		}
	} else {
		isOrderIdParam = false
		clientOrderId = req.ClientOrderId
	}

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderGet().Symbol(req.Symbol)
		if isOrderIdParam {
			api = api.OrderId(orderId)
		} else {
			api = api.OrigClientOrderId(clientOrderId)
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		order = &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			OrderId:       strconv.FormatInt(res.OrderId, 10),
			ClientOrderId: res.ClientOrderId,
			Price:         res.Price,
			Quantity:      res.OrigQty,
			ExecutedQty:   res.ExecutedQty,
			CumQuoteQty:   res.CummulativeQuoteQty,
			Status:        b.bnConverter.FromBNOrderStatus(res.Status),
			Type:          b.bnConverter.FromBNOrderType(res.Type),
			Side:          b.bnConverter.FromBNOrderSide(res.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
			CreateTime:    res.Time,
			UpdateTime:    res.UpdateTime,
		}
	case BN_AC_FUTURE:
		api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderGet().Symbol(req.Symbol)
		if isOrderIdParam {
			api = api.OrderId(orderId)
		} else {
			api = api.OrigClientOrderId(clientOrderId)
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		order = &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			OrderId:       strconv.FormatInt(res.OrderId, 10),
			ClientOrderId: res.ClientOrderId,
			Price:         res.Price,
			Quantity:      res.OrigQty,
			ExecutedQty:   res.ExecutedQty,
			CumQuoteQty:   res.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(res.Status),
			Type:          b.bnConverter.FromBNOrderType(res.Type),
			Side:          b.bnConverter.FromBNOrderSide(res.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
			ReduceOnly:    res.ReduceOnly,
			CreateTime:    res.Time,
			UpdateTime:    res.UpdateTime,
		}
	case BN_AC_SWAP:
		api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderGet().Symbol(req.Symbol)
		if isOrderIdParam {
			api = api.OrderId(orderId)
		} else {
			api = api.OrigClientOrderId(clientOrderId)
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		order = &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			OrderId:       strconv.FormatInt(res.OrderId, 10),
			ClientOrderId: res.ClientOrderId,
			Price:         res.Price,
			Quantity:      res.OrigQty,
			ExecutedQty:   res.ExecutedQty,
			CumQuoteQty:   res.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(res.Status),
			Type:          b.bnConverter.FromBNOrderType(res.Type),
			Side:          b.bnConverter.FromBNOrderSide(res.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
			ReduceOnly:    res.ReduceOnly,
			CreateTime:    res.Time,
			UpdateTime:    res.UpdateTime,
		}
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b BinanceTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	var trades []*Trade
	binance := mybinanceapi.MyBinance{}
	var orderId int64
	var err error
	if req.OrderId != "" {
		orderId, err = strconv.ParseInt(req.OrderId, 10, 64)
		if err != nil {
			return nil, ErrorInvalid("order id")
		}
	}
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMyTrades().
			Symbol(req.Symbol)
		if req.OrderId != "" {
			api = api.OrderId(orderId)
		}
		if req.StartTime != 0 {
			api = api.StartTime(req.StartTime)
		}
		if req.EndTime != 0 {
			api = api.EndTime(req.EndTime)
		}
		if req.Limit != 0 {
			api = api.Limit(req.Limit)
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		for _, trade := range *res {
			var orderSide OrderSide
			if trade.IsBuyer {
				orderSide = ORDER_SIDE_BUY
			} else {
				orderSide = ORDER_SIDE_SELL
			}
			trades = append(trades, &Trade{
				Exchange:    BINANCE_NAME.String(),
				AccountType: req.AccountType,
				Symbol:      req.Symbol,
				TradeId:     strconv.FormatInt(trade.Id, 10),
				OrderId:     strconv.FormatInt(trade.OrderId, 10),
				Price:       trade.Price,
				Quantity:    trade.Qty,
				QuoteQty:    trade.QuoteQty,
				Side:        orderSide,
				FeeAmount:   trade.Commission,
				FeeCcy:      trade.CommissionAsset,
				IsMaker:     trade.IsMaker,
				Timestamp:   trade.Time,
			})
		}
	case BN_AC_FUTURE:
		api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureUserTrades().
			Symbol(req.Symbol)
		if req.OrderId != "" {
			api = api.OrderId(orderId)
		}
		if req.StartTime != 0 {
			api = api.StartTime(req.StartTime)
		}
		if req.EndTime != 0 {
			api = api.EndTime(req.EndTime)
		}
		if req.Limit != 0 {
			api = api.Limit(int64(req.Limit))
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		for _, trade := range *res {
			trades = append(trades, &Trade{
				Exchange:     BINANCE_NAME.String(),
				AccountType:  req.AccountType,
				Symbol:       req.Symbol,
				TradeId:      strconv.FormatInt(trade.Id, 10),
				OrderId:      strconv.FormatInt(trade.OrderId, 10),
				Price:        trade.Price,
				Quantity:     trade.Qty,
				QuoteQty:     trade.QuoteQty,
				Side:         b.bnConverter.FromBNOrderSide(trade.Side),
				PositionSide: b.bnConverter.FromBNPositionSide(trade.PositionSide),
				FeeAmount:    trade.Commission,
				FeeCcy:       trade.CommissionAsset,
				RealizedPnl:  trade.RealizedPnl,
				IsMaker:      trade.Maker,
				Timestamp:    trade.Time,
			})
		}
	case BN_AC_SWAP:
		api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapUserTrades().
			Symbol(req.Symbol)
		if req.OrderId != "" {
			api = api.OrderId(orderId)
		}
		if req.StartTime != 0 {
			api = api.StartTime(req.StartTime)
		}
		if req.EndTime != 0 {
			api = api.EndTime(req.EndTime)
		}
		if req.Limit != 0 {
			api = api.Limit(int64(req.Limit))
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		for _, trade := range *res {
			trades = append(trades, &Trade{
				Exchange:     BINANCE_NAME.String(),
				AccountType:  req.AccountType,
				Symbol:       req.Symbol,
				TradeId:      strconv.FormatInt(trade.Id, 10),
				OrderId:      strconv.FormatInt(trade.OrderId, 10),
				Price:        trade.Price,
				Quantity:     trade.Qty,
				QuoteQty:     trade.BaseQty,
				Side:         b.bnConverter.FromBNOrderSide(trade.Side),
				PositionSide: b.bnConverter.FromBNPositionSide(trade.PositionSide),
				FeeAmount:    trade.Commission,
				FeeCcy:       trade.CommissionAsset,
				RealizedPnl:  trade.RealizedPnl,
				IsMaker:      trade.Maker,
				Timestamp:    trade.Time,
			})
		}
	default:
		return nil, ErrorAccountType
	}

	return trades, nil
}

func (b BinanceTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	var order *Order
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		api := b.apiSpotOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSpotOrderCreate(req, res)
	case BN_AC_FUTURE:
		api := b.apiFutureOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderCreate(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSwapOrderCreate(req, res)
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}
func (b BinanceTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	var order *Order

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		api := b.apiSpotOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		if res.CancelResult != "SUCCESS" {
			return nil, errors.New("cancel order failed")
		}
		if res.NewOrderResult != "SUCCESS" {
			return nil, errors.New("cancel order success and amend order failed")
		}
		order = b.handleOrderFromSpotOrderAmend(req, res)
	case BN_AC_FUTURE:
		api := b.apiFutureOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderAmend(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSwapOrderAmend(req, res)
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b BinanceTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	var order *Order
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		api := b.apiSpotOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSpotOrderCancel(req, res)
	case BN_AC_FUTURE:
		api := b.apiFutureOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderCancel(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSwapOrderCancel(req, res)
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}

func (b BinanceTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order
	//检测长度，最多批量下5个订单
	if len(reqs) > 5 {
		return nil, ErrorInvalid("order param length require less than 5")

	}

	//检测类型是否相同
	for _, req := range reqs {
		if req.AccountType != reqs[0].AccountType {
			return nil, ErrorInvalid("order param account type require same")
		}
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT:
		//现货无批量接口，直接并发下单
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, req := range reqs {
			req := req
			wg.Add(1)
			go func() {
				defer wg.Done()
				order, err := b.CreateOrder(req)
				if err != nil {
					mu.Lock()
					orders = append(orders, b.handleOrderFromSpotBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BN_AC_FUTURE:
		api := b.apiFutureBatchOrderCreate(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderCreate(reqs, res)
	case BN_AC_SWAP:
		api := b.apiSwapBatchOrderCreate(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSwapBatchOrderCreate(reqs, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b BinanceTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order
	//检测长度，最多批量改5个订单
	if len(reqs) > 5 {
		return nil, ErrorInvalid("order param length require less than 5")

	}

	//检测类型是否相同
	for _, req := range reqs {
		if req.AccountType != reqs[0].AccountType {
			return nil, ErrorInvalid("order param account type require same")
		}
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT:
		//现货无批量接口，直接并发改单
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, req := range reqs {
			req := req
			wg.Add(1)
			go func() {
				defer wg.Done()
				order, err := b.AmendOrder(req)
				if err != nil {
					mu.Lock()
					orders = append(orders, b.handleOrderFromSpotBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BN_AC_FUTURE:
		api := b.apiFutureBatchOrderAmend(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderAmend(reqs, res)
	case BN_AC_SWAP:
		api := b.apiSwapBatchOrderAmend(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSwapBatchOrderAmend(reqs, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b BinanceTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order
	//检测长度，最多批量撤单10个订单
	if len(reqs) > 10 {
		return nil, ErrorInvalid("order param length require less than 10")

	}

	//检测类型是否相同
	for _, req := range reqs {
		if req.AccountType != reqs[0].AccountType {
			return nil, ErrorInvalid("order param account type require same")
		}
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT:
		//现货无批量接口，直接并发撤单
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, req := range reqs {
			req := req
			wg.Add(1)
			go func() {
				defer wg.Done()
				order, err := b.CancelOrder(req)
				if err != nil {
					mu.Lock()
					orders = append(orders, b.handleOrderFromSpotBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BN_AC_FUTURE:
		api, err := b.apiFutureBatchOrderCancel(reqs)
		if err != nil {
			return nil, err
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderCancel(reqs, res)
	case BN_AC_SWAP:
		api, err := b.apiSwapBatchOrderCancel(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSwapBatchOrderCancel(reqs, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}

func (b BinanceTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (b BinanceTradeEngine) SubscribeOrder(r *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	req := *r
	binance := &mybinanceapi.MyBinance{}
	var err error
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if b.wsSpotAccount == nil {
			b.wsSpotAccount, err = binance.NewSpotWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey, mybinanceapi.SPOT_WS_TYPE)
			if err != nil {
				return nil, err
			}
			err := b.wsSpotAccount.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		newPayload, err := b.wsSpotAccount.CreatePayload()
		if err != nil {
			return nil, err
		}
		//构建一个推送订单数据的中转订阅
		newSub := &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		}

		//处理不需要的订阅数据
		go func() {
			for {
				select {
				case <-newPayload.BalanceUpdatePayload.ErrChan():
					continue
				case <-newSub.closeChan:
					return
				case r := <-newPayload.BalanceUpdatePayload.ResultChan():
					_ = r
				}
			}
		}()
		go func() {
			for {
				select {
				case <-newPayload.OutboundAccountPositionPayload.ErrChan():
					continue
				case <-newSub.closeChan:
					return
				case r := <-newPayload.OutboundAccountPositionPayload.ResultChan():
					_ = r
				}
			}
		}()

		//处理订单推送订阅
		go func() {
			for {
				select {
				case err := <-newPayload.ExecutionReportPayload.ErrChan():
					newSub.errChan <- err
				case <-newSub.closeChan:
					newSub.CloseChan() <- struct{}{}
					return
				case r := <-newPayload.ExecutionReportPayload.ResultChan():
					order := Order{
						Exchange:      BINANCE_NAME.String(),
						AccountType:   req.AccountType,
						Symbol:        r.Symbol,
						OrderId:       strconv.FormatInt(r.OrderId, 10),
						ClientOrderId: r.ClientOrderId,
						Price:         r.Price,
						Quantity:      r.OrigQty,
						ExecutedQty:   r.ExecutedQty,
						CumQuoteQty:   r.CummulativeQuoteQty,
						Status:        b.bnConverter.FromBNOrderStatus(r.Status),
						Type:          b.bnConverter.FromBNOrderType(r.Type),
						Side:          b.bnConverter.FromBNOrderSide(r.Side),
						TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
						FeeAmount:     r.FeeQty,
						FeeCcy:        r.FeeAsset,
						CreateTime:    r.OrderCreateTime,
						UpdateTime:    r.Timestamp,
					}
					newSub.resultChan <- order
				}
			}
		}()

		return newSub, nil
	case BN_AC_FUTURE:
		if b.wsFutureAccount == nil {
			b.wsFutureAccount, err = binance.NewFutureWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsFutureAccount.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		newPayload, err := b.wsFutureAccount.CreatePayload()
		if err != nil {
			return nil, err
		}

		//构建一个推送订单数据的中转订阅
		newSub := &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		}

		//处理不需要的订阅数据
		go func() {
			for {
				select {
				case <-newPayload.AccountUpdatePayload.ErrChan():
					continue
				case <-newSub.closeChan:
					return
				case r := <-newPayload.AccountUpdatePayload.ResultChan():
					_ = r
				}
			}
		}()

		//处理订单推送订阅
		go func() {
			for {
				select {
				case err := <-newPayload.OrderTradeUpdatePayload.ErrChan():
					newSub.errChan <- err
				case <-newSub.closeChan:
					newSub.CloseChan() <- struct{}{}
					return
				case result := <-newPayload.OrderTradeUpdatePayload.ResultChan():
					r := result.Order
					CumQuoteQty := decimal.Zero
					avgPrice, err := decimal.NewFromString(r.AvgPrice)
					if err != nil {
						newSub.ErrChan() <- err
					}
					CumQuoteQty = avgPrice.Mul(decimal.RequireFromString(r.ExecutedQty))
					order := Order{
						Exchange:      BINANCE_NAME.String(),
						AccountType:   req.AccountType,
						Symbol:        r.Symbol,
						OrderId:       strconv.FormatInt(r.OrderId, 10),
						ClientOrderId: r.ClientOrderId,
						Price:         r.Price,
						Quantity:      r.OrigQty,
						ExecutedQty:   r.ExecutedQty,
						CumQuoteQty:   CumQuoteQty.String(),
						Status:        b.bnConverter.FromBNOrderStatus(r.Status),
						Type:          b.bnConverter.FromBNOrderType(r.Type),
						Side:          b.bnConverter.FromBNOrderSide(r.Side),
						PositionSide:  b.bnConverter.FromBNPositionSide(r.PositionSide),
						TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
						FeeAmount:     r.FeeQty,
						FeeCcy:        r.FeeAsset,
						ReduceOnly:    r.IsReduceOnly,
						CreateTime:    result.Timestamp,
						UpdateTime:    result.Timestamp,
					}
					newSub.resultChan <- order
				}
			}
		}()

		return newSub, nil
	case BN_AC_SWAP:
		if b.wsSwapAccount == nil {
			b.wsSwapAccount, err = binance.NewSwapWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsSwapAccount.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		newPayload, err := b.wsSwapAccount.CreatePayload()
		if err != nil {
			return nil, err
		}

		//构建一个推送订单数据的中转订阅
		newSub := &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		}

		//处理不需要的订阅数据
		go func() {
			for {
				select {
				case <-newPayload.AccountUpdatePayload.ErrChan():
					continue
				case <-newSub.closeChan:
					return
				case r := <-newPayload.AccountUpdatePayload.ResultChan():
					_ = r
				}
			}
		}()

		//处理订单推送订阅
		go func() {
			for {
				select {
				case err := <-newPayload.OrderTradeUpdatePayload.ErrChan():
					newSub.errChan <- err
				case <-newSub.closeChan:
					newSub.CloseChan() <- struct{}{}
					return
				case result := <-newPayload.OrderTradeUpdatePayload.ResultChan():
					r := result.Order
					CumQuoteQty := decimal.Zero
					avgPrice, err := decimal.NewFromString(r.AvgPrice)
					if err != nil {
						newSub.ErrChan() <- err
					}
					CumQuoteQty = avgPrice.Mul(decimal.RequireFromString(r.ExecutedQty))
					order := Order{
						Exchange:      BINANCE_NAME.String(),
						AccountType:   req.AccountType,
						Symbol:        r.Symbol,
						OrderId:       strconv.FormatInt(r.OrderId, 10),
						ClientOrderId: r.ClientOrderId,
						Price:         r.Price,
						Quantity:      r.OrigQty,
						ExecutedQty:   r.ExecutedQty,
						CumQuoteQty:   CumQuoteQty.String(),
						Status:        b.bnConverter.FromBNOrderStatus(r.Status),
						Type:          b.bnConverter.FromBNOrderType(r.Type),
						Side:          b.bnConverter.FromBNOrderSide(r.Side),
						PositionSide:  b.bnConverter.FromBNPositionSide(r.PositionSide),
						TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
						FeeAmount:     r.FeeQty,
						FeeCcy:        r.FeeAsset,
						ReduceOnly:    r.IsReduceOnly,
						CreateTime:    result.Timestamp,
						UpdateTime:    result.Timestamp,
					}
					newSub.resultChan <- order
				}
			}
		}()

		return newSub, nil
	default:
		return nil, ErrorAccountType
	}
}

func (b BinanceTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	var order *Order
	binance := &mybinanceapi.MyBinance{}
	var err error
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if b.wsSpotAccount == nil {
			b.wsSpotAccount, err = binance.NewSpotWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsSpotAccount.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		res, err := b.wsSpotAccount.CreateOrder(b.apiSpotOrderCreate(req))
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSpotOrderCreate(req, &res.Result)

	}
	return order, nil
}

func (b BinanceTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsCancelOrder(req *OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsCancelOrders(reqs []*OrderParam) error {
	//TODO implement me
	panic("implement me")
}
