package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	"github.com/elys-network/elys/v7/x/commitment/types"
)

// accounting the liquid token as a claimed token in commitment module.
func (k Keeper) DepositLiquidTokensClaimed(ctx sdk.Context, denom string, amount math.Int, sender sdk.AccAddress) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "denom: %s, sender: %s", denom, sender.String())
	}

	if !assetProfile.CommitEnabled {
		return errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s, sender: %s", denom, sender.String())
	}

	depositCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// send the deposited coins to the module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, depositCoins)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "unable to send deposit tokens: %v, sender: %s", depositCoins, sender.String())
	}

	// Increase claimed amount
	commitments := k.GetCommitments(ctx, sender)
	commitments.AddClaimed(sdk.NewCoin(denom, amount))
	k.SetCommitments(ctx, commitments)

	return nil
}
