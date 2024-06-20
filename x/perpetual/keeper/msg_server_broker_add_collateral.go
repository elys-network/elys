package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) BrokerAddCollateral(goCtx context.Context, msg *types.MsgBrokerAddCollateral) (*types.MsgAddCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// fail if msg.Owner is empty
	if msg.Owner == "" {
		return nil, errors.Wrap(types.ErrUnauthorised, "owner is not defined")
	}

	// fail if msg.Creator is not broker address
	if msg.Creator != k.parameterKeeper.GetParams(ctx).BrokerAddress {
		return nil, errors.Wrap(types.ErrUnauthorised, "creator must be broker address")
	}

	// TODO
	//msgAddCollateral := types.NewMsgAddCollateral(msg.Owner, msg.Amount, uint64(msg.Id))
	//return k.Keeper.Add(ctx, msgAddCollateral)
	return &types.MsgAddCollateralResponse{}, nil
}
