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

			if version.Version == "v0.40.0" || version.Version == "v999.999.999" {
				// Retrieve the wasm module store key
				storeKey := app.keys[wasmmodule.StoreKey]

				// Retrieve the wasm module store
				store := ctx.KVStore(storeKey)

				// List of prefixes to clear
				prefixes := [][]byte{
					wasmtypes.GetContractStorePrefix(sdk.MustAccAddressFromBech32("elys1s37xz7tzrru2cpl96juu9lfqrsd4jh73j9slyv440q5vttx2uyesetjpne")), // AH
					wasmtypes.GetContractStorePrefix(sdk.MustAccAddressFromBech32("elys1g2xwx805epc897rwyrykskjque07yxfmc4qq2p4ef5dwd6znl30qnxje76")), // FS
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1s37xz7tzrru2cpl96juu9lfqrsd4jh73j9slyv440q5vttx2uyesetjpne")),  // AH
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1g2xwx805epc897rwyrykskjque07yxfmc4qq2p4ef5dwd6znl30qnxje76")),  // FS
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1x8gwn06l85q0lyncy7zsde8zzdn588k2dck00a8j6lkprydcutwqa9tv6n")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1657pee2jhf4jk8pq6yq64e758ngvum45gl866knmjkd83w6jgn3s923j5j")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1xhcxq4fvxth2hn3msmkpftkfpw73um7s4et3lh4r8cfmumk3qsmsmgjjrc")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1wr6vc3g4caz9aclgjacxewr0pjlre9wl2uhq73rp8mawwmqaczsq3ppn83")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys15m728qxvtat337jdu2f0uk6pu905kktrxclgy36c0wd822tpxcmqfzew4d")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1ul4msjc3mmaxsscdgdtjds85rg50qrepvrczp0ldgma5mm9xv8yqxhk8nu")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1mx32w9tnfxv0z5j000750h8ver7qf3xpj09w3uzvsr3hq68f4hxqte4gam")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys14zykjnz94dr9nj4v2yzpvnlrw5uurk5h7d5w0wug902vxdynm6xsue684e")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys175r6y463k8cdcte6dzrxydxnwfkhz9afdghzcjxxhzfmm6rgu64qdp9z37")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys1jyhyqjxf3pc7vzwyqhwe53up5pj0e53zw3xu2589uqgkvqngswnqtxfw4e")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys14see2dq4nu37yk9qhjn2laqxrmzzjyxwhfgnxw4nuzpm7vc6ztysxjv4p5")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys15qe27v4z7j78g5g4ak2ftftky3c078zvtr9qtv5lhxwc54ccf4asggmyyp")),  // old contract
					wasmtypes.GetContractAddressKey(sdk.MustAccAddressFromBech32("elys193dzg6ealfymax4pyrkge60swlr2tjupwegdemgalzhkkxc8kzyqh5qw9c")),  // old contract
				}

				// Add old code keys to the list of prefixes to clear
				for i := uint64(675); i < 680; i++ {
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
