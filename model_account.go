package mytrade

import "github.com/shopspring/decimal"

type AccountInfo struct {
}

type FeeRate struct {
	Maker decimal.Decimal
	Taker decimal.Decimal
}

type Position struct {
	Exchange               string       `json:"exchange"`               //交易所
	AccountType            string       `json:"accountType"`            //账户类型
	Symbol                 string       `json:"symbol"`                 //交易对
	InitialMargin          string       `json:"initialMargin"`          //当前所需起始保证金(基于最新标记价格)
	MaintMargin            string       `json:"maintMargin"`            //维持保证金
	UnrealizedProfit       string       `json:"unrealizedProfit"`       //持仓未实现盈亏
	PositionInitialMargin  string       `json:"positionInitialMargin"`  //持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin string       `json:"openOrderInitialMargin"` //当前挂单所需起始保证金(基于最新标记价格)
	Leverage               string       `json:"leverage"`               //杠杆倍数
	MarginMode             MarginMode   `json:"isolated"`               //是否逐仓
	EntryPrice             string       `json:"entryPrice"`             //开仓价格
	MaxNotional            string       `json:"maxNotional"`            //当前杠杆下用户可用的最大名义价值
	PositionSide           PositionSide `json:"positionSide"`           // 持仓方向
	PositionAmt            string       `json:"positionAmt"`            // 持仓数量
	MarkPrice              string       `json:"markPrice"`              // 标记价格
	LiquidationPrice       string       `json:"liquidationPrice"`       // 强平价格
	MarginRatio            string       `json:"marginRatio"`            // 保证金率
	UpdateTime             int64        `json:"updateTime"`             // 更新时间
}

type Asset struct {
	Exchange               string `json:"exchange"`               //交易所
	AccountType            string `json:"accountType"`            //账户类型
	Asset                  string `json:"asset"`                  //资产
	Free                   string `json:"free"`                   //可用余额
	Locked                 string `json:"locked"`                 //冻结余额
	WalletBalance          string `json:"walletBalance"`          //余额
	UnrealizedProfit       string `json:"unrealizedProfit"`       //未实现盈亏
	MarginBalance          string `json:"marginBalance"`          //保证金余额
	MaintMargin            string `json:"maintMargin"`            //维持保证金
	InitialMargin          string `json:"initialMargin"`          //当前所需起始保证金
	PositionInitialMargin  string `json:"positionInitialMargin"`  //持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` //当前挂单所需起始保证金(基于最新标记价格)
	CrossWalletBalance     string `json:"crossWalletBalance"`     //全仓账户余额
	CrossUnPnl             string `json:"crossUnPnl"`             //全仓持仓未实现盈亏
	AvailableBalance       string `json:"availableBalance"`       //可用余额
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`      //最大可转出余额
	MarginAvailable        bool   `json:"marginAvailable"`        //否可用作联合保证金
	UpdateTime             int64  `json:"updateTime"`             //更新时间
}
