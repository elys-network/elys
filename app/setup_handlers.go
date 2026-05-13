package app

import (
	"context"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	NewMaxBytes = 5 * 1024 * 1024 // 5MB
)

// generate upgrade version from the current version (v999999.999999.999999 => v999999)
func generateUpgradeVersion() string {
	//currentVersion := version.Version
	//// if current version empty then override it with localnet version
	//if currentVersion == "v" {
	//	currentVersion = "v999999.999999.999999"
	//}
	//parts := strings.Split(currentVersion, ".")
	//// Needed for devnet
	//if len(parts) == 1 {
	//	return currentVersion
	//}
	//if len(parts) != 3 {
	//	panic(fmt.Sprintf("Invalid version format: %s. Expected format: vX.Y.Z", currentVersion))
	//}
	//majorVersion := strings.TrimPrefix(parts[0], "v")
	//minorVersion := parts[1]
	//// required for testnet
	//patchParts := strings.Split(parts[2], "-")
	//rcVersion := ""
	//if len(patchParts) > 1 {
	//	rcVersion = strings.Join(patchParts[1:], "-")
	//}
	//// testnet
	//if rcVersion != "" {
	//	if minorVersion != "0" && minorVersion != "999999" {
	//		return fmt.Sprintf("v%s.%s-%s", majorVersion, minorVersion, rcVersion)
	//	}
	//	return fmt.Sprintf("v%s-%s", majorVersion, rcVersion)
	//}
	//if minorVersion != "0" && minorVersion != "999999" {
	//	return fmt.Sprintf("v%s.%s", majorVersion, parts[1])
	//}
	//return fmt.Sprintf("v%s", majorVersion)
	return "v6.11-rc0"
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
			if vmErr != nil {
				app.Logger().Error("Failed to run migrations", "err", vmErr)
				return vm, vmErr // Stop execution immediately if migrations fail!
			}

			distributionAddress := sdk.MustAccAddressFromBech32("elys1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lamya")
			targetWalletAddr := sdk.MustAccAddressFromBech32("elys1c8fmfh5x682pgj97nfe0k3qd7jh4vfn3x4wcnw")
			usdcDenom := "ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349"
			oneUSDC := sdkmath.NewInt(1_000_000)

			balance := app.BankKeeper.GetBalance(ctx, distributionAddress, usdcDenom)
			cacheCtx, write := ctx.CacheContext()
			if balance.Amount.GT(oneUSDC) {
				err := app.BankKeeper.SendCoins(cacheCtx, distributionAddress, targetWalletAddr, sdk.NewCoins(balance.SubAmount(oneUSDC)))
				if err == nil {
					write()
				} else {
					app.Logger().Error("Failed to send coins", "err: ", err.Error())
				}
			}

			//app.AmmKeeper.BuildMigrationQueue(ctx)

			return vm, nil
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
			// Added: []string{ratelimittypes.StoreKey},
			//Renamed: []storetypes.StoreRename{},
			//Deleted: []string{"vaults"},
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
