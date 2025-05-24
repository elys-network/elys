package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
)

// accounting the liquid token as a claimed token in commitment module.
func (k Keeper) DepositLiquidTokensClaimed(ctx sdk.Context, denom string, amount math.Int, sender sdk.AccAddress) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.CommitEnabled {
		return errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s", denom)
	}

	depositCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// send the deposited coins to the module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, depositCoins)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, fmt.Sprintf("unable to send deposit tokens: %v", depositCoins))
	}

	// Increase claimed amount
	commitments := k.GetCommitments(ctx, sender)
	commitments.AddClaimed(sdk.NewCoin(denom, amount))
	k.SetCommitments(ctx, commitments)

	return nil
}
