package mytrade

import (
	mybitgetapi "github.com/Hongssd/mybitgetapi"
	"github.com/shopspring/decimal"
)

type BitgetTradeAccount struct {
	ExchangeBase

	converter BitgetEnumConverter

	apiKey        string
	secretKey     string
	passphrase    string
	privateClient *mybitgetapi.PrivateRestClient
	isClassic     bool
	modeDetectErr error
}

func (a *BitgetTradeAccount) checkMode() error {
	if a.modeDetectErr != nil {
		return a.modeDetectErr
	}
	return nil
}

func (a *BitgetTradeAccount) GetAccountMode() (AccountMode, error) {
	if err := a.checkMode(); err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	if a.isClassic {
		res, err := a.privateClient.NewPrivateRestClassicFuturesAccountSingleAccount().
			ProductType(mybitgetapi.INST_TYPE_USDT_FUTURES.String()).
			Symbol("BTCUSDT").
			MarginCoin("USDT").
			Do()
		if err != nil {
			return ACCOUNT_MODE_UNKNOWN, err
		}
		return a.converter.FromBitgetAccountMode(*res.Data.AssetMode), nil
	} else {
		// TODO UTA
	}
	return ACCOUNT_MODE_UNKNOWN, ErrorNotSupport
}

func (a *BitgetTradeAccount) GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	if err := a.checkMode(); err != nil {
		return MARGIN_MODE_UNKNOWN, err
	}
	if a.isClassic {
		marginCoin := bitgetMarginCoinFromSymbol(symbol, "")
		if marginCoin == "" {
			return MARGIN_MODE_UNKNOWN, ErrorSymbolNotFound
		}
		res, err := a.privateClient.NewPrivateRestClassicFuturesAccountSingleAccount().
			ProductType(mybitgetapi.INST_TYPE_USDT_FUTURES.String()).
			Symbol(symbol).
			MarginCoin(marginCoin).
			Do()
		if err != nil {
			return MARGIN_MODE_UNKNOWN, err
		}
		return a.converter.FromBitgetMarginMode(res.Data.MarginMode), nil
	} else {
		// TODO UTA
	}
	return MARGIN_MODE_UNKNOWN, ErrorNotSupport
}

func (a *BitgetTradeAccount) GetPositionMode(accountType, symbol string) (PositionMode, error) {
	if err := a.checkMode(); err != nil {
		return POSITION_MODE_UNKNOWN, err
	}
	if a.isClassic {
		switch accountType {
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			marginCoin := bitgetMarginCoinFromSymbol(symbol, "")
			if marginCoin == "" {
				return POSITION_MODE_UNKNOWN, ErrorSymbolNotFound
			}
			res, err := a.privateClient.NewPrivateRestClassicFuturesAccountSingleAccount().
				ProductType(accountType).
				Symbol(symbol).
				MarginCoin(marginCoin).
				Do()
			if err != nil {
				return POSITION_MODE_UNKNOWN, err
			}
			return a.converter.FromBitgetPositionMode(res.Data.PosMode), nil
		}
	} else {
		// TODO UTA
	}
	return POSITION_MODE_UNKNOWN, ErrorNotSupport
}

func (a *BitgetTradeAccount) GetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	if err := a.checkMode(); err != nil {
		return decimal.Zero, err
	}
	if a.isClassic {
		marginCoin := bitgetMarginCoinFromSymbol(symbol, "")
		if marginCoin == "" {
			return decimal.Zero, ErrorSymbolNotFound
		}
		switch accountType {
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			res, err := a.privateClient.NewPrivateRestClassicFuturesAccountSingleAccount().
				ProductType(accountType).
				Symbol(symbol).
				MarginCoin(marginCoin).
				Do()
			if err != nil {
				return decimal.Zero, err
			}
			switch marginMode {
			case MARGIN_MODE_CROSSED:
				return decimal.NewFromInt(int64(res.Data.CrossedMarginLeverage)), nil
			case MARGIN_MODE_ISOLATED:
				switch positionSide {
				case POSITION_SIDE_LONG:
					return decimal.NewFromInt(int64(res.Data.IsolatedLongLever)), nil
				case POSITION_SIDE_SHORT:
					return decimal.NewFromInt(int64(res.Data.IsolatedShortLever)), nil
				}
			}
		}
		return decimal.Zero, ErrorNotSupport
	} else {
		// TODO UTA
	}
	return decimal.Zero, ErrorNotSupport
}

