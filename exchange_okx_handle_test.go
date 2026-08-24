package mytrade

import (
	"testing"

	"github.com/Hongssd/myokxapi"
)

func TestOkxOrderExecutedQty(t *testing.T) {
	t.Parallel()
	if got := okxOrderExecutedQty(""); got != "0" {
		t.Fatalf("empty accFillSz=%q want 0", got)
	}
	if got := okxOrderExecutedQty("2.56"); got != "2.56" {
		t.Fatalf("accFillSz=%q want 2.56", got)
	}
}

func TestHandleOrderFromQueryOrderGetUsesAccFillSz(t *testing.T) {
	t.Parallel()
	engine := &OkxTradeEngine{}
	res := &myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderGetRes]{
		OkxErrorRes: myokxapi.OkxErrorRes{Code: "0"},
		Data: myokxapi.PrivateRestTradeOrderGetRes{
			{
				InstType:  "SWAP",
				InstId:    "SKHY-USDT-SWAP",
				OrdId:     "1",
				ClOrdId:   "c1",
				Sz:        "8.21",
				FillSz:    "0.01",
				AccFillSz: "2.56",
				AvgPx:     "166.71",
				State:     OKX_ORDER_STATUS_CANCELED,
				OrdType:   OKX_ORDER_TYPE_POST_ONLY,
				Side:      OKX_ORDER_SIDE_SELL,
				PosSide:   "long",
			},
		},
	}
	order, err := engine.handleOrderFromQueryOrderGet(&QueryOrderParam{}, res)
	if err != nil {
		t.Fatal(err)
	}
	if order.ExecutedQty != "2.56" {
		t.Fatalf("ExecutedQty=%q want accFillSz 2.56 (not fillSz 0.01)", order.ExecutedQty)
	}
	if order.Status != ORDER_STATUS_CANCELED {
		t.Fatalf("Status=%s want CANCELED", order.Status)
	}
}

func TestHandleOrderFromWsOrderUsesAccFillSz(t *testing.T) {
	t.Parallel()
	engine := &OkxTradeEngine{}
	order := engine.handleOrderFromWsOrder(myokxapi.WsOrders{
		Orders: myokxapi.Orders{
			InstType:  "SWAP",
			InstId:    "SKHY-USDT-SWAP",
			OrdId:     "1",
			ClOrdId:   "c1",
			Sz:        "8.21",
			FillSz:    "0.01",
			AccFillSz: "2.56",
			AvgPx:     "166.71",
			State:     OKX_ORDER_STATUS_CANCELED,
			OrdType:   OKX_ORDER_TYPE_POST_ONLY,
			Side:      OKX_ORDER_SIDE_SELL,
			PosSide:   "long",
		},
	})
	if order.ExecutedQty != "2.56" {
		t.Fatalf("WS ExecutedQty=%q want accFillSz 2.56 (not fillSz 0.01)", order.ExecutedQty)
	}
}
