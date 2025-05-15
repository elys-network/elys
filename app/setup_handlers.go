package app

import (
	"context"
	"fmt"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"

	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	LocalNetVersion = "v999999"
	NewMaxBytes     = 5 * 1024 * 1024 // 5MB
)

// make sure to update these when you upgrade the version
var NextVersion = "vNEXT"

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
	// required for testnet
	patchParts := strings.Split(parts[2], "-")
	rcVersion := ""
	if len(patchParts) > 1 {
		rcVersion = strings.Join(patchParts[1:], "-")
	}
	if rcVersion != "" {
		return fmt.Sprintf("v%s-%s", majorVersion, rcVersion)
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

			app.AssetprofileKeeper.FixEntries(ctx)

			vm, vmErr := app.mm.RunMigrations(ctx, app.configurator, vm)

			//oracleParams := app.OracleKeeper.GetParams(ctx)
			//if len(oracleParams.MandatoryList) == 0 {
			//	err := app.ojoOracleMigration(ctx, plan.Height+1)
			//	if err != nil {
			//		return nil, err
			//	}
			//}

			// 250USDC from protocol account to masterchef
			params := app.MasterchefKeeper.GetParams(ctx)
			protocolRevenueAddress, err := sdk.AccAddressFromBech32(params.ProtocolRevenueAddress)
			if err != nil {
				return vm, errorsmod.Wrapf(err, "invalid protocol revenue address")
			}

			// Create 250 USDC coin
			// get usdc denom
			usdcDenom, _ := app.AssetprofileKeeper.GetUsdcDenom(ctx)
			usdcAmount := sdk.NewCoin(usdcDenom, sdkmath.NewInt(250000000)) // 250 USDC with 6 decimals

			// Send coins from protocol revenue address to masterchef module
			err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, protocolRevenueAddress, types.ModuleName, sdk.NewCoins(usdcAmount))
			if err != nil {
				// log error
				app.Logger().Error("failed to send USDC to masterchef", "error", err)
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
			//Added:   []string{},
			//Renamed: []storetypes.StoreRename{},
			Deleted: []string{"itransferhook"},
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
