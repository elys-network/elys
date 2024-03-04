package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) BrokerClose(goCtx context.Context, msg *types.MsgBrokerClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// fail if msg.Owner is empty
	if msg.Owner == "" {
		return nil, errors.Wrap(types.ErrUnauthorised, "owner is not defined")
	}

	// fail if msg.Creator is not broker address
	if msg.Creator != k.parameterKeeper.GetParams(ctx).BrokerAddress {
		return nil, errors.Wrap(types.ErrUnauthorised, "creator must be broker address")
	}

	msgClose := types.NewMsgClose(msg.Owner, msg.Id, msg.Amount)

	return k.Keeper.Close(ctx, msgClose)
}
