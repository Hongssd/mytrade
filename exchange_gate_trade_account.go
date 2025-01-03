package mytrade

import (
	"github.com/shopspring/decimal"
)

type GateTradeAccount struct {
	ExchangeBase

	gateConverter GateEnumConverter
	apiKey        string
	secretKey     string
	passphrase    string
}

func (o GateTradeAccount) GetAccountMode() (AccountMode, error) {
	return ACCOUNT_MODE_UNKNOWN, nil
}

func (o GateTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	return MARGIN_MODE_UNKNOWN, nil
}

func (o GateTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	return POSITION_MODE_UNKNOWN, nil
}

func (o GateTradeAccount) GetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (o GateTradeAccount) SetAccountMode(mode AccountMode) error {
	return nil
}

func (o GateTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	return nil
}

func (o GateTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	return nil
}

func (o GateTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {
	return nil
}

func (o GateTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	return nil, nil
}

func (o GateTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	return nil, nil
}

func (o GateTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	return nil, nil
}

func (o GateTradeAccount) AssetTransfer(req *AssetTransferParams) ([]*AssetTransfer, error) {
	return nil, nil
}

func (o GateTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	return nil, nil
}
