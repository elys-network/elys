package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AirDrop(goCtx context.Context, req *types.QueryAirDropRequest) (*types.QueryAirDropResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	atomStaker := k.GetAtomStaker(ctx, address)
	nftHolder := k.GetNFTHolder(ctx, address)
	cadet := k.GetCadet(ctx, address)
	governor := k.GetGovernor(ctx, address)
	return &types.QueryAirDropResponse{
		AtomStaking: atomStaker.Amount,
		Cadet:       cadet.Amount,
		NftHolder:   nftHolder.Amount,
		Governor:    governor.Amount,
	}, nil
}
