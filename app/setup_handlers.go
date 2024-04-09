package app

import (
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

			if version.Version == "v0.29.31" {
				validators := []string{
					"elysvalcons1j7047ewlfa75dv0q93lnqkctr9afgfayyvmhc4", // euphoria
					"elysvalcons1a58n8t00elj7g4v8lm7rd9q06xu4nz3dgy723q", // shangrila
				}
				for _, val := range validators {
					addr, err := sdk.ConsAddressFromBech32(val)
					if err != nil {
						app.Logger().Error("failed to convert validator address", "error", err)
						continue
					}
					signingInfo, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addr)
					if !found {
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

func loadUpgradeStore(app *ElysApp) {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if shouldLoadUpgradeStore(app, upgradeInfo) {
		storeUpgrades := storetypes.StoreUpgrades{
			// Added: []string{},
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
