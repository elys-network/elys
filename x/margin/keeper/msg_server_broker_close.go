package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) BrokerClose(goCtx context.Context, msg *types.MsgBrokerClose) (*types.MsgBrokerCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get broker address
	brokerAddress := k.parameterKeeper.GetParams(ctx).BrokerAddress

	// check if broker is allowed
	if brokerAddress != msg.Creator {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid broker address; expected %s, got %s", brokerAddress, msg.Creator)
	}

	return k.Keeper.BrokerClose(ctx, msg)
}
