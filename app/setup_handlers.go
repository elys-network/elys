package app

import (
	"context"
	"fmt"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	NewMaxBytes = 5 * 1024 * 1024 // 5MB
)

// generate upgrade version from the current version (v999999.999999.999999 => v999999)
func generateUpgradeVersion() string {
	currentVersion := version.Version
	// if current version empty then override it with localnet version
	if currentVersion == "v" {
		currentVersion = "v999999.999999.999999"
	}
	parts := strings.Split(currentVersion, ".")
	// Needed for devnet
	if len(parts) == 1 {
		return currentVersion
	}
	if len(parts) != 3 {
		panic(fmt.Sprintf("Invalid version format: %s. Expected format: vX.Y.Z", currentVersion))
	}
	majorVersion := strings.TrimPrefix(parts[0], "v")
	minorVersion := parts[1]
	// required for testnet
	patchParts := strings.Split(parts[2], "-")
	rcVersion := ""
	if len(patchParts) > 1 {
		rcVersion = strings.Join(patchParts[1:], "-")
	}
	// testnet
	if rcVersion != "" {
		if minorVersion != "0" && minorVersion != "999999" {
			return fmt.Sprintf("v%s.%s-%s", majorVersion, minorVersion, rcVersion)
		}
		return fmt.Sprintf("v%s-%s", majorVersion, rcVersion)
	}
	if minorVersion != "0" && minorVersion != "999999" {
		return fmt.Sprintf("v%s.%s", majorVersion, parts[1])
	}
	return fmt.Sprintf("v%s", majorVersion)
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

			for _, profile := range app.AssetprofileKeeper.GetAllEntry(ctx) {
				if profile.DisplayName == "WBTC" || profile.DisplayName == "wBTC" {
					profile.DisplayName = "BTC"
				}
				if profile.DisplayName == "WETH" || profile.DisplayName == "wETH" {
					profile.DisplayName = "ETH"
				}
				app.AssetprofileKeeper.SetEntry(ctx, profile)
			}

			for _, assetInfo := range app.LegacyOracleKeepper.GetAllAssetInfo(ctx) {
				if assetInfo.Display == "WBTC" || assetInfo.Display == "wBTC" {
					assetInfo.Display = "BTC"
					assetInfo.BandTicker = "BTC"
					assetInfo.ElysTicker = "BTC"
				}
				if assetInfo.Display == "WETH" || assetInfo.Display == "wETH" {
					assetInfo.Display = "ETH"
					assetInfo.BandTicker = "ETH"
					assetInfo.ElysTicker = "ETH"
				}
				app.LegacyOracleKeepper.SetAssetInfo(ctx, assetInfo)
			}

			for _, price := range app.LegacyOracleKeepper.GetAllAssetPrice(ctx, "WBTC") {
				price.Asset = "BTC"
				app.LegacyOracleKeepper.SetPrice(ctx, price)
			}

			for _, price := range app.LegacyOracleKeepper.GetAllAssetPrice(ctx, "WETH") {
				price.Asset = "ETH"
				app.LegacyOracleKeepper.SetPrice(ctx, price)
			}

			oracleParams := app.OracleKeeper.GetParams(ctx)
			if len(oracleParams.MandatoryList) == 0 {
				err := app.ojoOracleMigration(ctx, plan.Height+1)
				if err != nil {
					return nil, err
				}
			}

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
			// Added: []string{ratelimittypes.StoreKey},
			//Renamed: []storetypes.StoreRename{},
			//Deleted: []string{ratelimittypes.StoreKey},
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
