package keeper

import (
	"errors"

	errorsmod "cosmossdk.io/errors"

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
}

// Rewards distribution
func (k Keeper) ProcessRewardsDistribution(ctx sdk.Context) {
	// Read tokenomics time based inflation params and update incentive module params.
	if !k.ProcessUpdateIncentiveParams(ctx) {
		ctx.Logger().Error("Invalid tokenomics params", "error", errors.New("invalid tokenomics params"))
		return
	}

	canDistribute := k.CanDistributeStakingRewards(ctx)
	if canDistribute {
		err := k.UpdateStakersRewards(ctx)
		if err != nil {
			ctx.Logger().Error("Failed to update staker rewards unclaimed", "error", err)
		}
	}
}

func (k Keeper) ProcessUpdateIncentiveParams(ctx sdk.Context) bool {
	// Non-linear inflation per year happens and this includes yearly inflation data
	listTimeBasedInflations := k.tokenomicsKeeper.GetAllTimeBasedInflation(ctx)
	if len(listTimeBasedInflations) < 1 {
		return false
	}

	params := k.GetParams(ctx)

	for _, inflation := range listTimeBasedInflations {
		// Finding only current inflation data - and skip rest
		if inflation.StartBlockHeight > uint64(ctx.BlockHeight()) || inflation.EndBlockHeight < uint64(ctx.BlockHeight()) {
			continue
		}

		totalBlocksPerYear := sdk.NewInt(int64(inflation.EndBlockHeight - inflation.StartBlockHeight + 1))

		// If totalDistributionEpochPerYear is zero, we skip this inflation to avoid division by zero
		if totalBlocksPerYear == sdk.ZeroInt() {
			continue
		}

		// ------------- Stakers parameter -------------
		blocksDistributed := sdk.NewInt(ctx.BlockHeight() - int64(inflation.StartBlockHeight))
		incentiveInfo := types.IncentiveInfo{
			EdenAmountPerYear:      sdk.NewInt(int64(inflation.Inflation.IcsStakingRewards)),
			DistributionStartBlock: sdk.NewInt(int64(inflation.StartBlockHeight)),
			TotalBlocksPerYear:     totalBlocksPerYear,
			BlocksDistributed:      blocksDistributed,
		}

		if params.StakeIncentives == nil {
			params.StakeIncentives = &incentiveInfo
		} else {
			// If any of block number related parameter changed, we re-calculate the current epoch
			if params.StakeIncentives.DistributionStartBlock != incentiveInfo.DistributionStartBlock ||
				params.StakeIncentives.TotalBlocksPerYear != incentiveInfo.TotalBlocksPerYear {
				params.StakeIncentives.BlocksDistributed = blocksDistributed
			}
			params.StakeIncentives.EdenAmountPerYear = incentiveInfo.EdenAmountPerYear
			params.StakeIncentives.DistributionStartBlock = incentiveInfo.DistributionStartBlock
			params.StakeIncentives.TotalBlocksPerYear = incentiveInfo.TotalBlocksPerYear
		}
		break
	}

	k.SetParams(ctx, params)
	return true
}

func (k Keeper) CanDistributeStakingRewards(ctx sdk.Context) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false
	}

	// If we don't have enough params
	if params.StakeIncentives == nil {
		return false
	}

	// Incentive params initialize
	stakeIncentive := params.StakeIncentives

	curBlockHeight := sdk.NewInt(ctx.BlockHeight())
	if stakeIncentive.DistributionStartBlock.GT(curBlockHeight) {
		return false
	}

	// Increase current epoch of Stake incentive param
	stakeIncentive.BlocksDistributed = stakeIncentive.BlocksDistributed.Add(sdk.OneInt())
	if stakeIncentive.BlocksDistributed.GTE(stakeIncentive.TotalBlocksPerYear) || curBlockHeight.GT(stakeIncentive.TotalBlocksPerYear.Add(stakeIncentive.DistributionStartBlock)) {
		params.StakeIncentives = nil
		k.SetParams(ctx, params)
		return false
	}

	params.StakeIncentives.BlocksDistributed = stakeIncentive.BlocksDistributed
	k.SetParams(ctx, params)

	// return found, stake incentive params
	return true
}

// Update unclaimed token amount
// Called back through epoch hook
func (k Keeper) UpdateStakersRewards(ctx sdk.Context) error {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// USDC amount in sdk.Dec type
	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	totalFeesCollected := k.commKeeper.GetAllBalances(ctx, feeCollectorAddr)
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(totalFeesCollected...)
	dexRevenueStakersAmount := gasFeeCollectedDec.AmountOf(baseCurrency)

	// Calculate eden amount per epoch
	params := k.GetParams(ctx)
	stakeIncentive := params.StakeIncentives

	// Ensure stakeIncentive.TotalBlocksPerYear or stakeIncentive.EpochNumBlocks are not zero to avoid division by zero
	if stakeIncentive.TotalBlocksPerYear.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	// Calculate
	epochStakersEdenAmount := stakeIncentive.EdenAmountPerYear.
		Quo(stakeIncentive.TotalBlocksPerYear)

	// Maximum eden based per distribution epoch on maximum APR - 30% by default
	// Allocated for staking per day = (0.3/365)* ( total elys staked + total Eden committed + total Eden boost committed)

	totalElysEdenEdenBStake := k.TotalBondedTokens(ctx)

	epochStakersMaxEdenAmount := params.MaxEdenRewardAprStakers.
		MulInt(totalElysEdenEdenBStake).
		QuoInt(stakeIncentive.TotalBlocksPerYear)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	epochStakersEdenAmount = sdk.MinInt(epochStakersEdenAmount, epochStakersMaxEdenAmount.TruncateInt())

	epochStakersEdenBAmount := sdk.NewDecFromInt(totalElysEdenEdenBStake).
		Mul(params.EdenBoostApr).
		QuoInt(stakeIncentive.TotalBlocksPerYear).
		RoundInt()

	// Set block number and total dex rewards given
	params.DexRewardsStakers.NumBlocks = sdk.OneInt()
	params.DexRewardsStakers.Amount = dexRevenueStakersAmount

	k.SetParams(ctx, params)

	coins := sdk.NewCoins(
		sdk.NewCoin(ptypes.Eden, epochStakersEdenAmount),
		sdk.NewCoin(ptypes.EdenB, epochStakersEdenBAmount),
	)
	return k.commKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins.Sort())
}
