package mytrade

import (
	"errors"
	"fmt"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

type BitgetExchange struct {
	ExchangeBase
}

func isBitgetClassicMode(client *mybitgetapi.PrivateRestClient) (isClassic bool, modeDetectErr error) {
	utaRes, errUta := client.NewPrivateRestUtaAccountSwitchStatus().Do()
	claRes, errClassic := client.NewPrivateRestClassicSpotAccountUpgradeStatus().Do()

	if errUta != nil && errClassic != nil {
		return false, errors.New("detect bitget account stack failed")
	}
	if errUta == nil && errClassic != nil {
		return false, nil
	}
	if errUta != nil && errClassic == nil {
		return true, nil
	}

	if utaRes != nil && utaRes.Data.Status == "success" {
		return false, nil
	}
	if claRes != nil && claRes.Data.Status == "success" {
		return true, nil
	}
	utaSt, claSt := "", ""
	if utaRes != nil {
		utaSt = utaRes.Data.Status
	}
	if claRes != nil {
		claSt = claRes.Data.Status
	}
	return false, fmt.Errorf("ambiguous bitget mode: utaStatus=%q classicStatus=%q", utaSt, claSt)
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
	return &BitgetTradeEngine{
		ExchangeBase:  b.ExchangeBase,
		apiKey:        apiKey,
		secretKey:     secretKey,
		passphrase:    passphrase,
		privateClient: client,
		isClassic:     isClassic,
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
