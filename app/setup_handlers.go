package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	LocalNetVersion = "v999.999.999"
	NewMaxBytes     = 5 * 1024 * 1024 // 5MB
)

// make sure to update these when you upgrade the version
var NextVersion = "v0.50.0"

func (app *ElysApp) setUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		version.Version,
		func(goCtx context.Context, plan upgradetypes.Plan, vm m.VersionMap) (m.VersionMap, error) {
			ctx := sdk.UnwrapSDKContext(goCtx)
			app.Logger().Info("Running upgrade handler for " + version.Version)

			if version.Version == NextVersion || version.Version == LocalNetVersion {

				// Add any logic here to run when the chain is upgraded to the new version
				// Update consensus params in order to safely enable comet pruning
				consensusParams, err := app.ConsensusParamsKeeper.ParamsStore.Get(ctx)
				if err != nil {
					return nil, err
				}
				consensusParams.Block.MaxBytes = NewMaxBytes
				app.ConsensusParamsKeeper.ParamsStore.Set(ctx, consensusParams)

				// Iterate over all the keys in the wasm module store
				// and delete them
				// TODO: Delete wasm code after deleting wasm module store
				// Retrieve the wasm module store key
				storeKey := app.GetMemKey(wasmtypes.StoreKey)
				store := ctx.KVStore(storeKey)

				iterator := store.Iterator(nil, nil)
				defer iterator.Close()
				for ; iterator.Valid(); iterator.Next() {
					store.Delete(iterator.Key())
				}
			}

			return app.mm.RunMigrations(ctx, app.configurator, vm)
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
			// Added: []string{},
			Deleted: []string{"clock"},
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
	return upgradeInfo.Name == version.Version && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height)
}
