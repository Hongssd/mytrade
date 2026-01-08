package mytrade

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Hongssd/mysunxapi"
)

// 广播 用于异步接收订单操作的返回结果
type sunxOrderBroadcaster struct {
	symbol           string
	accountType      string
	sunxWsAccount    *mysunxapi.PrivateWsStreamClient
	currentSubscribe *mysunxapi.Subscription[mysunxapi.WsOrdersReq, mysunxapi.WsOrders]
	subscribers      *MySyncMap[string, *sunxOrderSubscriber]
	keys             *MySyncMap[*sunxOrderSubscriber, string]
	mu               sync.RWMutex
}

// 订阅者 用于异步接收订单操作的返回结果
type sunxOrderSubscriber struct {
	symbol        string
	clientOrderId string
	orderId       string
	ch            *subscription[Order]
}

// 新建订阅者
func (s *SunxTradeEngine) newOrderSubscriber(broadcaster *sunxOrderBroadcaster, clientOrderId, orderId, symbol string) (*sunxOrderSubscriber, error) {

	sub := &sunxOrderSubscriber{
		symbol:        symbol,
		clientOrderId: clientOrderId,
		orderId:       orderId,
		ch: &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		},
	}

	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()

	// 使用 clientOrderId 或 orderId 作为 key
	key := clientOrderId
	if key == "" {
		key = orderId
	}
	broadcaster.subscribers.Store(key, sub)
	broadcaster.keys.Store(sub, key)
	return sub, nil
}

// 等待广播消息，超时返回
func (s *SunxTradeEngine) waitSubscribeReturn(sub *sunxOrderSubscriber, timeout time.Duration) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	//超时返回
	case <-ctx.Done():
		return nil, errors.New("api msg timeout")
	case <-sub.ch.CloseChan():
		//链接关闭的情况，使用接口查询
		queryReq := s.NewQueryOrderReq().SetSymbol(sub.symbol)
		if sub.clientOrderId != "" {
			queryReq.SetClientOrderId(sub.clientOrderId)
		} else if sub.orderId != "" {
			queryReq.SetOrderId(sub.orderId)
		}
		return s.QueryOrder(queryReq)
	case order := <-sub.ch.ResultChan():
		return &order, nil
	}
}

// 新建广播者，本质上是一条ws链接接收指定交易对的订单推送结果进行广播
func (s *SunxTradeEngine) newOrderBroadcaster(symbol string) (*sunxOrderBroadcaster, error) {
	// 创建新的 ws 连接
	wsClient := sunx.NewPrivateWsStreamClient(s.accessKey, s.secretKey, mysunxapi.WsAPITypeNotification)
	err := wsClient.OpenConn()
	if err != nil {
		return nil, err
	}

	broadcaster := &sunxOrderBroadcaster{
		symbol:        symbol,
		accountType:   SUNX_ACCOUNT_TYPE_SWAP.String(),
		sunxWsAccount: wsClient,
		subscribers:   GetPointer(NewMySyncMap[string, *sunxOrderSubscriber]()),
		keys:          GetPointer(NewMySyncMap[*sunxOrderSubscriber, string]()),
		mu:            sync.RWMutex{},
	}

	// 订阅指定交易对的订单推送
	sub, err := broadcaster.sunxWsAccount.SubscribeOrders([]string{symbol}, true)
	if err != nil {
		return nil, err
	}

	broadcaster.currentSubscribe = sub

	go func() {
		for {
			sub := broadcaster.currentSubscribe
			select {
			case err := <-sub.ErrChan():
				broadcaster.subscribers.Range(func(key string, value *sunxOrderSubscriber) bool {
					value.ch.ErrChan() <- err
					return true
				})
			case result := <-sub.ResultChan():
				// log.Infof("订单频道订阅接收到消息：%v", result)
				order := s.handleOrderFromWsOrder(result)
				order.AccountType = broadcaster.accountType
				broadcaster.subscribers.Range(func(key string, value *sunxOrderSubscriber) bool {
					// 根据 clientOrderId 或 orderId 匹配订阅者
					if (value.clientOrderId != "" && order.ClientOrderId == value.clientOrderId) ||
						(value.orderId != "" && order.OrderId == value.orderId) {
						value.ch.ResultChan() <- *order
					}
					return true
				})
			case <-sub.CloseChan():
				broadcaster.subscribers.Range(func(key string, value *sunxOrderSubscriber) bool {
					value.ch.CloseChan() <- struct{}{}
					return true
				})
				return
			}
		}
	}()

	return broadcaster, nil
}

// 关闭订阅者
func (s *SunxTradeEngine) closeSubscribe(broadcaster *sunxOrderBroadcaster, sub *sunxOrderSubscriber) {
	if broadcaster == nil {
		return
	}
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	key, _ := broadcaster.keys.Load(sub)
	broadcaster.subscribers.Delete(key)
	broadcaster.keys.Delete(sub)
}

// 从 WS 订单数据转换为 Order
func (s *SunxTradeEngine) handleOrderFromWsOrder(wsOrder mysunxapi.WsOrders) *Order {
	r := wsOrder.Data
	return &Order{
		Exchange:      s.ExchangeType().String(),
		AccountType:   SUNX_ACCOUNT_TYPE_SWAP.String(),
		Symbol:        wsOrder.ContractCode,
		IsMargin:      false,
		IsIsolated:    false,
		OrderId:       r.OrderId,
		ClientOrderId: r.ClientOrderId,
		Price:         r.Price,
		Quantity:      r.Volume,
		ExecutedQty:   r.TradeVolume,
		CumQuoteQty:   r.TradeTurnover,
		AvgPrice:      r.TradeAvgPrice,
		Status:        s.sunxConverter.FromSunxOrderStatus(r.State),
		Type:          s.sunxConverter.FromSunxOrderType(r.Type),
		Side:          s.sunxConverter.FromSunxOrderSide(r.Side),
		PositionSide:  s.sunxConverter.FromSunxPositionSide(r.PositionSide),
		TimeInForce:   s.sunxConverter.FromSunxTimeInForce(r.Type, r.TimeInForce),
		ReduceOnly:    r.ReduceOnly,
		FeeAmount:     r.Fee,
		FeeCcy:        r.FeeCurrency,
		CreateTime:    stringToInt64(r.CreatedTime),
		UpdateTime:    stringToInt64(r.UpdatedTime),
	}
}
