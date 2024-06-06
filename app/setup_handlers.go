package app

import (
	"cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	membershiptiertypes "github.com/elys-network/elys/x/membershiptier/types"
)

func SetupHandlers(app *ElysApp) {
	setUpgradeHandler(app)

	loadUpgradeStore(app)
}

func setUpgradeHandler(app *ElysApp) {
	app.UpgradeKeeper.SetUpgradeHandler(
		version.Version,
		func(ctx sdk.Context, plan upgradetypes.Plan, vm m.VersionMap) (m.VersionMap, error) {
			app.Logger().Info("Running upgrade handler for " + version.Version)

			if version.Version == "v0.32.0" || version.Version == "v999.999.999" {
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

				app.Logger().Info("Deleting proposals with ID <= 193")
				store := ctx.KVStore(app.keys[govtypes.StoreKey])
				for i := uint64(1); i <= 193; i++ {
					store.Delete(govtypes.ProposalKey(i))
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
			Added: []string{membershiptiertypes.StoreKey},
			// Deleted: []string{},
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
