package keeper

import (
	"context"

	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
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
	poolsTVL := elystypes.ZeroDec34()
	totalTVL := math.ZeroInt()

	for _, pool := range allPools {
		tvl, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}
		poolsTVL = poolsTVL.Add(tvl)
	}
	totalTVL = totalTVL.Add(poolsTVL.ToInt())

	baseCurrencyEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, status.Error(codes.NotFound, "asset profile not found")
	}

	stableStakeTVL := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrencyEntry.Denom)
	totalTVL = totalTVL.Add(stableStakeTVL.ToInt())

	elysPrice, decimals := k.amm.GetTokenPrice(ctx, ptypes.Elys, baseCurrencyEntry.Denom)

	stakedElys := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(stakingtypes.BondedPoolName), ptypes.Elys).Amount
	stakedElysValue := elysPrice.MulInt(stakedElys).QuoInt(ammtypes.OneTokenUnit(decimals))
	totalTVL = totalTVL.Add(stakedElysValue.ToInt())

	commitmentParams := k.commitmentKeeper.GetParams(ctx)
	stakedEden := commitmentParams.TotalCommitted.AmountOf(ptypes.Eden)
	stakedEdenValue := elysPrice.MulInt(stakedEden)
	totalTVL = totalTVL.Add(stakedEdenValue.ToInt())

	return &types.QueryChainTVLResponse{
		Total:       totalTVL,
		Pools:       poolsTVL.ToInt(),
		UsdcStaking: stableStakeTVL.ToInt(),
		StakedElys:  stakedElysValue.ToInt(),
		StakedEden:  stakedEdenValue.ToInt(),
	}, nil
}
