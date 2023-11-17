package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) BrokerOpen(ctx sdk.Context, msg *types.MsgBrokerOpen) (*types.MsgBrokerOpenResponse, error) {
	msgOpen := types.NewMsgOpen(msg.Owner, msg.CollateralAsset, msg.CollateralAmount, msg.BorrowAsset, msg.Position, msg.Leverage, msg.TakeProfitPrice)

	res, err := k.Open(ctx, msgOpen)
	if err != nil {
		return nil, err
	}

	return &types.MsgBrokerOpenResponse{
		Id: res.Id,
	}, nil
}
