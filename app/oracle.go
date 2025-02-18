package app

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	leagacyOracletypes "github.com/elys-network/elys/x/oracle/types"
	ojooracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

var currencyPairProviders = []ojooracletypes.CurrencyPairProviders{
	{
		BaseDenom:  "USDT",
		QuoteDenom: "USD",
		Providers: []string{
			"kraken",
			"coinbase",
			"crypto",
		},
	},
	{
		BaseDenom:  "USDC",
		QuoteDenom: "USD",
		Providers: []string{
			"kraken",
		},
	},
	{
		BaseDenom:  "ATOM",
		QuoteDenom: "USD",
		Providers: []string{
			"coinbase",
			"kraken",
		},
	},
	{
		BaseDenom:  "ATOM",
		QuoteDenom: "USDT",
		Providers: []string{
			"binance",
			//"okx",
			//"bitget",
			"gate",
		},
		ExternLiquidityProvider: "binance",
		PoolId:                  1,
	},
	{
		BaseDenom:  "ATOM",
		QuoteDenom: "USDC",
		Providers: []string{
			"mexc",
		},
	},
	{
		BaseDenom:  "AKT",
		QuoteDenom: "USD",
		Providers: []string{
			"coinbase",
		},
	},
	{
		BaseDenom:  "AKT",
		QuoteDenom: "USDT",
		Providers: []string{
			"gate",
		},
		ExternLiquidityProvider: "gate",
		PoolId:                  3,
	},
	{
		BaseDenom:  "TIA",
		QuoteDenom: "USDT",
		Providers: []string{
			"binance",
			"mexc",
			"coinbase",
			"gate",
		},
		BaseProxyDenom:          "TIA",
		QuoteProxyDenom:         "USDC",
		ExternLiquidityProvider: "binance",
		PoolId:                  2,
	},
	{
		BaseDenom:  "KAVA",
		QuoteDenom: "USDT",
		Providers: []string{
			"binance",
			"mexc",
			//"bitget",
		},
		ExternLiquidityProvider: "binance",
	},
	{
		BaseDenom:  "SAGA",
		QuoteDenom: "USDT",
		Providers: []string{
			"binance",
			"mexc",
			//"bitget",
		},
		ExternLiquidityProvider: "binance",
	},
	{
		BaseDenom:  "XION",
		QuoteDenom: "USDT",
		Providers: []string{
			"mexc",
			//"bitget",
			"gate",
		},
		ExternLiquidityProvider: "gate",
	},
	{
		BaseDenom:  "SCRT",
		QuoteDenom: "USDT",
		Providers: []string{
			"binance",
			"mexc",
			"gate",
		},
		ExternLiquidityProvider: "binance",
	},
	{
		BaseDenom:  "SCRT",
		QuoteDenom: "USD",
		Providers: []string{
			"kraken",
		},
	},
	{
		BaseDenom:  "OSMO",
		QuoteDenom: "USDT",
		Providers: []string{
			"binance",
			"mexc",
			"gate",
			"huobi",
		},
		ExternLiquidityProvider: "binance",
	},
	{
		BaseDenom:  "NTRN",
		QuoteDenom: "USDT",
		Providers: []string{
			"mexc",
			"binance",
			//"bitget",
		},
		ExternLiquidityProvider: "binance",
	},
}

func updateAndGetCurrencyProviders(ammPools []ammtypes.Pool, assetInfos []leagacyOracletypes.AssetInfo) []ojooracletypes.CurrencyPairProviders {
	poolMap := make(map[string]uint64)
	for _, assetInfo := range assetInfos {
		for _, ammPool := range ammPools {
			for _, poolAsset := range ammPool.PoolAssets {
				if poolAsset.Token.Denom == assetInfo.Denom && assetInfo.Display != ptypes.USDC_DISPLAY {
					poolMap[assetInfo.Display] = ammPool.PoolId
				}
			}
		}
	}
	for i, _ := range currencyPairProviders {
		poolId := poolMap[currencyPairProviders[i].BaseDenom]
		if currencyPairProviders[i].ExternLiquidityProvider != "" {
			currencyPairProviders[i].PoolId = poolId
		}
	}
	return currencyPairProviders
}

