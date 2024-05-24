package app

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
)

func SetupHandlers(app *ElysApp) {
	setUpgradeHandler(app)

	loadUpgradeStore(app)
}

func setUpgradeHandler(app *ElysApp) {
	// Set param key table for params module migration
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		subspace := subspace

		app.Logger().Info("Setting up upgrade handler for " + subspace.Name())

		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
		default:
			continue
		}

		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}

	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())

	app.UpgradeKeeper.SetUpgradeHandler(
		version.Version,
		func(ctx sdk.Context, plan upgradetypes.Plan, vm m.VersionMap) (m.VersionMap, error) {
			app.Logger().Info("Running upgrade handler for " + version.Version)

			// Migrate Tendermint consensus parameters from x/params module to a
			// dedicated x/consensus module.
			baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

			if version.Version == "v0.32.0" {
				// Since invariant is broken sending missing amount to bonded pool
				sumOfValTokens := math.ZeroInt()
				app.EstakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
					sumOfValTokens = sumOfValTokens.Add(validator.GetTokens())
					return false
				})

				totalBondedTokens := app.EstakingKeeper.TotalBondedTokens(ctx)
				if totalBondedTokens.LT(sumOfValTokens) {
					missingAmount := sumOfValTokens.Sub(totalBondedTokens)
					missingCoins := sdk.Coins{sdk.NewCoin("uelys", missingAmount)}
					err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, missingCoins)
					if err != nil {
						panic(err)
					}
					err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, stakingtypes.BondedPoolName, missingCoins)
					if err != nil {
						panic(err)
					}
				}
			}

			return app.mm.RunMigrations(ctx, app.configurator, vm)
		},
	)
}

func loadUpgradeStore(app *ElysApp) {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if shouldLoadUpgradeStore(app, upgradeInfo) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added:   []string{distrtypes.StoreKey, mastercheftypes.StoreKey, estakingtypes.StoreKey},
			Deleted: []string{minttypes.StoreKey},
		}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func shouldLoadUpgradeStore(app *ElysApp, upgradeInfo upgradetypes.Plan) bool {
	return upgradeInfo.Name == version.Version && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height)
}
