package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// EndBlocker of incentive module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// Rewards distribution
	k.ProcessRewardsDistribution(ctx)
	// Burn EdenB tokens if staking changed
	k.BurnEdenBIfElysStakingReduced(ctx)
}

func (k Keeper) TakeDelegationSnapshot(ctx sdk.Context, addr sdk.AccAddress) {
	// Calculate delegated amount per delegator
	delAmount := k.CalcDelegationAmount(ctx, addr)

	elysStaked := types.ElysStaked{
		Address: addr.String(),
		Amount:  delAmount,
	}

	// Set Elys staked amount
	k.SetElysStaked(ctx, elysStaked)
}

func (k Keeper) BurnEdenBIfElysStakingReduced(ctx sdk.Context) {
	addrs := k.GetAllElysStakeChange(ctx)

	// Handle addresses recorded on AfterDelegationModified
	// This hook is exposed for genesis delegations as well
	for _, delAddr := range addrs {
		k.BurnEdenBFromElysUnstaking(ctx, delAddr)
		k.TakeDelegationSnapshot(ctx, delAddr)
		k.RemoveElysStakeChange(ctx, delAddr)
	}
}

// Rewards distribution
func (k Keeper) ProcessRewardsDistribution(ctx sdk.Context) {
	// Read tokenomics time based inflation params and update incentive module params.
	k.ProcessUpdateIncentiveParams(ctx)

	err := k.UpdateStakersRewards(ctx)
	if err != nil {
		ctx.Logger().Error("Failed to update staker rewards unclaimed", "error", err)
	}
}

func (k Keeper) ProcessUpdateIncentiveParams(ctx sdk.Context) {
	// Non-linear inflation per year happens and this includes yearly inflation data
	listTimeBasedInflations := k.tokenomicsKeeper.GetAllTimeBasedInflation(ctx)
	if len(listTimeBasedInflations) == 0 {
		return
	}

	params := k.GetParams(ctx)

	for _, inflation := range listTimeBasedInflations {
		// Finding only current inflation data - and skip rest
		if inflation.StartBlockHeight > uint64(ctx.BlockHeight()) || inflation.EndBlockHeight < uint64(ctx.BlockHeight()) {
			continue
		}

		totalBlocks := inflation.EndBlockHeight - inflation.StartBlockHeight + 1

		// If totalBlocks is zero, we skip this inflation to avoid division by zero
		if totalBlocks == 0 {
			continue
		}

		// ------------- Stakers parameter -------------
		blocksDistributed := ctx.BlockHeight() - int64(inflation.StartBlockHeight)
		params.StakeIncentives = &types.IncentiveInfo{
			EdenAmountPerYear: math.NewInt(int64(inflation.Inflation.IcsStakingRewards)),
			BlocksDistributed: blocksDistributed,
		}
		k.SetParams(ctx, params)
		return
	}

	params.StakeIncentives = nil
	k.SetParams(ctx, params)
}

func (k Keeper) UpdateStakersRewards(ctx sdk.Context) error {
	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// USDC amount in math.LegacyDec type
	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	totalFeesCollected := k.commKeeper.GetAllBalances(ctx, feeCollectorAddr)
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(totalFeesCollected...)
	dexRevenueStakersAmount := gasFeeCollectedDec.AmountOf(baseCurrency)

	// Calculate eden amount per block
	params := k.GetParams(ctx)
	stakeIncentive := params.StakeIncentives

	// Ensure totalBlocksPerYear are not zero to avoid division by zero
	totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)

	// Calculate
	edenAmountPerYear := math.ZeroInt()
	if stakeIncentive != nil && stakeIncentive.EdenAmountPerYear.IsPositive() {
		edenAmountPerYear = stakeIncentive.EdenAmountPerYear
	}
	stakersEdenAmount := edenAmountPerYear.Quo(math.NewInt(totalBlocksPerYear))

	// Maximum eden APR - 30% by default
	totalElysEdenEdenBStake, err := k.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	stakersMaxEdenAmount := params.MaxEdenRewardAprStakers.
		MulInt(totalElysEdenEdenBStake).
		QuoInt64(totalBlocksPerYear)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	stakersEdenAmount = math.MinInt(stakersEdenAmount, stakersMaxEdenAmount.TruncateInt())

	stakersEdenBAmount := math.LegacyNewDecFromInt(totalElysEdenEdenBStake).
		Mul(params.EdenBoostApr).
		QuoInt64(totalBlocksPerYear).
		RoundInt()

	// Set block number and total dex rewards given
	params.DexRewardsStakers.NumBlocks = 1
	params.DexRewardsStakers.Amount = dexRevenueStakersAmount
	k.SetParams(ctx, params)

	coins := sdk.NewCoins(
		sdk.NewCoin(ptypes.Eden, stakersEdenAmount),
		sdk.NewCoin(ptypes.EdenB, stakersEdenBAmount),
	)
	return k.commKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins.Sort())
}
