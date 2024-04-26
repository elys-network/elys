package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
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

			if version.Version == "v0.31.0" {
				app.Logger().Info("Deleting proposals with ID <= 185")
				store := ctx.KVStore(app.keys[govtypes.StoreKey])
				for i := uint64(1); i <= 185; i++ {
					store.Delete(govtypes.ProposalKey(i))
				}

				// update the signing info for the validators
				signers := []string{
					"elysvalcons1j7047ewlfa75dv0q93lnqkctr9afgfayyvmhc4", // euphoria
					"elysvalcons1a58n8t00elj7g4v8lm7rd9q06xu4nz3dgy723q", // shangrila
					"elysvalcons1t0cm443g88ns9rl7ac45a5u9cs54thtww7w4ag", // ottersync
				}
				for _, signer := range signers {
					addr, err := sdk.ConsAddressFromBech32(signer)
					if err != nil {
						app.Logger().Error("failed to convert signer address", "error", err)
						continue
					}
					signingInfo, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addr)
					if !found {
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
					"elysvaloper1cduy8wep22mdnsmml8w9gn94ek8hqnsdffy098", // euphoria
					"elysvaloper17wc3s7am5qgjk4pm0k96kg6laxq8hkyq0dzq5n", // shangrila
					"elysvaloper1uf8m4jga7akv25lf2lem50alu3kxdd8lzlvak6", // ottersync
				}
				for _, operator := range operators {
					addr, err := sdk.ValAddressFromBech32(operator)
					if err != nil {
						app.Logger().Error("failed to convert operator address", "error", err)
						continue
					}
					validator, found := app.StakingKeeper.GetValidator(ctx, addr)
					if !found {
						app.Logger().Error("failed to get validator", "operator", operator)
						continue
					}

					for _, unbondingId := range validator.UnbondingIds {
						unbondingDelegation, found := app.StakingKeeper.GetUnbondingDelegationByUnbondingID(ctx, unbondingId)
						if !found {
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

func loadUpgradeStore(app *ElysApp) {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if shouldLoadUpgradeStore(app, upgradeInfo) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{distrtypes.StoreKey, mastercheftypes.StoreKey, estakingtypes.StoreKey},
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
