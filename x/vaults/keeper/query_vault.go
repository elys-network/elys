package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	tiertypes "github.com/elys-network/elys/v6/x/tier/types"
	"github.com/elys-network/elys/v6/x/vaults/types"
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

	_, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, status.Error(codes.NotFound, "vault not found")
	}

	usdValue, err := k.VaultUsdValue(ctx, req.VaultId)
	if err != nil {
		return nil, err
	}

	return &types.QueryVaultValueResponse{UsdValue: usdValue.Dec()}, nil
}

func (k Keeper) VaultPositions(goCtx context.Context, req *types.QueryVaultPositionsRequest) (*types.QueryVaultPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryVaultPositionsResponse{Positions: k.GetVaultPositions(ctx, req.VaultId)}, nil
}

func (k Keeper) GetVaultPositions(ctx sdk.Context, vaultId uint64) []types.PositionToken {
	vaultAddress := types.NewVaultAddress(vaultId)
	positions := []types.PositionToken{}
	commitments := k.commitment.GetCommitments(ctx, vaultAddress)
	// TODO: Handle zero values for denom, we should issue shares if price is not available
	for _, commitment := range commitments.CommittedTokens {
		// Pool balance
		if strings.HasPrefix(commitment.Denom, "amm/pool") {
			poolId, err := ammtypes.GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				continue
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				continue
			}
			info := k.amm.PoolExtraInfo(ctx, pool, tiertypes.OneDay)
			amount := osmomath.BigDecFromSDKInt(commitment.Amount)
			token := types.PositionToken{
				TokenDenom:    commitment.Denom,
				TokenAmount:   amount.Dec(),
				TokenUsdValue: amount.Mul(osmomath.BigDecFromDec(info.LpTokenPrice)).Quo(osmomath.BigDecFromSDKInt(ammtypes.OneShare)).Dec(),
			}
			positions = append(positions, token)
		}
	}
	// Get all balances of vault
	balances := k.bk.GetAllBalances(ctx, vaultAddress)
	for _, balance := range balances {
		token := types.PositionToken{
			TokenDenom:    balance.Denom,
			TokenAmount:   osmomath.BigDecFromSDKInt(balance.Amount).Dec(),
			TokenUsdValue: k.amm.CalculateUSDValue(ctx, balance.Denom, balance.Amount).Dec(),
		}
		positions = append(positions, token)
	}

	return positions
}
