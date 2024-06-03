package mytrade

import "github.com/Hongssd/mybybitapi"

func (b *BybitTradeEngine) handleOrderFromWsOrder(orders mybybitapi.WsOrder) []*Order {
	// 从ws订单信息转换为本地订单信息
	var res []*Order
	for _, order := range orders.Data {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			AccountType:   order.Category,
			Symbol:        order.Symbol,
			OrderId:       order.OrderId,
			ClientOrderId: order.OrderLinkId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.CumExecQty,
			CumQuoteQty:   order.CumExecValue,
			Status:        b.bybitConverter.FromBYBITOrderStatus(order.OrderStatus),
			Type:          b.bybitConverter.FromBYBITOrderType(order.OrderType),
			Side:          b.bybitConverter.FromBYBITOrderSide(order.Side),
			PositionSide:  b.bybitConverter.FromBYBITPositionSide(order.PositionIdx),
			TimeInForce:   b.bybitConverter.FromBYBITTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),

			FeeAmount: order.CumExecFee,
			FeeCcy:    order.FeeCurrency,
		}
		res = append(res, order)
	}
	return res
}
