package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (k msgServer) UpdatePoolMultipliers(goCtx context.Context, msg *types.MsgUpdatePoolMultipliers) (*types.MsgUpdatePoolMultipliersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.UpdatePoolMultipliers(ctx, msg.PoolMultipliers)

	return &types.MsgUpdatePoolMultipliersResponse{}, nil
}
