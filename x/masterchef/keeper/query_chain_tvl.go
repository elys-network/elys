package keeper

import (
	"context"
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ChainTVL(goCtx context.Context, req *types.QueryChainTVLRequest) (*types.QueryChainTVLResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	allPools := k.amm.GetAllPool(ctx)
	totalTVL := math.LegacyZeroDec()

	for _, pool := range allPools {
		tvl, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}
		totalTVL = totalTVL.Add(tvl)
	}

	baseCurrencyEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, status.Error(codes.NotFound, "asset profile not found")
	}

	stableStakeTVL := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrencyEntry.Denom)
	totalTVL = totalTVL.Add(stableStakeTVL)

	elysPrice := k.amm.GetTokenPrice(ctx, ptypes.Elys, baseCurrencyEntry.Denom)

	stakedElys := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(stakingtypes.BondedPoolName), ptypes.Elys).Amount
	stakedElysValue := elysPrice.MulInt(stakedElys)
	totalTVL = totalTVL.Add(stakedElysValue)

	commitmentParams := k.commitmentKeeper.GetParams(ctx)
	stakedEden := commitmentParams.TotalCommitted.AmountOf(ptypes.Eden)
	stakedEdenValue := elysPrice.MulInt(stakedEden)
	totalTVL = totalTVL.Add(stakedEdenValue)

	return &types.QueryChainTVLResponse{Total: totalTVL.TruncateInt()}, nil
}
