package mytrade

import (
	"fmt"
	"strings"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

type BitgetExchange struct {
	ExchangeBase
}

func isBitgetClassicMode(client *mybitgetapi.PrivateRestClient) (isClassic bool, modeDetectErr error) {
	_, errUta := client.NewPrivateRestUtaAccountSwitchStatus().Do()

	// 经典账户调用 UTA 开关状态接口会返回 40084，属于可预期结果，直接判定为 classic。
	if errUta != nil {
		if strings.Contains(errUta.Error(), "40084") {
			return true, nil
		}
		return false, fmt.Errorf("detect bitget account stack failed, uta error: %s", errUta.Error())
	}

	return false, nil
}

func isBitgetPosModeHedge(client *mybitgetapi.PrivateRestClient) (bool, error) {
	res, err := client.NewPrivateRestClassicFuturesAccountSingleAccount().
		ProductType(mybitgetapi.INST_TYPE_USDT_FUTURES.String()).
		Symbol("BTCUSDT").
		MarginCoin("USDT").
		Do()
	if err != nil {
		return false, err
	}
	return strings.EqualFold(res.Data.PosMode, BITGET_POSITION_MODE_HEDGE), nil
}

func (b *BitgetExchange) NewExchangeInfo() TradeExchangeInfo {
	return &BitgetExchangeInfo{
		isLoaded: false,
	}
}

func (b *BitgetExchange) NewMarketData() TradeMarketData {
	return &BitgetMarketData{}
}

func (b *BitgetExchange) NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine {
	client := mybitgetapi.NewRestClient(apiKey, secretKey, passphrase).PrivateRestClient()
	isClassic, modeErr := isBitgetClassicMode(client)
	posModeHedge := false
	if isClassic {
		var err error
		posModeHedge, err = isBitgetPosModeHedge(client)
		if err != nil && modeErr == nil {
			modeErr = err
		}
	}
	return &BitgetTradeEngine{
		ExchangeBase:  b.ExchangeBase,
		apiKey:        apiKey,
		secretKey:     secretKey,
		passphrase:    passphrase,
		privateClient: client,
		isClassic:     isClassic,
		posModeHedge:  posModeHedge,
		modeDetectErr: modeErr,
	}
}

func (b *BitgetExchange) NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount {
	client := mybitgetapi.NewRestClient(apiKey, secretKey, passphrase).PrivateRestClient()
	isClassic, modeErr := isBitgetClassicMode(client)
	return &BitgetTradeAccount{
		ExchangeBase:  b.ExchangeBase,
		converter:     BitgetEnumConverter{},
		apiKey:        apiKey,
		secretKey:     secretKey,
		passphrase:    passphrase,
		privateClient: client,
		isClassic:     isClassic,
		modeDetectErr: modeErr,
	}
}
