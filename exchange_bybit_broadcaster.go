package mytrade

import (
	"context"
	"errors"
	"github.com/Hongssd/mybybitapi"
	"sync"
	"time"
)

// 广播 用于异步接收订单操作的返回结果
type bybitOrderBroadcaster struct {
	accountType      string
	bybitWsAccount   *mybybitapi.PrivateWsStreamClient
	currentSubscribe *mybybitapi.Subscription[mybybitapi.WsOrder]
	subscribers      *MySyncMap[string, *bybitOrderSubscriber]
	keys             *MySyncMap[*bybitOrderSubscriber, string]
	mu               sync.RWMutex
}

// 订阅者 用于异步接收订单操作的返回结果
type bybitOrderSubscriber struct {
	symbol        string
	clientOrderId string
	ch            *subscription[Order]
}

// 新建订阅者
func (o *BybitTradeEngine) newOrderSubscriber(ob **bybitOrderBroadcaster, clientOrderId, accountType, symbol string) (*bybitOrderSubscriber, error) {

	sub := &bybitOrderSubscriber{
		symbol:        symbol,
		clientOrderId: clientOrderId,
		ch: &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		},
	}
	if *ob == nil {
		newOb, err := o.newOrderBroadcaster(accountType)
		if err != nil {
			return nil, err
		}
		*ob = newOb
	}

	(*ob).mu.Lock()
	defer (*ob).mu.Unlock()

	key := clientOrderId
	(*ob).subscribers.Store(key, sub)
	(*ob).keys.Store(sub, key)
	return sub, nil
}

// 关闭订阅者
func (o *BybitTradeEngine) closeSubscribe(b **bybitOrderBroadcaster, sub *bybitOrderSubscriber) {
	(*b).mu.Lock()
	defer (*b).mu.Unlock()
	key, _ := (*b).keys.Load(sub)
	(*b).subscribers.Delete(key)
	(*b).keys.Delete(sub)
}

// 等待广播消息，超时返回
func (o *BybitTradeEngine) waitSubscribeReturn(sub *bybitOrderSubscriber, timeout time.Duration) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	//超时返回
	case <-ctx.Done():
		return nil, errors.New("api msg timeout")
	case <-sub.ch.CloseChan():
		//链接关闭的情况，使用接口查询
		return o.QueryOrder(o.NewQueryOrderReq().SetSymbol(sub.symbol).SetClientOrderId(sub.clientOrderId))
	case order := <-sub.ch.ResultChan():
		return &order, nil
	}
}

// 新建广播者，本质上是一条ws链接接收所有的订单推送结果进行广播
func (o *BybitTradeEngine) newOrderBroadcaster(accountType string) (*bybitOrderBroadcaster, error) {
	broadcaster := &bybitOrderBroadcaster{
		accountType:    accountType,
		bybitWsAccount: mybybitapi.NewPrivateWsStreamClient(),
		subscribers:    GetPointer(NewMySyncMap[string, *bybitOrderSubscriber]()),
		keys:           GetPointer(NewMySyncMap[*bybitOrderSubscriber, string]()),
		mu:             sync.RWMutex{},
	}
	err := broadcaster.bybitWsAccount.OpenConn()
	if err != nil {
		return nil, err
	}

	err = broadcaster.bybitWsAccount.Auth(mybybitapi.NewRestClient(o.apiKey, o.secretKey))
	if err != nil {
		return nil, err
	}

	sub, err := broadcaster.bybitWsAccount.SubscribeOrder(accountType)
	if err != nil {
		return nil, err
	}

	broadcaster.currentSubscribe = sub

	go func() {
		for {
			sub := broadcaster.currentSubscribe
			select {
			case err := <-sub.ErrChan():
				broadcaster.subscribers.Range(func(key string, value *bybitOrderSubscriber) bool {
					value.ch.ErrChan() <- err
					return true
				})
			case result := <-sub.ResultChan():
				//log.Infof("订单频道订阅接收到消息：%s", result)
				orders := o.handleOrderFromWsOrder(result)
				for _, order := range orders {
					order.AccountType = broadcaster.accountType
					broadcaster.subscribers.Range(func(key string, value *bybitOrderSubscriber) bool {
						if value.clientOrderId == "" || order.ClientOrderId == value.clientOrderId {
							value.ch.ResultChan() <- *order
						}
						return true
					})
				}
			case <-sub.CloseChan():
				broadcaster.subscribers.Range(func(key string, value *bybitOrderSubscriber) bool {
					value.ch.CloseChan() <- struct{}{}
					return true
				})
				return
			}
		}
	}()

	return broadcaster, nil
}

func (o *BybitTradeEngine) getBroadcastFromAccountType(accountType string) **bybitOrderBroadcaster {
	switch BybitAccountType(accountType) {
	case BYBIT_AC_SPOT:
		return &o.broadcasterSpot
	case BYBIT_AC_LINEAR:
		return &o.broadcasterLinear
	case BYBIT_AC_INVERSE:
		return &o.broadcasterInverse
	default:
		return nil
	}
}
