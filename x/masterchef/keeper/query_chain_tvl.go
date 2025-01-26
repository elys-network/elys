package keeper

import (
	"context"
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"

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
	poolsTVL := math.LegacyZeroDec()
	totalTVL := math.ZeroInt()

	for _, pool := range allPools {
		tvl, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}
		poolsTVL = poolsTVL.Add(tvl)
	}
	totalTVL = totalTVL.Add(poolsTVL.TruncateInt())

	baseCurrencyEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, status.Error(codes.NotFound, "asset profile not found")
	}

	baseCurrencyPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrencyEntry.Denom)

	stableStakeTVL := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrencyEntry.Denom)
	totalTVL = totalTVL.Add(stableStakeTVL.TruncateInt())

	elysPrice := k.amm.GetTokenPrice(ctx, ptypes.Elys, baseCurrencyEntry.Denom)

	stakedElys := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(stakingtypes.BondedPoolName), ptypes.Elys).Amount
	stakedElysValue := elysPrice.MulInt(stakedElys)
	totalTVL = totalTVL.Add(stakedElysValue.TruncateInt())

	commitmentParams := k.commitmentKeeper.GetParams(ctx)
	stakedEden := commitmentParams.TotalCommitted.AmountOf(ptypes.Eden)
	stakedEdenValue := elysPrice.MulInt(stakedEden)
	totalTVL = totalTVL.Add(stakedEdenValue.TruncateInt())

	stableStakeBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(stablestaketypes.ModuleName), baseCurrencyEntry.Denom)

	return &types.QueryChainTVLResponse{
		Total:       totalTVL,
		Pools:       poolsTVL.TruncateInt(),
		UsdcStaking: stableStakeTVL.TruncateInt(),
		StakedElys:  stakedElysValue.TruncateInt(),
		StakedEden:  stakedEdenValue.TruncateInt(),
		NetStakings: sdk.NewCoins(sdk.NewCoin(baseCurrencyEntry.DisplayName, (stableStakeBalance.Amount.ToLegacyDec().Mul(baseCurrencyPrice).TruncateInt()))),
	}, nil
}
