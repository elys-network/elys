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

		for index := len(commitments.VestingTokens) - 1; index >= 0; index-- {
			vesting := commitments.VestingTokens[index]
			vesting.CurrentEpoch = vesting.CurrentEpoch + 1

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

			if vesting.CurrentEpoch == vesting.NumEpochs {
				commitments.VestingTokens = append(commitments.VestingTokens[:index], commitments.VestingTokens[index+1:]...)
			}

			// update commitments
			k.SetCommitments(ctx, commitments)
		}

		return false
	})

	return nil
}
