package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	LocalNetVersion = "v999.999.999"
	NewMaxBytes     = 5 * 1024 * 1024 // 5MB
)

// make sure to update these when you upgrade the version
var NextVersion = "v0.51.0"

func (app *ElysApp) setUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		version.Version,
		func(goCtx context.Context, plan upgradetypes.Plan, vm m.VersionMap) (m.VersionMap, error) {
			ctx := sdk.UnwrapSDKContext(goCtx)
			app.Logger().Info("Running upgrade handler for " + version.Version)

			if version.Version == NextVersion || version.Version == LocalNetVersion {

				// Add any logic here to run when the chain is upgraded to the new version
				consensusParams, err := app.ConsensusParamsKeeper.ParamsStore.Get(ctx)
				if err != nil {
					return nil, err
				}
				consensusParams.Block.MaxBytes = NewMaxBytes
				app.ConsensusParamsKeeper.ParamsStore.Set(ctx, consensusParams)

				// untombstone validators
				validators := []string{
					"elysvalcons1f9lzcfxxu6l9yj9uf0lqjc0qa82raypnlk58ej", // Synergy Nodes
					"elysvalcons1frn2njtny6gzdjl2df9rvz3atcds2vl2fhxg8s", // Regenerator | Green Validator
				}
				for _, val := range validators {
					addr, err := sdk.ConsAddressFromBech32(val)
					if err != nil {
						app.Logger().Error("failed to convert validator address", "error", err)
						continue
					}
					signingInfo, err := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addr)
					if err != nil {
						app.Logger().Error("failed to get validator signing info", "validator", val)
						continue
					}
					signingInfo.Tombstoned = false
					app.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, signingInfo)
					app.Logger().Info("reset tombstoned status for validator", "validator", val)
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
			// Deleted: []string{},
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