func (a *BitgetTradeAccount) GetFeeRate(accountType, symbol string) (*FeeRate, error) {
	if err := a.checkMode(); err != nil {
		return nil, err
	}
	if a.isClassic {
		api := a.privateClient.NewPrivateRestClassicTradeRate().Symbol(symbol)
		switch accountType {
		case BITGET_AC_SPOT:
			api.BusinessType("spot")
		case BITGET_AC_MARGIN:
			api.BusinessType("margin")
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			api.BusinessType("mix")
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		feeRate := &FeeRate{
			Maker: decimal.RequireFromString(res.Data.MakerFeeRate),
			Taker: decimal.RequireFromString(res.Data.TakerFeeRate),
		}
		return feeRate, nil
	}
	return nil, ErrorNotSupport
}

func (a *BitgetTradeAccount) GetPositions(accountType string, symbols ...string) ([]*Position, error) {
	if err := a.checkMode(); err != nil {
		return nil, err
	}
	if a.isClassic {
		switch accountType {
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			res, err := a.privateClient.NewPrivateRestClassicFuturesPositionAllPosition().
				ProductType(accountType).
				Do()
			if err != nil {
				return nil, err
			}
			positions := make([]*Position, 0, len(res.Data))
			for _, p := range res.Data {
				if len(symbols) > 0 && !stringInSlice(p.Symbol, symbols) {
					continue
				}
				position := &Position{
					Exchange:               BITGET_NAME.String(),
					AccountType:            accountType,
					Symbol:                 p.Symbol,
					MarginCcy:              p.MarginCoin,
					InitialMargin:          p.MarginSize,
					MaintMargin:            p.KeepMarginRate,
					UnrealizedProfit:       p.UnrealizedPL,
					PositionInitialMargin:  p.MarginSize,
					OpenOrderInitialMargin: "0",
					Leverage:               p.Leverage,
					MarginMode:             a.converter.FromBitgetMarginMode(p.MarginMode),
					EntryPrice:             p.OpenPriceAvg,
					MaxNotional:            "0",
					PositionSide:           a.converter.FromBitgetPositionSide(p.HoldSide),
					PositionAmt:            p.Total,
					MarkPrice:              p.MarkPrice,
					LiquidationPrice:       p.LiquidationPrice,
					MarginRatio:            p.MarginRatio,
					UpdateTime:             stringToInt64(p.UTime),
				}
				positions = append(positions, position)
			}
			return positions, nil
		}
	}
	return nil, ErrorNotSupport
}

func (a *BitgetTradeAccount) GetAssets(accountType string, currencies ...string) ([]*Asset, error) {
	if err := a.checkMode(); err != nil {
		return nil, err
	}
	if a.isClassic {
		switch accountType {
		case BITGET_AC_SPOT:
			res, err := a.privateClient.NewPrivateRestClassicSpotAccountAssets().AssetType("all").Do()
			if err != nil {
				return nil, err
			}
			assets := make([]*Asset, 0, len(res.Data))
			for _, a := range res.Data {
				if len(currencies) > 0 && !stringInSlice(a.Coin, currencies) {
					continue
				}
				available, err := decimal.NewFromString(a.Available)
				if err != nil {
					return nil, err
				}
				frozen, err := decimal.NewFromString(a.Frozen)
				if err != nil {
					return nil, err
				}
				locked, err := decimal.NewFromString(a.Locked)
				if err != nil {
					return nil, err
				}
				totalLocked := frozen.Add(locked)
				walletBalance := available.Add(totalLocked)
				asset := &Asset{
					Exchange:               BITGET_NAME.String(),
					AccountType:            accountType,
					Asset:                  a.Coin,
					Free:                   a.Available,
					Locked:                 totalLocked.String(),
					WalletBalance:          walletBalance.String(),
					UnrealizedProfit:       "0",
					MarginBalance:          "0",
					MaintMargin:            "0",
					InitialMargin:          "0",
					PositionInitialMargin:  "0",
					OpenOrderInitialMargin: "0",
					CrossWalletBalance:     "0",
					CrossUnPnl:             "0",
					AvailableBalance:       a.Available,
					MaxWithdrawAmount:      a.Available,
					MarginAvailable:        false,
					UpdateTime:             stringToInt64(a.UTime),
				}
				assets = append(assets, asset)
			}
			return assets, nil
		case BITGET_AC_MARGIN_CROSSED:
			res, err := a.privateClient.NewPrivateRestClassicMarginCrossAccountAssets().Do()
			if err != nil {
				return nil, err
			}
			assets := make([]*Asset, 0, len(res.Data))
			for _, a := range res.Data {
				if len(currencies) > 0 && !stringInSlice(a.Coin, currencies) {
					continue
				}
				available, err := decimal.NewFromString(a.Available)
				if err != nil {
					return nil, err
				}
				frozen, err := decimal.NewFromString(a.Frozen)
				if err != nil {
					return nil, err
				}
				walletBalance := available.Add(frozen)
				asset := &Asset{
					Exchange:               BITGET_NAME.String(),
					AccountType:            accountType,
					Asset:                  a.Coin,
					Borrowed:               a.Borrow,
					Interest:               a.Interest,
					Free:                   a.Available,
					Locked:                 a.Frozen,
					WalletBalance:          walletBalance.String(),
					UnrealizedProfit:       "0",
					MarginBalance:          a.Net,
					MaintMargin:            "0",
					InitialMargin:          "0",
					PositionInitialMargin:  "0",
					OpenOrderInitialMargin: "0",
					CrossWalletBalance:     a.Net,
					CrossUnPnl:             "0",
					AvailableBalance:       a.Available,
					MaxWithdrawAmount:      a.Available,
					MarginAvailable:        false,
					UpdateTime:             stringToInt64(a.UTime),
				}
				assets = append(assets, asset)
			}
			return assets, nil
		case BITGET_AC_MARGIN_ISOLATED:
			res, err := a.privateClient.NewPrivateRestClassicMarginIsolatedAccountAssets().Do()
			if err != nil {
				return nil, err
			}
			assets := make([]*Asset, 0, len(res.Data))
			for _, a := range res.Data {
				if len(currencies) > 0 && !stringInSlice(a.Coin, currencies) {
					continue
				}
				available, err := decimal.NewFromString(a.Available)
				if err != nil {
					return nil, err
				}
				frozen, err := decimal.NewFromString(a.Frozen)
				if err != nil {
					return nil, err
				}
				walletBalance := available.Add(frozen)
				asset := &Asset{
					Exchange:               BITGET_NAME.String(),
					AccountType:            accountType,
					Asset:                  a.Coin,
					Borrowed:               a.Borrow,
					Interest:               a.Interest,
					Free:                   a.Available,
					Locked:                 a.Frozen,
					WalletBalance:          walletBalance.String(),
					UnrealizedProfit:       "0",
					MarginBalance:          a.Net,
					MaintMargin:            "0",
					InitialMargin:          "0",
					PositionInitialMargin:  "0",
					OpenOrderInitialMargin: "0",
					CrossWalletBalance:     "0",
					CrossUnPnl:             "0",
					AvailableBalance:       a.Available,
					MaxWithdrawAmount:      a.Available,
					MarginAvailable:        false,
					UpdateTime:             stringToInt64(a.UTime),
				}
				assets = append(assets, asset)
			}
			return assets, nil
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			res, err := a.privateClient.NewPrivateRestClassicFuturesAccountList().ProductType(accountType).Do()
			if err != nil {
				return nil, err
			}
			assets := make([]*Asset, 0, len(res.Data))
			for _, row := range res.Data {
				// 联合保证金模式优先使用 assetList 按币种返回
				if len(row.AssetList) > 0 {
					for _, item := range row.AssetList {
						if len(currencies) > 0 && !stringInSlice(item.Coin, currencies) {
							continue
						}
						balance, err := decimal.NewFromString(item.Balance)
						if err != nil {
							return nil, err
						}
						available, err := decimal.NewFromString(item.Available)
						if err != nil {
							return nil, err
						}
						locked := balance.Sub(available)
						if locked.IsNegative() {
							locked = decimal.Zero
						}
						assets = append(assets, &Asset{
							Exchange:               BITGET_NAME.String(),
							AccountType:            accountType,
							Asset:                  item.Coin,
							Free:                   item.Available,
							Locked:                 locked.String(),
							WalletBalance:          item.Balance,
							UnrealizedProfit:       row.UnrealizedPL,
							MarginBalance:          item.Balance,
							MaintMargin:            row.CrossedRiskRate,
							InitialMargin:          "0",
							PositionInitialMargin:  "0",
							OpenOrderInitialMargin: "0",
							CrossWalletBalance:     item.Balance,
							CrossUnPnl:             row.UnrealizedPL,
							AvailableBalance:       item.Available,
							MaxWithdrawAmount:      item.Available,
							MarginAvailable:        false,
							UpdateTime:             0,
						})
					}
					continue
				}

				coin := row.MarginCoin
				if len(currencies) > 0 && !stringInSlice(coin, currencies) {
					continue
				}
				crossUnPnl := "0"
				if row.CrossedUnrealizedPL != nil {
					crossUnPnl = *row.CrossedUnrealizedPL
				}
				assets = append(assets, &Asset{
					Exchange:               BITGET_NAME.String(),
					AccountType:            accountType,
					Asset:                  coin,
					Free:                   row.Available,
					Locked:                 row.Locked,
					WalletBalance:          row.AccountEquity,
					UnrealizedProfit:       row.UnrealizedPL,
					MarginBalance:          row.AccountEquity,
					MaintMargin:            row.CrossedRiskRate,
					InitialMargin:          "0",
					PositionInitialMargin:  "0",
					OpenOrderInitialMargin: "0",
					CrossWalletBalance:     row.AccountEquity,
					CrossUnPnl:             crossUnPnl,
					AvailableBalance:       row.Available,
					MaxWithdrawAmount:      row.MaxTransferOut,
					MarginAvailable:        false,
					UpdateTime:             0,
				})
			}
			return assets, nil
		}
	}
	return nil, ErrorNotSupport
}

func (a *BitgetTradeAccount) SetAccountMode(mode AccountMode) error {
	if err := a.checkMode(); err != nil {
		return err
	}
	return ErrorNotSupport
}

func (a *BitgetTradeAccount) SetMarginMode(accountType, symbol string, mode MarginMode) error {
	if err := a.checkMode(); err != nil {
		return err
	}
	if a.isClassic {
		marginCoin := bitgetMarginCoinFromSymbol(symbol, "")
		if marginCoin == "" {
			return ErrorSymbolNotFound
		}
		switch accountType {
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			_, err := a.privateClient.NewPrivateRestClassicFuturesAccountSetMarginMode().
				ProductType(accountType).
				Symbol(symbol).
				MarginMode(a.converter.ToBitgetMarginMode(mode)).
				MarginCoin(marginCoin).
				Do()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrorNotSupport
}

func (a *BitgetTradeAccount) SetPositionMode(accountType, symbol string, mode PositionMode) error {
	if err := a.checkMode(); err != nil {
		return err
	}
	if a.isClassic {
		switch accountType {
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			_, err := a.privateClient.NewPrivateRestClassicFuturesAccountSetPositionMode().
				ProductType(accountType).
				PosMode(a.converter.ToBitgetPositionMode(mode)).
				Do()
			if err != nil {
				return err
			}
			return nil
		}
	}

	return ErrorNotSupport
}

func (a *BitgetTradeAccount) SetLeverage(accountType, symbol string, marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {
	if err := a.checkMode(); err != nil {
		return err
	}
	// MARGIN 现货杠杆不支持调整杠杆倍数
	if a.isClassic {
		switch accountType {
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			posMode, err := a.GetPositionMode(accountType, symbol)
			if err != nil {
				return err
			}
			marginCoin := bitgetMarginCoinFromSymbol(symbol, "")
			if marginCoin == "" {
				return ErrorSymbolNotFound
			}
			api := a.privateClient.NewPrivateRestClassicFuturesAccountSetLeverage().
				ProductType(accountType).
				Symbol(symbol).
				MarginCoin(marginCoin)

			switch marginMode {
			case MARGIN_MODE_CROSSED:
				// 全仓模式不需要 holdSide
				api.Leverage(leverage.String())
			case MARGIN_MODE_ISOLATED:
				// 逐仓模式：
				// - 单向持仓：不需要 holdSide，直接 leverage
				// - 双向持仓单边设置：需要 holdSide + long/shortLeverage
				// - 双向持仓多空同杠杆：同时设置 long/shortLeverage，holdSide 非必填
				if posMode == POSITION_MODE_HEDGE {
					switch positionSide {
					case POSITION_SIDE_LONG:
						api.LongLeverage(leverage.String()).
							HoldSide(a.converter.ToBitgetPositionSide(positionSide))
					case POSITION_SIDE_SHORT:
						api.ShortLeverage(leverage.String()).
							HoldSide(a.converter.ToBitgetPositionSide(positionSide))
					default:
						api.LongLeverage(leverage.String()).
							ShortLeverage(leverage.String())
					}
				} else {
					api.Leverage(leverage.String())
				}
			default:
				return ErrorNotSupport
			}

			_, err = api.Do()
			if err != nil {
				return err
			}
			return nil
		default:
			return ErrorNotSupport
		}
	}

	return ErrorNotSupport
}

func (a *BitgetTradeAccount) AssetTransfer(params *AssetTransferParams) ([]*AssetTransfer, error) {
	if err := a.checkMode(); err != nil {
		return nil, err
	}
	return nil, ErrorNotSupport
}

func (a *BitgetTradeAccount) QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	if err := a.checkMode(); err != nil {
		return nil, err
	}
	return nil, ErrorNotSupport
}
