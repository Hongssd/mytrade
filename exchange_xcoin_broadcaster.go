package mytrade

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Hongssd/myxcoinapi"
)

// 广播 用于异步接收订单操作的返回结果
type xcoinOrderBroadcaster struct {
	accountType      string
	xcoinWsAccount   *myxcoinapi.PrivateWsStreamClient
	currentSubscribe *myxcoinapi.Subscription[myxcoinapi.WsOrder]
	subscribers      *MySyncMap[string, *xcoinOrderSubscriber]
	keys             *MySyncMap[*xcoinOrderSubscriber, string]
	mu               sync.RWMutex
}

// 订阅者 用于异步接收订单操作的返回结果
type xcoinOrderSubscriber struct {
	accountType   string
	symbol        string
	clientOrderId string
	orderId       string
	ch            *subscription[Order]
}

// 新建订阅者
func (x *XcoinTradeEngine) newOrderSubscriber(ob **xcoinOrderBroadcaster, clientOrderId, orderId, accountType, symbol string) (*xcoinOrderSubscriber, error) {
	sub := &xcoinOrderSubscriber{
		accountType:   accountType,
		symbol:        symbol,
		clientOrderId: clientOrderId,
		orderId:       orderId,
		ch: &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		},
	}
	if *ob == nil {
		newOb, err := x.newOrderBroadcaster(accountType, symbol)
		if err != nil {
			return nil, err
		}
		*ob = newOb
	}

	(*ob).mu.Lock()
	defer (*ob).mu.Unlock()

	key := clientOrderId
	if key == "" {
		key = orderId
	}
	(*ob).subscribers.Store(key, sub)
	(*ob).keys.Store(sub, key)
	return sub, nil
}

// 关闭订阅者
func (x *XcoinTradeEngine) closeSubscribe(b **xcoinOrderBroadcaster, sub *xcoinOrderSubscriber) {
	if b == nil || *b == nil || sub == nil {
		return
	}
	(*b).mu.Lock()
	defer (*b).mu.Unlock()
	key, _ := (*b).keys.Load(sub)
	(*b).subscribers.Delete(key)
	(*b).keys.Delete(sub)
}

// 等待广播消息，超时返回
func (x *XcoinTradeEngine) waitSubscribeReturn(sub *xcoinOrderSubscriber, timeout time.Duration) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	queryOrder := func() (*Order, error) {
		queryReq := x.NewQueryOrderReq().SetAccountType(sub.accountType).SetSymbol(sub.symbol)
		if sub.clientOrderId != "" {
			queryReq.SetClientOrderId(sub.clientOrderId)
		} else if sub.orderId != "" {
			queryReq.SetOrderId(sub.orderId)
		} else {
			return nil, errors.New("api msg timeout")
		}
		return x.QueryOrder(queryReq)
	}
	select {
	case <-ctx.Done():
		return queryOrder()
	case err := <-sub.ch.ErrChan():
		return nil, err
	case <-sub.ch.CloseChan():
		return queryOrder()
	case order := <-sub.ch.ResultChan():
		return &order, nil
	}
}

// 新建广播者，本质上是一条ws链接接收订单推送结果进行广播
func (x *XcoinTradeEngine) newOrderBroadcaster(accountType, symbol string) (*xcoinOrderBroadcaster, error) {
	wsClient := xcoin.NewPrivateWsStreamClient(xcoin.NewRestClient(x.apiKey, x.apiSecret))
	err := wsClient.OpenConn()
	if err != nil {
		return nil, err
	}

	broadcaster := &xcoinOrderBroadcaster{
		accountType:    accountType,
		xcoinWsAccount: wsClient,
		subscribers:    GetPointer(NewMySyncMap[string, *xcoinOrderSubscriber]()),
		keys:           GetPointer(NewMySyncMap[*xcoinOrderSubscriber, string]()),
		mu:             sync.RWMutex{},
	}

	sub, err := broadcaster.xcoinWsAccount.SubscribeOrder(accountType, symbol)
	if err != nil {
		return nil, err
	}
	broadcaster.currentSubscribe = sub

	go func() {
		for {
			sub := broadcaster.currentSubscribe
			select {
			case err := <-sub.ErrChan():
				broadcaster.subscribers.Range(func(key string, value *xcoinOrderSubscriber) bool {
					value.ch.ErrChan() <- err
					return true
				})
			case result := <-sub.ResultChan():
				order := x.handleOrderFromWsOrder(result)
				order.AccountType = broadcaster.accountType
				broadcaster.subscribers.Range(func(key string, value *xcoinOrderSubscriber) bool {
					if (value.clientOrderId != "" && order.ClientOrderId == value.clientOrderId) ||
						(value.orderId != "" && order.OrderId == value.orderId) ||
						(value.clientOrderId == "" && value.orderId == "") {
						value.ch.ResultChan() <- *order
					}
					return true
				})
			case <-sub.CloseChan():
				broadcaster.subscribers.Range(func(key string, value *xcoinOrderSubscriber) bool {
					value.ch.CloseChan() <- struct{}{}
					return true
				})
				return
			}
		}
	}()

	return broadcaster, nil
}

// 从 WS 订单数据转换为 Order
func (x *XcoinTradeEngine) handleOrderFromWsOrder(wsOrder myxcoinapi.WsOrder) *Order {
	fee := wsOrder.QuoteFee
	if fee == "" || fee == "0" {
		fee = wsOrder.BaseFee
	}
	orderType, timeInForce := x.xcoinConverter.FromXcoinOrderType(wsOrder.OrderType, wsOrder.TimeInForce)
	return &Order{
		Exchange:      x.ExchangeType().String(),
		AccountType:   wsOrder.BusinessType,
		Symbol:        wsOrder.Symbol,
		IsMargin:      false,
		IsIsolated:    false,
		OrderId:       wsOrder.OrderId,
		ClientOrderId: wsOrder.ClientOrderId,
		Price:         wsOrder.Price,
		Quantity:      wsOrder.Qty,
		ExecutedQty:   wsOrder.TotalFillQty,
		CumQuoteQty:   wsOrder.QuoteQty,
		AvgPrice:      wsOrder.AvgPrice,
		Status:        x.xcoinConverter.FromXcoinOrderStatus(wsOrder.Status),
		Type:          orderType,
		Side:          x.xcoinConverter.FromXcoinOrderSide(wsOrder.Side),
		PositionSide:  x.xcoinConverter.FromXcoinPositionSide(wsOrder.PosSide),
		TimeInForce:   timeInForce,
		FeeAmount:     fee,
		FeeCcy:        "",
		ReduceOnly:    wsOrder.ReduceOnly,
		CreateTime:    stringToInt64(wsOrder.CreateTime),
		UpdateTime:    stringToInt64(wsOrder.UpdateTime),
	}
}

func (x *XcoinTradeEngine) getBroadcastFromAccountType(accountType string) **xcoinOrderBroadcaster {
	switch XcoinAccountType(accountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT:
		return &x.broadcasterSpot
	case XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL:
		return &x.broadcasterLinearPerpetual
	case XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
		return &x.broadcasterLinearFutures
	default:
		return nil
	}
}
