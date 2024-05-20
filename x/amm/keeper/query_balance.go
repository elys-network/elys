package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Balance(goCtx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	denom := req.Denom
	addr := req.Address
	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	// Calculates balance in bank module
	// For Eden/EdenB, calculates commitment claimed amount
	balance := k.bankKeeper.GetBalance(ctx, address, denom)
	if denom == paramtypes.Eden || denom == paramtypes.EdenB {
		commitment := k.commitmentKeeper.GetCommitments(ctx, addr)
		claimed := commitment.GetClaimedForDenom(denom)
		balance = sdk.NewCoin(denom, claimed)
	}

	return &types.QueryBalanceResponse{
		Balance: balance,
	}, nil
}
