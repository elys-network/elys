package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) BrokerClose(ctx sdk.Context, msg *types.MsgBrokerClose) (*types.MsgBrokerCloseResponse, error) {
	msgClose := types.NewMsgClose(msg.Owner, msg.Id, msg.Amount)

	res, err := k.Close(ctx, msgClose)
	if err != nil {
		return nil, err
	}

	return &types.MsgBrokerCloseResponse{
		Id: res.Id,
	}, nil
}
