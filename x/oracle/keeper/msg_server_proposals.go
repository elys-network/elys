package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/oracle/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) RemoveAssetInfo(goCtx context.Context, msg *types.MsgRemoveAssetInfo) (*types.MsgRemoveAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.RemoveAssetInfo(ctx, msg.Denom)
	return &types.MsgRemoveAssetInfoResponse{}, nil
}

func (k msgServer) AddPriceFeeders(goCtx context.Context, msg *types.MsgAddPriceFeeders) (*types.MsgAddPriceFeedersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	for _, feeder := range msg.Feeders {
		k.Keeper.SetPriceFeeder(ctx, types.PriceFeeder{
			Feeder:   feeder,
			IsActive: true,
		})
	}
	return &types.MsgAddPriceFeedersResponse{}, nil
}

func (k msgServer) RemovePriceFeeders(goCtx context.Context, msg *types.MsgRemovePriceFeeders) (*types.MsgRemovePriceFeedersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	for _, feeder := range msg.Feeders {
		k.Keeper.RemovePriceFeeder(ctx, sdk.MustAccAddressFromBech32(feeder))
	}
	return &types.MsgRemovePriceFeedersResponse{}, nil
}
