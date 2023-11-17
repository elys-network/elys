package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) BrokerOpen(goCtx context.Context, msg *types.MsgBrokerOpen) (*types.MsgBrokerOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get broker address
	brokerAddress := k.GetBrokerAddress(ctx).String()

	// check if broker is allowed
	if brokerAddress != msg.Creator {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid broker address; expected %s, got %s", brokerAddress, msg.Creator)
	}

	return k.Keeper.BrokerOpen(ctx, msg)
}
