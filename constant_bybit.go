package mytrade

type BybitAccountType string

func (b BybitAccountType) String() string {
	return string(b)
}

const (
	BYBIT_AC_SPOT    BybitAccountType = "spot"    //现货
	BYBIT_AC_LINEAR  BybitAccountType = "linear"  //永续合约
	BYBIT_AC_INVERSE BybitAccountType = "inverse" //交割合约
	BYBIT_AC_OPTION  BybitAccountType = "option"  //期权
)

// 支持的訂單類型 (orderType):
// 限價單: orderType=Limit, 需要指定訂單數量和價格.
// 市價單: orderType=Market, 以Bybit市場內最優的價格一直執行到成交. 選擇市價單時，price 參數為空。在期貨交易系統，為了保護市價單產生嚴重的滑點，Bybit 交易系統會將市價單轉為限價單進行撮合，如果市場價格轉限價時，超過滑點設置的閾值，該筆市場價格訂單將會被取消。滑點閾值是指訂單價格偏離最新成交價格的百分比，目前閾值設置為：BTC 合約3%，其他合約5%。
const (
	BYBIT_ORDER_TYPE_MARKET = "Market"
	BYBIT_ORDER_TYPE_LIMIT  = "Limit"
)

const (
	BYBIT_ORDER_SIDE_BUY  = "Buy"
	BYBIT_ORDER_SIDE_SELL = "Sell"
)

// 倉位標識, 用戶不同倉位模式. 該字段對於雙向持倉模式(僅USDT永續和反向期貨有雙向模式)是必傳:
// 0: 單向持倉
// 1: 買側雙向持倉
// 2: 賣側雙向持倉
// 僅對linear和inverse有效
const (
	BYBIT_POSITION_SIDE_BOTH  = 0
	BYBIT_POSITION_SIDE_LONG  = 1
	BYBIT_POSITION_SIDE_SHORT = 2
)

// 订单状态
// New 訂單成功下達
// PartiallyFilled 部分成交
// Filled 完全成交
// Rejected 訂單被拒絕
// Cancelled 已撤销
const (
	BYBIT_ORDER_STATUS_NEW                       = "New"
	BYBIT_ORDER_STATUS_PARTIALLY_FILLED          = "PartiallyFilled"
	BYBIT_ORDER_STATUS_FILLED                    = "Filled"
	BYBIT_ORDER_STATUS_CANCELED                  = "Cancelled"
	BYBIT_ORDER_STATUS_REJECTED                  = "Rejected"
	BYBIT_ORDER_STATUS_PARTIALLY_FILLED_CANCELED = "PartiallyFilledCanceled"
)

// 支持的timeInForce策略:
// GTC 一直有效至取消
// IOC 立即成交或取消
// FOK 完全成交或取消
// PostOnly: 被動委托類型，如果該訂單在提交時會被立即執行成交，它將被取消. 這樣做的目的是為了保護您的訂單在提交的過程中，如果因為場內的價格變化，而撮合系統無法委託該筆訂單到訂單簿，因此會被取消。針對 PostOnly 訂單類型，單筆訂單可提交的數量比其他類型的訂單更多，請參考該接口中的lotSizeFilter > postOnlyMaxOrderQty參數.
const (
	BYBIT_TIME_IN_FORCE_GTC       = "GTC"
	BYBIT_TIME_IN_FORCE_IOC       = "IOC"
	BYBIT_TIME_IN_FORCE_FOK       = "FOK"
	BYBIT_TIME_IN_FORCE_POST_ONLY = "PostOnly"
)

// ISOLATED_MARGIN(逐倉保證金模式), REGULAR_MARGIN（全倉保證金模式）PORTFOLIO_MARGIN（組合保證金模式）默認常規，傳常規則返回設置成功
const (
	BYBIT_ACCOUNT_MOED_ISOLATED_MARGIN  = "ISOLATED_MARGIN"
	BYBIT_ACCOUNT_MODE_REGULAR_MARGIN   = "REGULAR_MARGIN"
	BYBIT_ACCOUNT_MODE_PORTFOLIO_MARGIN = "PORTFOLIO_MARGIN"
)

// 0: 全倉, 1: 逐倉
const (
	BYBIT_MARGIN_MODE_ISOLATED = 1
	BYBIT_MARGIN_MODE_CROSSED  = 0
)

// mode	true	integer	倉位模式. 0: 單向持倉. 3: 雙向持倉
const (
	BYBIT_POSITION_MODE_HEDGE  = 3
	BYBIT_POSITION_MODE_ONEWAY = 0
)

// 時間粒度. 1,3,5,15,30,60,120,240,360,720,D,M,W
const (
	BYBIT_KLINE_INTERVAL_1m  = "1"
	BYBIT_KLINE_INTERVAL_3m  = "3"
	BYBIT_KLINE_INTERVAL_5m  = "5"
	BYBIT_KLINE_INTERVAL_15m = "15"
	BYBIT_KLINE_INTERVAL_30m = "30"
	BYBIT_KLINE_INTERVAL_1H  = "60"
	BYBIT_KLINE_INTERVAL_2H  = "120"
	BYBIT_KLINE_INTERVAL_4H  = "240"
	BYBIT_KLINE_INTERVAL_6H  = "360"
	BYBIT_KLINE_INTERVAL_12H = "720"
	BYBIT_KLINE_INTERVAL_1D  = "D"
	BYBIT_KLINE_INTERVAL_1W  = "W"
	BYBIT_KLINE_INTERVAL_1M  = "M"
)

func bybitGetMillisecondFromInterval(interval string) int64 {
	switch interval {
	case BYBIT_KLINE_INTERVAL_1m:
		return 60 * 1000
	case BYBIT_KLINE_INTERVAL_3m:
		return 3 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_5m:
		return 5 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_15m:
		return 15 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_30m:
		return 30 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1H:
		return 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_2H:
		return 2 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_4H:
		return 4 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_6H:
		return 6 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_12H:
		return 12 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1D:
		return 24 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1W:
		return 7 * 24 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1M:
		return 30 * 24 * 60 * 60 * 1000
	default:
		return 60 * 1000
	}
}
func bybitGetKlineCloseTime(ts int64, interval string) int64 {
	return ts + bybitGetMillisecondFromInterval(interval) - 1
}
