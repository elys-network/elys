package app

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	oracletypes "github.com/elys-network/elys/v7/x/oracle/types"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vaulttypes "github.com/elys-network/elys/v7/x/vaults/types"

	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
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

			if ctx.ChainID() == "elys-1" {
				newPriceFeeders := []sdk.AccAddress{
					sdk.MustAccAddressFromBech32("elys1y5jeztqtf7vwqe6wd7tv2z8mzjg38hn6zkpmvq"),
					sdk.MustAccAddressFromBech32("elys16r4y6hdehvntgg70avg6f0s2x55k3wekeyw8vw"),
				}

				for _, accAddress := range newPriceFeeders {
					v := oracletypes.PriceFeeder{
						Feeder:   accAddress.String(),
						IsActive: true,
					}
					app.OracleKeeper.SetPriceFeeder(ctx, v)
				}
			}

			receiver := sdk.MustAccAddressFromBech32("elys1kxgan4uq0m8gqztd09n6qm627p6v4ayzngxyx8")
			sender := sdk.MustAccAddressFromBech32("elys1p2fhrn9zfra9lv5062nvzkmp9hduhm9hkk6kz6a3ucwktjvuzv9smry5sk")
			amount := sdk.NewCoins(
				sdk.NewInt64Coin("ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349", 85_000_000),                                                    //USDC
				sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 1_036_566_310),                                                 // ATOM
				sdk.NewInt64Coin("ibc/8BFE59DCD5A7054F0A97CF91F3E3ABCA8C5BA454E548FA512B729D4004584D47", 3_000_000),                                                     // NTRN
				sdk.NewCoin("ibc/8464A63954C0350A26C8588E20719F3A0AC8705E4CA0F7450B60C3F16B2D3421", math.LegacyMustNewDecFromStr("10000000000000000000").TruncateInt()), // XRP
				sdk.NewCoin("ibc/ADF401C952ADD9EE232D52C8303B8BE17FE7953C8D420F20769AF77240BD0C58", math.LegacyMustNewDecFromStr("51059039885106925").TruncateInt()),    // INJ
				sdk.NewInt64Coin("ibc/45D6B52CAD911A15BD9C2F5FFDA80E26AFCB05C7CD520070790ABC86D2B24229", 9_960_000),                                                     // TIA
			)
			if ctx.ChainID() == "elysicstestnet-1" {
				sender = sdk.MustAccAddressFromBech32("elys1jeqlq99ustyug8rxgsrtmf3awzl83x535v3svykfwkkkr7049wxqgdt6ss")
				amount = sdk.NewCoins(sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 100000000), sdk.NewInt64Coin("uelys", 1000000000))
			}
			cacheCtx, write := ctx.CacheContext()
			err := app.BankKeeper.SendCoins(cacheCtx, sender, receiver, amount)
			if err == nil {
				write()
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
			Added: []string{vaulttypes.StoreKey},
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
