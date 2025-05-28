package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
	"github.com/elys-network/elys/v5/x/masterchef/types"
)

// AfterPoolCreated is called after CreatePool
func (k Keeper) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64) error {
	_, found := k.GetPoolInfo(ctx, poolId)
	if found {
		return nil
	}
	// Initiate a new pool info
	poolInfo := types.PoolInfo{
		// reward amount
		PoolId: poolId,
		// reward wallet address
		RewardWallet: ammtypes.NewPoolRevenueAddress(poolId).String(),
		// multiplier for lp rewards
		Multiplier: math.LegacyNewDec(1),
		// Eden APR, updated at every distribution
		EdenApr: math.LegacyZeroDec(),
		// Dex APR, updated at every distribution
		DexApr: math.LegacyZeroDec(),
		// Gas APR, updated at every distribution
		GasApr: math.LegacyZeroDec(),
		// External Incentive APR, updated at every distribution
		ExternalIncentiveApr: math.LegacyZeroDec(),
		// external reward denoms on the pool
		ExternalRewardDenoms: []string{},
		// Eden rewards is false by default
		EnableEdenRewards: false,
	}
	k.SetPoolInfo(ctx, poolInfo)
	return nil
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (k Keeper) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, enterCoins sdk.Coins, shareOutAmount math.Int) error {
	k.AfterDeposit(ctx, poolId, sender, shareOutAmount)
	return nil
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (k Keeper) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount math.Int, exitCoins sdk.Coins) error {
	k.AfterWithdraw(ctx, poolId, sender, shareInAmount)

	return nil
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (k Keeper) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) error {
	return nil
}

// Hooks wrapper struct for incentive keeper
type AmmHooks struct {
	k Keeper
}

var _ ammtypes.AmmHooks = AmmHooks{}

// Return the wrapper struct
func (k Keeper) AmmHooks() AmmHooks {
	return AmmHooks{k}
}

// AfterPoolCreated is called after CreatePool
func (h AmmHooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool) error {
	return h.k.AfterPoolCreated(ctx, sender, pool.PoolId)
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (h AmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, enterCoins sdk.Coins, shareOutAmount math.Int) error {
	return h.k.AfterJoinPool(ctx, sender, pool.PoolId, enterCoins, shareOutAmount)
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (h AmmHooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, shareInAmount math.Int, exitCoins sdk.Coins) error {
	return h.k.AfterExitPool(ctx, sender, pool.PoolId, shareInAmount, exitCoins)
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (h AmmHooks) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, input sdk.Coins, output sdk.Coins) error {
	return h.k.AfterSwap(ctx, sender, pool.PoolId, input, output)
}
