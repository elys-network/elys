package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	LocalNetVersion = "v999.999.999"
	NewMaxBytes     = 5 * 1024 * 1024 // 5MB
)

// make sure to update these when you upgrade the version
var NextVersion = "v0.52.0"

func (app *ElysApp) setUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		version.Version,
		func(goCtx context.Context, plan upgradetypes.Plan, vm m.VersionMap) (m.VersionMap, error) {
			ctx := sdk.UnwrapSDKContext(goCtx)
			app.Logger().Info("Running upgrade handler for " + version.Version)

			if version.Version == NextVersion || version.Version == LocalNetVersion {

				// Add any logic here to run when the chain is upgraded to the new version

				// update the signing info for the validators
				signers := []string{
					"elysvalcons1f9lzcfxxu6l9yj9uf0lqjc0qa82raypnlk58ej", // Synergy Nodes
					"elysvalcons1frn2njtny6gzdjl2df9rvz3atcds2vl2fhxg8s", // Regenerator | Green Validator
				}
				for _, signer := range signers {
					addr, err := sdk.ConsAddressFromBech32(signer)
					if err != nil {
						app.Logger().Error("failed to convert signer address", "error", err)
						continue
					}
					signingInfo, err := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addr)
					if err != nil {
						app.Logger().Error("failed to get signer signing info", "signer", signer)
						continue
					}
					signingInfo.JailedUntil = ctx.BlockTime() // set jailed until to current block time
					signingInfo.Tombstoned = false
					app.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, signingInfo)
					app.Logger().Info("reset tombstoned status and jailed until date for signer", "signer", signer)
				}
				// update the unbonded status for the validators
				operators := []string{
					"elysvaloper1xesqr8vjvy34jhu027zd70ypl0nnev5ewa6r7h", // Synergy Nodes
					"elysvaloper19r0mcqdgserlx4v9htqh8erp8r2fc4ry30vl3j", // Regenerator | Green Validator
				}
				for _, operator := range operators {
					addr, err := sdk.ValAddressFromBech32(operator)
					if err != nil {
						app.Logger().Error("failed to convert operator address", "error", err)
						continue
					}
					validator, err := app.StakingKeeper.GetValidator(ctx, addr)
					if err != nil {
						app.Logger().Error("failed to get validator", "operator", operator)
						continue
					}
					for _, unbondingId := range validator.UnbondingIds {
						unbondingDelegation, err := app.StakingKeeper.GetUnbondingDelegationByUnbondingID(ctx, unbondingId)
						if err != nil {
							app.Logger().Error("failed to get unbonding delegation", "operator", operator, "unbondingId", unbondingId)
							continue
						}
						app.StakingKeeper.RemoveUnbondingDelegation(ctx, unbondingDelegation)
						app.Logger().Info("removed unbonding delegation", "operator", operator, "unbondingId", unbondingId)
					}
					validator.Jailed = false
					validator.Status = stakingtypes.Bonded
					validator.UnbondingTime = ctx.BlockTime()
					validator.UnbondingIds = []uint64{}
					app.StakingKeeper.SetValidator(ctx, validator)
					app.Logger().Info("reset unbonded status for validator", "operator", operator)
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
