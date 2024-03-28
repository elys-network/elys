package app

import (
	"fmt"

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

			if version.Version == "v0.29.29" {
				app.Logger().Info("Decommission deprecated pool #1")
				pool, found := app.AmmKeeper.GetPool(ctx, 1)
				if found {
					// withdraw all liquidity from the pool
					allCommitments := app.CommitmentKeeper.GetAllCommitments(ctx)
					for _, commitments := range allCommitments {
						addr, err := sdk.AccAddressFromBech32(commitments.Creator)
						if err != nil {
							continue
						}
						for _, committedToken := range commitments.CommittedTokens {
							if committedToken.Denom == fmt.Sprintf("amm/pool/%d", pool.PoolId) {
								app.AmmKeeper.ExitPool(ctx, addr, pool.PoolId, committedToken.Amount, sdk.NewCoins(), "")
							}
						}
					}
					// remove the pool
					app.AmmKeeper.RemovePool(ctx, pool.PoolId)
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
