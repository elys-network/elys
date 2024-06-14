package app

import (
	"fmt"

	"cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

			if version.Version == "v0.32.0" || version.Version == "v999.999.999" {
				// Since invariant is broken sending missing amount to bonded pool
				sumOfValTokens := math.ZeroInt()
				app.EstakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
					sumOfValTokens = sumOfValTokens.Add(validator.GetTokens())
					return false
				})

				totalBondedTokens := app.EstakingKeeper.TotalBondedTokens(ctx)
				if totalBondedTokens.LT(sumOfValTokens) {
					missingAmount := sumOfValTokens.Sub(totalBondedTokens)
					missingCoins := sdk.Coins{sdk.NewCoin("uelys", missingAmount)}
					err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, missingCoins)
					if err != nil {
						panic(err)
					}
					err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, stakingtypes.BondedPoolName, missingCoins)
					if err != nil {
						panic(err)
					}
				}

				app.Logger().Info("Deleting proposals with ID <= 193")
				store := ctx.KVStore(app.keys[govtypes.StoreKey])
				for i := uint64(1); i <= 193; i++ {
					store.Delete(govtypes.ProposalKey(i))
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
			Added: []string{"membershiptier"},
		}
		fmt.Printf("Setting store loader with height %d and store upgrades: %+v\n", upgradeInfo.Height, storeUpgrades)

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))

		fmt.Println("Store loader set successfully.")
	} else {
		fmt.Println("No need to load upgrade store.")
	}
}

func shouldLoadUpgradeStore(app *ElysApp, upgradeInfo upgradetypes.Plan) bool {
	currentHeight := app.LastBlockHeight()
	fmt.Printf("Current block height: %d, Upgrade height: %d\n", currentHeight, upgradeInfo.Height)
	return upgradeInfo.Name == version.Version && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height)
}
