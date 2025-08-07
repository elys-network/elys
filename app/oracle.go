package app

//
//import (
//	"cosmossdk.io/math"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
//	leagacyOracletypes "github.com/elys-network/elys/v7/x/oracle/types"
//	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
//)
//
//var currencyPairProviders = []ojooracletypes.CurrencyPairProviders{
//	{
//		BaseDenom:  "USDT",
//		QuoteDenom: "USD",
//		Providers: []string{
//			"kraken",
//			"coinbase",
//		},
//	},
//	{
//		BaseDenom:  "USDC",
//		QuoteDenom: "USD",
//		Providers: []string{
//			"kraken",
//		},
//	},
//	{
//		BaseDenom:  "ATOM",
//		QuoteDenom: "USD",
//		Providers: []string{
//			"coinbase",
//			"kraken",
//		},
//	},
//	{
//		BaseDenom:  "ATOM",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"okx",
//			"bitget",
//			"mexc",
//			"huobi",
//			"gate",
//		},
//		BaseProxyDenom:          "ATOM",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//		PoolId:                  1,
//	},
//	{
//		BaseDenom:  "AKT",
//		QuoteDenom: "USD",
//		Providers: []string{
//			"coinbase",
//			"kraken",
//		},
//	},
//	{
//		BaseDenom:  "AKT",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"gate",
//			"huobi",
//		},
//		BaseProxyDenom:          "AKT",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "gate",
//		PoolId:                  3,
//	},
//	{
//		BaseDenom:  "TIA",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"okx",
//			"huobi",
//			"bitget",
//			"gate",
//		},
//		BaseProxyDenom:          "TIA",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//		PoolId:                  2,
//	},
//	{
//		BaseDenom:  "TIA",
//		QuoteDenom: "USD",
//		Providers: []string{
//			"coinbase",
//			"kraken",
//		},
//	},
//	{
//		BaseDenom:  "KAVA",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"mexc",
//			"gate",
//			"bitget",
//		},
//		BaseProxyDenom:          "KAVA",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:  "SAGA",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"huobi",
//			"bitget",
//			"gate",
//		},
//		BaseProxyDenom:          "SAGA",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:  "XION",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"bitget",
//			"gate",
//			"huobi",
//			"mexc",
//		},
//		BaseProxyDenom:          "XION",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "gate",
//	},
//	{
//		BaseDenom:  "SCRT",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"huobi",
//			"gate",
//			"mexc",
//		},
//		BaseProxyDenom:          "SCRT",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:  "SCRT",
//		QuoteDenom: "USD",
//		Providers: []string{
//			"kraken",
//		},
//	},
//	{
//		BaseDenom:  "OSMO",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"bitget",
//			"gate",
//			"huobi",
//		},
//		BaseProxyDenom:          "OSMO",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:  "NTRN",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"bitget",
//			"gate",
//			"mexc",
//		},
//		BaseProxyDenom:          "NTRN",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:  "OM",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
//			"gate",
//			"okx",
//			"bitget",
//		},
//		BaseProxyDenom:          "OM",
//		QuoteProxyDenom:         "USDC",
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "BLD",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "BLD",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"gate",
//			"huobi",
//		},
//		ExternLiquidityProvider: "gate",
//	},
//	{
//		BaseDenom:       "BTC",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "BTC",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"mexc",
//			"gate",
//			"huobi",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "ETH",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "ETH",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"mexc",
//			"gate",
//			"huobi",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "PAXG",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "PAXG",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"gate",
//			"mexc",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "BABY",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "BABY",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"gate",
//			"okx",
//			"mexc",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "FET",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "FET",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"gate",
//			"huobi",
//			"okx",
//			"mexc",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:  "STRD",
//		QuoteDenom: "USD",
//		Providers: []string{
//			"kraken",
//		},
//	},
//	{
//		BaseDenom:       "STRD",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "STRD",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"mexc",
//		},
//		ExternLiquidityProvider: "mexc",
//	},
//	{
//		BaseDenom:       "INJ",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "INJ",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"okx",
//			"mexc",
//			"gate",
//			"huobi",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "XRP",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "XRP",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"gate",
//			"okx",
//			"mexc",
//			"huobi",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "LINK",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "LINK",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"gate",
//			"mexc",
//			"okx",
//			"huobi",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//	{
//		BaseDenom:       "ONDO",
//		QuoteDenom:      "USDT",
//		BaseProxyDenom:  "ONDO",
//		QuoteProxyDenom: "USDC",
//		Providers: []string{
//			"binance",
//			"gate",
//			"mexc",
//			"okx",
//			"huobi",
//		},
//		ExternLiquidityProvider: "binance",
//	},
//}
//
//func addDenomToList(denom, display string, decimal uint64, denomList []ojooracletypes.Denom, rewardBand []ojooracletypes.RewardBand, deviationThreshold []ojooracletypes.CurrencyDeviationThreshold) ([]ojooracletypes.Denom, []ojooracletypes.RewardBand, []ojooracletypes.CurrencyDeviationThreshold) {
//	denomList = append(denomList, ojooracletypes.Denom{
//		BaseDenom:   denom,
//		SymbolDenom: display,
//		Exponent:    uint32(decimal),
//	})
//	deviationThreshold = append(deviationThreshold, ojooracletypes.CurrencyDeviationThreshold{
//		BaseDenom: display,
//		Threshold: "2",
//	})
//	rewardBand = append(rewardBand, ojooracletypes.RewardBand{
//		SymbolDenom: display,
//		RewardBand:  math.LegacyNewDecWithPrec(2, 2),
//	})
//	return denomList, rewardBand, deviationThreshold
//}
//
//func (app *ElysApp) enableVoteExtensions(ctx sdk.Context, height int64) error {
//	consensusParams, err := app.ConsensusParamsKeeper.ParamsStore.Get(ctx)
//	if err != nil {
//		return err
//	}
//	if consensusParams.Abci.VoteExtensionsEnableHeight == 0 {
//		consensusParams.Abci.VoteExtensionsEnableHeight = height
//		_, err = app.ConsensusParamsKeeper.UpdateParams(ctx, &consensustypes.MsgUpdateParams{
//			Authority: app.ConsensusParamsKeeper.GetAuthority(),
//			Block:     consensusParams.Block,
//			Evidence:  consensusParams.Evidence,
//			Validator: consensusParams.Validator,
//			Abci:      consensusParams.Abci,
//		})
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (app *ElysApp) ojoOracleMigration(ctx sdk.Context, height int64) error {
//
//	for _, profile := range app.AssetprofileKeeper.GetAllEntry(ctx) {
//		if profile.DisplayName == "WBTC" || profile.DisplayName == "wBTC" {
//			profile.DisplayName = "BTC"
//		}
//		if profile.DisplayName == "WETH" || profile.DisplayName == "wETH" {
//			profile.DisplayName = "ETH"
//		}
//		app.AssetprofileKeeper.SetEntry(ctx, profile)
//	}
//
//	for _, assetInfo := range app.LegacyOracleKeepper.GetAllAssetInfo(ctx) {
//		if assetInfo.Display == "WBTC" || assetInfo.Display == "wBTC" {
//			assetInfo.Display = "BTC"
//			assetInfo.BandTicker = "BTC"
//			assetInfo.ElysTicker = "BTC"
//		}
//		if assetInfo.Display == "WETH" || assetInfo.Display == "wETH" {
//			assetInfo.Display = "ETH"
//			assetInfo.BandTicker = "ETH"
//			assetInfo.ElysTicker = "ETH"
//		}
//		app.LegacyOracleKeepper.SetAssetInfo(ctx, assetInfo)
//	}
//
//	for _, price := range app.LegacyOracleKeepper.GetAllAssetPrice(ctx, "WBTC") {
//		price.Asset = "BTC"
//		app.LegacyOracleKeepper.SetPrice(ctx, price)
//	}
//
//	for _, price := range app.LegacyOracleKeepper.GetAllAssetPrice(ctx, "WETH") {
//		price.Asset = "ETH"
//		app.LegacyOracleKeepper.SetPrice(ctx, price)
//	}
//
//	err := app.enableVoteExtensions(ctx, height)
//	if err != nil {
//		return err
//	}
//
//	// Add any logic here to run when the chain is upgraded to the new version
//	app.Logger().Info("Migrating legacy oracle prices to new oracle")
//	prices := app.LegacyOracleKeepper.GetAllLegacyPrice(ctx)
//	assetInfos := app.LegacyOracleKeepper.GetAllAssetInfo(ctx)
//	priceFeeders := app.LegacyOracleKeepper.GetAllPriceFeeder(ctx)
//	allAmmPool := app.AmmKeeper.GetAllPool(ctx)
//
//	var denomList []ojooracletypes.Denom
//	var rewardBand []ojooracletypes.RewardBand
//	usdtAssetInfo := leagacyOracletypes.AssetInfo{}
//	usdcAssetInfo := leagacyOracletypes.AssetInfo{}
//	var deviationThreshold []ojooracletypes.CurrencyDeviationThreshold
//
//	assetInfoMap := make(map[string]leagacyOracletypes.AssetInfo, len(assetInfos))
//	for _, assetInfo := range assetInfos {
//		assetInfoMap[assetInfo.Denom] = assetInfo
//
//		if assetInfo.Display == ptypes.USDC_DISPLAY {
//			usdcAssetInfo = assetInfo
//
//		}
//
//		if assetInfo.Display == ptypes.USDT_DISPLAY {
//			usdtAssetInfo = assetInfo
//		}
//	}
//	if usdtAssetInfo.Denom == "" {
//		denomList, rewardBand, deviationThreshold = addDenomToList("uusdt", ptypes.USDT_DISPLAY, 6, denomList, rewardBand, deviationThreshold)
//	} else {
//		denomList, rewardBand, deviationThreshold = addDenomToList(usdtAssetInfo.Denom, usdtAssetInfo.Display, usdtAssetInfo.Decimal, denomList, rewardBand, deviationThreshold)
//	}
//	denomList, rewardBand, deviationThreshold = addDenomToList(usdcAssetInfo.Denom, usdcAssetInfo.Display, usdcAssetInfo.Decimal, denomList, rewardBand, deviationThreshold)
//
//	// Set all pool id for external liquidity factor
//	for _, ammPool := range allAmmPool {
//		// Only oracle pools
//		if ammPool.PoolParams.UseOracle {
//			for _, poolAsset := range ammPool.PoolAssets {
//				assetInfo := assetInfoMap[poolAsset.Token.Denom]
//				if assetInfo.Display != ptypes.USDC_DISPLAY && assetInfo.Display != ptypes.USDT_DISPLAY {
//
//					denomList, rewardBand, deviationThreshold = addDenomToList(assetInfo.Denom, assetInfo.Display, assetInfo.Decimal, denomList, rewardBand, deviationThreshold)
//
//					for i, _ := range currencyPairProviders {
//						if currencyPairProviders[i].BaseDenom == assetInfo.Display && currencyPairProviders[i].ExternLiquidityProvider != "" {
//							currencyPairProviders[i].PoolId = ammPool.PoolId
//							break
//						}
//					}
//				}
//			}
//		}
//	}
//
//	newParams := ojooracletypes.DefaultParams()
//	newParams.MandatoryList = denomList
//	newParams.AcceptList = denomList
//	newParams.CurrencyPairProviders = currencyPairProviders
//	newParams.CurrencyDeviationThresholds = deviationThreshold
//	newParams.RewardBands = rewardBand
//
//	if err = newParams.Validate(); err != nil {
//		return err
//	}
//	app.OracleKeeper.SetParams(ctx, newParams)
//
//	// Delete Port
//	app.LegacyOracleKeepper.DeletePort(ctx)
//	// Migrate Feeders
//	for _, feeder := range priceFeeders {
//		app.OracleKeeper.SetPriceFeeder(ctx, ojooracletypes.PriceFeeder{
//			Feeder:   feeder.Feeder,
//			IsActive: feeder.IsActive,
//		})
//		app.LegacyOracleKeepper.RemovePriceFeeder(ctx, feeder.GetFeederAccount())
//	}
//
//	//Migrate Prices
//	for _, price := range prices {
//		app.OracleKeeper.SetPrice(ctx, ojooracletypes.Price{
//			Asset:       price.Asset,
//			Price:       price.Price,
//			Provider:    price.Provider,
//			Timestamp:   price.Timestamp,
//			BlockHeight: price.BlockHeight,
//		})
//		app.LegacyOracleKeepper.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
//	}
//	//Migrate AssetInfos
//	for _, assetInfo := range assetInfos {
//		app.OracleKeeper.SetAssetInfo(ctx, ojooracletypes.AssetInfo{
//			Denom:   assetInfo.Denom,
//			Display: assetInfo.Display,
//			Ticker:  assetInfo.ElysTicker,
//			Decimal: assetInfo.Decimal,
//		})
//		app.LegacyOracleKeepper.RemoveAssetInfo(ctx, assetInfo.Denom)
//	}
//
//	return nil
//}
