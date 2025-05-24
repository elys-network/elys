package app

//import (
//	"cosmossdk.io/math"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
//	leagacyOracletypes "github.com/elys-network/elys/v5/x/oracle/types"
//	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
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
//		},
//	},
//	{
//		BaseDenom:  "AKT",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"gate",
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
//		},
//	},
//	{
//		BaseDenom:  "KAVA",
//		QuoteDenom: "USDT",
//		Providers: []string{
//			"binance",
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
//	err := app.enableVoteExtensions(ctx, height)
//	if err != nil {
//		return err
//	}
//
//	// Add any logic here to run when the chain is upgraded to the new version
//	app.Logger().Info("Migrating legacy oracle prices to new oracle")
//	legacyParams := app.LegacyOracleKeepper.GetParams(ctx)
//	prices := app.LegacyOracleKeepper.GetAllPrice(ctx)
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
//	newParams.LifeTimeInBlocks = 2
//	newParams.PriceExpiryTime = legacyParams.PriceExpiryTime
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
//			Source:      price.Source,
//			Provider:    price.Provider,
//			Timestamp:   price.Timestamp,
//			BlockHeight: price.BlockHeight,
//		})
//		app.LegacyOracleKeepper.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
//	}
//	//Migrate AssetInfos
//	for _, assetInfo := range assetInfos {
//		app.OracleKeeper.SetAssetInfo(ctx, ojooracletypes.AssetInfo{
//			Denom:      assetInfo.Denom,
//			Display:    assetInfo.Display,
//			BandTicker: assetInfo.BandTicker,
//			ElysTicker: assetInfo.ElysTicker,
//			Decimal:    assetInfo.Decimal,
//		})
//		app.LegacyOracleKeepper.RemoveAssetInfo(ctx, assetInfo.Denom)
//	}
//
//	return nil
//}
