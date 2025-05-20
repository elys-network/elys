package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/oracle/types"
)

func (k msgServer) CreateAssetInfo(goCtx context.Context, msg *types.MsgCreateAssetInfo) (*types.MsgCreateAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.GetAssetInfo(ctx, msg.Denom)

	if found {
		return nil, errors.Wrapf(types.ErrAssetWasCreated, "%s", msg.Denom)
	}

	k.Keeper.SetAssetInfo(ctx, types.AssetInfo{
		Denom:      msg.Denom,
		Display:    msg.Display,
		BandTicker: msg.BandTicker,
		ElysTicker: msg.ElysTicker,
		Decimal:    msg.Decimal,
	})

	return &types.MsgCreateAssetInfoResponse{}, nil
}
