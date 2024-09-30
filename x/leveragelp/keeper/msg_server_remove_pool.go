package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) RemovePool(goCtx context.Context, msg *types.MsgRemovePool) (*types.MsgRemovePoolResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	pool, found := k.GetPool(ctx, msg.Id)

	if found {
		if pool.LeveragedLpAmount.GT(sdkmath.NewInt(0)) {
			return nil, errorsmod.Wrap(types.ErrPoolLeverageAmountNotZero, pool.LeveragedLpAmount.String())
		}
		k.DeletePool(ctx, msg.Id)
	}

	return &types.MsgRemovePoolResponse{}, nil
}
