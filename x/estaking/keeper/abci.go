package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// EndBlocker of incentive module
func (k Keeper) EndBlocker(ctx sdk.Context) error {
	// Rewards distribution
	err := k.ProcessRewardsDistribution(ctx)
	if err != nil {
		return err
	}
	// Burn EdenB tokens if staking changed
	err = k.BurnEdenBIfElysStakingReduced(ctx)
	if err != nil {
		return err
	}
	return nil
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

func (k Keeper) BurnEdenBIfElysStakingReduced(ctx sdk.Context) error {
	addrs := k.GetAllElysStakeChange(ctx)

	// Handle addresses recorded on AfterDelegationModified
	// This hook is exposed for genesis delegations as well
	for _, delAddr := range addrs {
		err := k.BurnEdenBFromElysUnstaking(ctx, delAddr)
		if err != nil {
			return err
		}
		k.TakeDelegationSnapshot(ctx, delAddr)
		k.RemoveElysStakeChange(ctx, delAddr)
	}
	return nil
}

// Rewards distribution
func (k Keeper) ProcessRewardsDistribution(ctx sdk.Context) error {
	// Read tokenomics time based inflation params and update incentive module params.
	k.ProcessUpdateIncentiveParams(ctx)

	return k.UpdateStakersRewards(ctx)
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

	providerEdenAmount := osmomath.BigDecFromSDKInt(stakersEdenAmount).Mul(params.GetBigDecProviderStakingRewardsPortion()).Dec().TruncateInt()
	err := k.commKeeper.MintCoins(ctx, ccvconsumertypes.ConsumerToSendToProviderName, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, providerEdenAmount)))
	if err != nil {
		return err
	}

	stakersEdenAmountAfterProvider := stakersEdenAmount.Sub(providerEdenAmount)

	totalElysEdenEdenBStake, err := k.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	totalElysEdenStake, err := k.TotalBondedElysEdenTokens(ctx)
	if err != nil {
		return err
	}

	// Maximum eden APR - 30% by default
	stakersMaxEdenAmount := osmomath.BigDecFromDec(params.MaxEdenRewardAprStakers).
		Mul(osmomath.BigDecFromSDKInt(totalElysEdenStake)).
		QuoInt64(totalBlocksPerYear)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	stakersEdenAmountForGovernors := math.MinInt(stakersEdenAmountAfterProvider, stakersMaxEdenAmount.Dec().TruncateInt())

	// EdenB should be mint based on Elys + Eden staked (should exclude edenB staked)
	stakersEdenBAmount := osmomath.BigDecFromSDKInt(totalElysEdenEdenBStake).
		Mul(params.GetBigDecEdenBoostApr()).
		QuoInt64(totalBlocksPerYear).
		Dec().
		RoundInt()

	consumerCoins := sdk.NewCoins(
		sdk.NewCoin(ptypes.Eden, stakersEdenAmountForGovernors),
		sdk.NewCoin(ptypes.EdenB, stakersEdenBAmount),
	)
	return k.commKeeper.MintCoins(ctx, ccvconsumertypes.ConsumerRedistributeName, consumerCoins.Sort())
}
