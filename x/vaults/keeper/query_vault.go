package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/v5/x/vaults/types"
)

func (k Keeper) Vault(goCtx context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, status.Error(codes.NotFound, "vault not found")
	}

	return &types.QueryVaultResponse{Vault: vault}, nil
}

func (k Keeper) Vaults(goCtx context.Context, req *types.QueryVaultsRequest) (*types.QueryVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryVaultsResponse{Vaults: k.GetAllVaults(ctx)}, nil
}

func (k Keeper) VaultValue(goCtx context.Context, req *types.QueryVaultValue) (*types.QueryVaultValueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	usdValue, err := k.VaultUsdValue(ctx, req.VaultId)
	if err != nil {
		return nil, err
	}

	return &types.QueryVaultValueResponse{UsdValue: usdValue.Dec()}, nil
}
