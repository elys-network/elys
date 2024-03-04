package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) BrokerOpen(goCtx context.Context, msg *types.MsgBrokerOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// fail if msg.Owner is empty
	if msg.Owner == "" {
		return nil, errors.Wrap(types.ErrUnauthorised, "owner is not defined")
	}

	// fail if msg.Creator is not broker address
	if msg.Creator != k.parameterKeeper.GetParams(ctx).BrokerAddress {
		return nil, errors.Wrap(types.ErrUnauthorised, "creator must be broker address")
	}

	msgOpen := types.NewMsgOpen(msg.Owner, msg.Position, msg.Leverage, msg.TradingAsset, msg.Collateral, msg.TakeProfitPrice)

	return k.Keeper.Open(ctx, msgOpen, true)
}
