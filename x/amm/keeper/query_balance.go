package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v4/x/amm/types"
	paramtypes "github.com/elys-network/elys/v4/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Balance(goCtx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	// Calculates balance in bank module
	// For Eden/EdenB, calculates commitment claimed amount
	balance := k.bankKeeper.GetBalance(ctx, address, req.Denom)
	if req.Denom == paramtypes.Eden || req.Denom == paramtypes.EdenB {
		commitment := k.commitmentKeeper.GetCommitments(ctx, address)
		claimed := commitment.GetClaimedForDenom(req.Denom)
		balance = sdk.NewCoin(req.Denom, claimed)
	}

	return &types.QueryBalanceResponse{
		Balance: balance,
	}, nil
}
