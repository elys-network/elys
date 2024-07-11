package app

import (
	"fmt"

	wasmmodule "github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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

			if version.Version == "v0.38.0" || version.Version == "v999.999.999" {
				// Retrieve the wasm module store key
				storeKey := app.keys[wasmmodule.StoreKey]

				// Retrieve the wasm module store
				store := ctx.KVStore(storeKey)

				// List of prefixes to clear
				prefixes := [][]byte{
					// Account History contract
					wasmtypes.GetContractStorePrefix(sdk.MustAccAddressFromBech32("elys17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgs98tvuy")),
				}

				// Add old code keys to the list of prefixes to clear
				for i := uint64(1); i < 671; i++ {
					codeKey := wasmtypes.GetCodeKey(i)
					// append the code key to the prefixes
					prefixes = append(prefixes, codeKey)
				}

				// Clear all keys in the store
				for _, prefix := range prefixes {
					iter := sdk.KVStorePrefixIterator(store, prefix)
					defer iter.Close()

					for ; iter.Valid(); iter.Next() {
						store.Delete(iter.Key())
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
		panic(fmt.Sprintf("Failed to read upgrade info from disk: %v", err))
	}

	fmt.Printf("Upgrade info: %+v\n", upgradeInfo)

	if shouldLoadUpgradeStore(app, upgradeInfo) {
		storeUpgrades := storetypes.StoreUpgrades{
			// Added: []string{},
			// Deleted: []string{},
		}
		fmt.Printf("Setting store loader with height %d and store upgrades: %+v\n", upgradeInfo.Height, storeUpgrades)

		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	} else {
		fmt.Println("No need to load upgrade store.")
	}
}

func shouldLoadUpgradeStore(app *ElysApp, upgradeInfo upgradetypes.Plan) bool {
	currentHeight := app.LastBlockHeight()
	fmt.Printf("Current block height: %d, Upgrade height: %d\n", currentHeight, upgradeInfo.Height)
	return upgradeInfo.Name == version.Version && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height)
}
