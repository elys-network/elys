package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) VestTokens(ctx sdk.Context, epochIdentifier string) error {
	// Future Improvement: get all VestingTokens by denom and iterate
	k.IterateCommitments(ctx, func(commitments types.Commitments) (stop bool) {
		logger := k.Logger(ctx)

		newVestingTokens := make([]*types.VestingTokens, 0, len(commitments.VestingTokens))

		for index := len(commitments.VestingTokens) - 1; index >= 0; index-- {
			vesting := commitments.VestingTokens[index]
			vesting.CurrentEpoch = vesting.CurrentEpoch + 1
			if vesting.CurrentEpoch > vesting.NumEpochs || vesting.UnvestedAmount.IsZero() {
				continue
			}

			epochAmount := vesting.TotalAmount.Quo(sdk.NewInt(vesting.NumEpochs))

			withdrawCoins := sdk.NewCoins(sdk.NewCoin(vesting.Denom, epochAmount))

			// mint coins if vesting token is ELYS
			if vesting.Denom == ptypes.Elys {
				err := k.bankKeeper.MintCoins(ctx, types.ModuleName, withdrawCoins)
				if err != nil {
					logger.Debug(
						"unable to mint vested tokens for ELYS token",
						"vestingtokens", vesting, commitments.Creator,
					)
				}
			}

			// Send the coins to the user's account
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(commitments.Creator), withdrawCoins)
			if err != nil {
				logger.Debug(
					"unable to send vested tokens",
					"vestingtokens", vesting, commitments.Creator,
				)
			}

			vesting.UnvestedAmount = vesting.UnvestedAmount.Sub(epochAmount)
		}

		// Remove completed vesting items.
		for index := 0; index < len(commitments.VestingTokens); index++ {
			vesting := commitments.VestingTokens[index]
			if vesting.CurrentEpoch >= vesting.NumEpochs || vesting.UnvestedAmount.IsZero() {
				continue
			}

			newVestingTokens = append(newVestingTokens, vesting)
		}

		commitments.VestingTokens = newVestingTokens
		// update commitments
		k.SetCommitments(ctx, commitments)
		return false
	})

	return nil
}
