package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// accounting the liquid token as a claimed token in commitment module.
func (k Keeper) DepositLiquidTokensClaimed(ctx sdk.Context, denom string, amount sdk.Int, sender string) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.CommitEnabled {
		return errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s", denom)
	}

	depositCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	addr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// send the deposited coins to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, depositCoins)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, fmt.Sprintf("unable to send deposit tokens: %v", depositCoins))
	}

	// Increase claimed amount
	commitments := k.GetCommitments(ctx, sender)
	commitments.AddClaimed(sdk.NewCoin(denom, amount))
	k.SetCommitments(ctx, commitments)

	return nil
}
