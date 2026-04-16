package mytrade

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

// 广播：用于异步接收 Bitget 私有 orders 频道推送并分发给订阅者
type bitgetOrderBroadcaster struct {
	accountType      string
	wsAccount        *mybitgetapi.PrivateWsStreamClient
	currentSubscribe *mybitgetapi.Subscription[mybitgetapi.WsOrder]
	subscribers      *MySyncMap[string, *bitgetOrderSubscriber]
	keys             *MySyncMap[*bitgetOrderSubscriber, string]
	mu               sync.RWMutex
}

// 订阅者：用于异步接收订单推送
type bitgetOrderSubscriber struct {
	accountType   string
	symbol        string
	clientOrderId string
	ch            *subscription[Order]
}

func bitgetWsInstTypeFromAccountType(accountType string) (mybitgetapi.InstType, error) {
	switch accountType {
	case BITGET_AC_SPOT:
		return mybitgetapi.INST_TYPE_SPOT, nil
	case BITGET_AC_MARGIN:
		return mybitgetapi.INST_TYPE_MARGIN, nil
	case BITGET_AC_USDT_FUTURES:
		return mybitgetapi.INST_TYPE_USDT_FUTURES, nil
	case BITGET_AC_COIN_FUTURES:
		return mybitgetapi.INST_TYPE_COIN_FUTURES, nil
	case BITGET_AC_USDC_FUTURES:
		return mybitgetapi.INST_TYPE_USDC_FUTURES, nil
	default:
		return "", ErrorAccountType
	}
}

func (o *BitgetTradeEngine) getBroadcastFromAccountType(accountType string) **bitgetOrderBroadcaster {
	switch BitgetAccountType(accountType) {
	case BITGET_AC_SPOT:
		return &o.broadcasterSpot
	case BITGET_AC_MARGIN:
		return &o.broadcasterMarginCrossed
	case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
		return &o.broadcasterFutures
	default:
		return nil
	}
}

// 新建订阅者
func (o *BitgetTradeEngine) newOrderSubscriber(ob **bitgetOrderBroadcaster, clientOrderId, accountType, symbol string) (*bitgetOrderSubscriber, error) {
	sub := &bitgetOrderSubscriber{
		accountType:   accountType,
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

// 等待广播消息，超时返回
func (o *BitgetTradeEngine) waitSubscribeReturn(sub *bitgetOrderSubscriber, timeout time.Duration) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	case <-ctx.Done():
		return nil, errors.New("api msg timeout")
	case <-sub.ch.CloseChan():
		// 连接关闭时退化为 REST 查询（若缺字段则由上层处理）
		return o.QueryOrder(o.NewQueryOrderReq().SetAccountType(sub.accountType).SetSymbol(sub.symbol).SetClientOrderId(sub.clientOrderId))
	case order := <-sub.ch.ResultChan():
		return &order, nil
	}
}

// 新建广播者：一条私有 WS 连接消费 orders 频道并广播
func (o *BitgetTradeEngine) newOrderBroadcaster(accountType string) (*bitgetOrderBroadcaster, error) {
	apiType := mybitgetapi.WS_UTA
	if o.isClassic {
		apiType = mybitgetapi.WS_CLASSIC
	}
	broadcaster := &bitgetOrderBroadcaster{
		accountType: accountType,
		wsAccount:   bitget.NewPrivateWsStreamClient(apiType),
		subscribers: GetPointer(NewMySyncMap[string, *bitgetOrderSubscriber]()),
		keys:        GetPointer(NewMySyncMap[*bitgetOrderSubscriber, string]()),
		mu:          sync.RWMutex{},
	}

	if err := broadcaster.wsAccount.OpenConn(); err != nil {
		return nil, err
	}
	if err := broadcaster.wsAccount.Login(mybitgetapi.NewRestClient(o.apiKey, o.secretKey, o.passphrase)); err != nil {
		return nil, err
	}

	instType, err := bitgetWsInstTypeFromAccountType(accountType)
	if err != nil {
		return nil, err
	}
	sub, err := broadcaster.wsAccount.SubscribeOrders(instType, "")
	if err != nil {
		return nil, err
	}
	broadcaster.currentSubscribe = sub

	go func() {
		for {
			sub := broadcaster.currentSubscribe
			select {
			case err := <-sub.ErrChan():
				broadcaster.subscribers.Range(func(_ string, v *bitgetOrderSubscriber) bool {
					v.ch.ErrChan() <- err
					return true
				})
			case result := <-sub.ResultChan():
				for i := range result.Data {
					order := o.handleOrderFromBitgetWsOrderData(accountType, &result.Data[i])
					if order == nil {
						continue
					}
					broadcaster.subscribers.Range(func(_ string, v *bitgetOrderSubscriber) bool {
						if v.clientOrderId == "" || order.ClientOrderId == v.clientOrderId {
							v.ch.ResultChan() <- *order
						}
						return true
					})
				}
			case <-sub.CloseChan():
				broadcaster.subscribers.Range(func(_ string, v *bitgetOrderSubscriber) bool {
					v.ch.CloseChan() <- struct{}{}
					return true
				})
				return
			}
		}
	}()

	return broadcaster, nil
}

func (o *BitgetTradeEngine) closeSubscribe(b **bitgetOrderBroadcaster, sub *bitgetOrderSubscriber) {
	(*b).mu.Lock()
	defer (*b).mu.Unlock()
	key, _ := (*b).keys.Load(sub)
	(*b).subscribers.Delete(key)
	(*b).keys.Delete(sub)
}

func (o *BitgetTradeEngine) handleOrderFromBitgetWsOrderData(accountType string, d *mybitgetapi.WsOrderData) *Order {
	if d == nil {
		return nil
	}
	orderType := d.OrderType
	if orderType == "" {
		orderType = d.OrdType
	}
	ot, tif := o.converter.FromBitgetOrderTypeWithTIF(orderType, d.Force)

	feeAmt := d.FillFee
	feeCcy := d.FillFeeCoin
	if feeAmt == "" && len(d.FeeDetail) > 0 {
		feeAmt = d.FeeDetail[0].Fee
		feeCcy = d.FeeDetail[0].FeeCoin
	}

	return &Order{
		Exchange:      BITGET_NAME.String(),
		AccountType:   accountType,
		Symbol:        d.InstId,
		OrderId:       d.OrderId,
		ClientOrderId: d.ClientOid,
		Price:         d.Price,
		Quantity:      d.Size,
		ExecutedQty:   d.AccBaseVolume,
		CumQuoteQty:   d.QuoteVolume,
		AvgPrice:      d.PriceAvg,
		Status:        o.converter.FromBitgetOrderStatusUTA(d.Status),
		Type:          ot,
		Side:          o.converter.FromBitgetOrderSide(d.Side),
		PositionSide:  o.converter.FromBitgetPositionSide(d.PosSide),
		TimeInForce:   tif,
		FeeAmount:     feeAmt,
		FeeCcy:        feeCcy,
		ReduceOnly:    strings.EqualFold(d.ReduceOnly, BITGET_REDUCE_ONLY_YES),
		CreateTime:    stringToInt64(d.CTime),
		UpdateTime:    stringToInt64(d.UTime),
		RealizedPnl:   firstNonEmpty(d.Pnl, d.TotalProfits),
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