func (app *ElysApp) ojoOracleMigration(ctx sdk.Context, plan upgradetypes.Plan) error {
	consensusParams, err := app.ConsensusParamsKeeper.ParamsStore.Get(ctx)
	if err != nil {
		return err
	}
	if consensusParams.Abci.VoteExtensionsEnableHeight == 0 {
		consensusParams.Abci.VoteExtensionsEnableHeight = plan.Height + 1
		_, err = app.ConsensusParamsKeeper.UpdateParams(ctx, &consensustypes.MsgUpdateParams{
			Authority: app.ConsensusParamsKeeper.GetAuthority(),
			Block:     consensusParams.Block,
			Evidence:  consensusParams.Evidence,
			Validator: consensusParams.Validator,
			Abci:      consensusParams.Abci,
		})
		if err != nil {
			return err
		}
	}

	// Add any logic here to run when the chain is upgraded to the new version
	app.Logger().Info("Migrating legacy oracle prices to new oracle")
	legacyParams := app.LegacyOracleKeepper.GetParams(ctx)
	prices := app.LegacyOracleKeepper.GetAllPrice(ctx)
	assetInfos := app.LegacyOracleKeepper.GetAllAssetInfo(ctx)
	priceFeeders := app.LegacyOracleKeepper.GetAllPriceFeeder(ctx)
	allAmmPool := app.AmmKeeper.GetAllPool(ctx)

	// Set Params
	var denomList []ojooracletypes.Denom
	// Need to add USDT as CEX provide prices in USDT
	denomList = append(denomList, ojooracletypes.Denom{
		BaseDenom:   "uusdt",
		SymbolDenom: ptypes.USDT_DISPLAY,
		Exponent:    6,
	})
	currencyPairs := updateAndGetCurrencyProviders(allAmmPool, assetInfos)
	var deviationThreshold []ojooracletypes.CurrencyDeviationThreshold
	for _, assetInfo := range assetInfos {
		if assetInfo.Denom != ptypes.Elys && assetInfo.Denom != ptypes.Eden && assetInfo.Denom != ptypes.EdenB {
			denomList = append(denomList, ojooracletypes.Denom{
				BaseDenom:   assetInfo.Denom,
				SymbolDenom: assetInfo.Display,
				Exponent:    uint32(assetInfo.Decimal),
			})
			deviationThreshold = append(deviationThreshold, ojooracletypes.CurrencyDeviationThreshold{
				BaseDenom: assetInfo.Display,
				Threshold: "2",
			})
		}
	}
	newParams := ojooracletypes.DefaultParams()
	newParams.LifeTimeInBlocks = legacyParams.LifeTimeInBlocks
	newParams.PriceExpiryTime = legacyParams.PriceExpiryTime
	newParams.MandatoryList = denomList
	newParams.AcceptList = denomList
	newParams.CurrencyPairProviders = currencyPairs
	newParams.CurrencyDeviationThresholds = deviationThreshold
	if err = newParams.Validate(); err != nil {
		return err
	}
	app.OracleKeeper.SetParams(ctx, newParams)

	// Delete Port
	app.LegacyOracleKeepper.DeletePort(ctx)
	// Migrate Feeders
	for _, feeder := range priceFeeders {
		app.OracleKeeper.SetPriceFeeder(ctx, ojooracletypes.PriceFeeder{
			Feeder:   feeder.Feeder,
			IsActive: feeder.IsActive,
		})
		app.LegacyOracleKeepper.RemovePriceFeeder(ctx, feeder.GetFeederAccount())
	}

	//Migrate Prices
	for _, price := range prices {
		app.OracleKeeper.SetPrice(ctx, ojooracletypes.Price{
			Asset:       price.Asset,
			Price:       price.Price,
			Source:      price.Source,
			Provider:    price.Provider,
			Timestamp:   price.Timestamp,
			BlockHeight: price.BlockHeight,
		})
		app.LegacyOracleKeepper.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
	}
	//Migrate AssetInfos
	for _, assetInfo := range assetInfos {
		app.OracleKeeper.SetAssetInfo(ctx, ojooracletypes.AssetInfo{
			Denom:      assetInfo.Denom,
			Display:    assetInfo.Display,
			BandTicker: assetInfo.BandTicker,
			ElysTicker: assetInfo.ElysTicker,
			Decimal:    assetInfo.Decimal,
		})
		app.LegacyOracleKeepper.RemoveAssetInfo(ctx, assetInfo.Denom)
	}

	return nil
}
