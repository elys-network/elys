package app

import (
	"context"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	ratelimittypes "github.com/cosmos/ibc-apps/modules/rate-limiting/v8/types"
)

const (
	NewMaxBytes = 5 * 1024 * 1024 // 5MB
)

// generate upgrade version from the current version (v999999.999999.999999 => v999999)
func generateUpgradeVersion() string {
	return "v6.3"
}

func (app *ElysApp) setUpgradeHandler() {
	upgradeVersion := generateUpgradeVersion()
	app.Logger().Info("Current version", "version", version.Version)
	app.Logger().Info("Upgrade version", "version", upgradeVersion)
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeVersion,
		func(goCtx context.Context, plan upgradetypes.Plan, vm m.VersionMap) (m.VersionMap, error) {
			ctx := sdk.UnwrapSDKContext(goCtx)
			app.Logger().Info("Running upgrade handler for " + upgradeVersion)

			vm, vmErr := app.mm.RunMigrations(ctx, app.configurator, vm)

			//oracleParams := app.OracleKeeper.GetParams(ctx)
			//if len(oracleParams.MandatoryList) == 0 {
			//	err := app.ojoOracleMigration(ctx, plan.Height+1)
			//	if err != nil {
			//		return nil, err
			//	}
			//}

			for _, pool := range app.LeveragelpKeeper.GetAllPools(ctx) {
				pool.MaxLeveragelpRatio = math.LegacyMustNewDecFromStr("0.35")
				app.LeveragelpKeeper.SetPool(ctx, pool)
			}

			perpetualParams := app.PerpetualKeeper.GetParams(ctx)
			perpetualParams.ExitBuffer = math.LegacyMustNewDecFromStr("0.1")
			err := app.PerpetualKeeper.SetParams(ctx, &perpetualParams)
			if err != nil {
				panic(err)
			}

			app.OracleKeeper.DeleteAXLPrices(ctx)

			return vm, vmErr
		},
	)
}

func (app *ElysApp) setUpgradeStore() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("Failed to read upgrade info from disk: %v", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	app.Logger().Debug("Upgrade info", "info", upgradeInfo)

	if shouldLoadUpgradeStore(app, upgradeInfo) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{ratelimittypes.StoreKey},
			//Renamed: []storetypes.StoreRename{},
			//Deleted: []string{ibcfeetypes.StoreKey},
		}
		app.Logger().Info(fmt.Sprintf("Setting store loader with height %d and store upgrades: %+v\n", upgradeInfo.Height, storeUpgrades))

		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else {
		app.Logger().Debug("No need to load upgrade store.")
	}
}

func shouldLoadUpgradeStore(app *ElysApp, upgradeInfo upgradetypes.Plan) bool {
	currentHeight := app.LastBlockHeight()
	app.Logger().Debug(fmt.Sprintf("Current block height: %d, Upgrade height: %d\n", currentHeight, upgradeInfo.Height))
	upgradeVersion := generateUpgradeVersion()
	app.Logger().Debug("Current version", "version", version.Version)
	app.Logger().Debug("Upgrade version", "version", upgradeVersion)
	return upgradeInfo.Name == upgradeVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height)
}
