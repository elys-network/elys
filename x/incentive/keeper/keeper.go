package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeKey     storetypes.StoreKey
		memKey       storetypes.StoreKey
		paramstore   paramtypes.Subspace
		cmk          types.CommitmentKeeper
		stk          types.StakingKeeper
		tci          *types.TotalCommitmentInfo
		authKeeper   types.AccountKeeper
		bankKeeper   types.BankKeeper
		amm          types.AmmKeeper
		oracleKeeper types.OracleKeeper

		feeCollectorName    string // name of the FeeCollector ModuleAccount
		dexRevCollectorName string // name of the Dex Revenue ModuleAccount
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ck types.CommitmentKeeper,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	amm types.AmmKeeper,
	ok types.OracleKeeper,
	feeCollectorName string,
	dexRevCollectorName string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		paramstore:          ps,
		cmk:                 ck,
		stk:                 sk,
		tci:                 &types.TotalCommitmentInfo{},
		feeCollectorName:    feeCollectorName,
		dexRevCollectorName: dexRevCollectorName,
		authKeeper:          ak,
		bankKeeper:          bk,
		amm:                 amm,
		oracleKeeper:        ok,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Update unclaimed token amount
// Called back through epoch hook
func (k Keeper) UpdateRewardsUnclaimed(ctx sdk.Context, epochIdentifier string, stakeIncentive types.IncentiveInfo, lpIncentive types.IncentiveInfo) {
	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx)

	// Collect DEX revenue while tracking 65% of it for LPs reward calculation
	// Assume these are collected in USDC
	dexRevenue, dexRevenueForLps := k.CollectDEXRevenue(ctx)

	// Calculate each portion of DEX revenue - stakers, LPs
	dexRevenueDec := sdk.NewDecCoinsFromCoins(dexRevenue...)
	dexRevenueForStakers := dexRevenueDec.Sub(dexRevenueForLps)

	// Calculate each portion of Gas fees collected - stakers, LPs
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(k.tci.TotalFeesCollected...)
	rewardPortionForLps := k.GetDEXRewardPortionForLPs(ctx)
	gasFeesForLps := gasFeeCollectedDec.MulDecTruncate(rewardPortionForLps)
	gasFeesForStakers := gasFeeCollectedDec.Sub(gasFeesForLps)

	// Sum Dex revenue for stakers + Gas fees for stakers and name it dex Revenus for stakers
	// But won't sum dex revenue for LPs and gas fees for LPs as the LP revenue will be rewared by pool.
	dexRevenueForStakers = dexRevenueForStakers.Add(gasFeesForStakers...)

	// Fund community pool based on the communtiy tax
	dexRevenueRemainedForStakers := k.UpdateCommunityPool(ctx, dexRevenueForStakers)

	// TODO:
	// Dummy denom for USDC
	// USDC amount in sdk.Dec type
	dexRevenueLPsAmt := dexRevenueForLps.AmountOf(ptypes.BaseCurrency)
	dexRevenueStakersAmt := dexRevenueRemainedForStakers.AmountOf(ptypes.BaseCurrency)
	gasFeesLPsAmt := gasFeesForLps.AmountOf(ptypes.BaseCurrency)

	// Calculate eden amount per epoch
	edenAmountPerEpochStakers := stakeIncentive.Amount.Quo(sdk.NewInt(stakeIncentive.NumEpochs))
	edenAmountPerEpochLPs := lpIncentive.Amount.Quo(sdk.NewInt(lpIncentive.NumEpochs))
	edenBoostAPR := stakeIncentive.EdenBoostApr

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have mulitplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL := k.CalculateProxyTVL(ctx)

	totalEdenGiven := sdk.ZeroInt()
	totalEdenGivenLP := sdk.ZeroInt()
	totalRewardsGiven := sdk.ZeroInt()
	totalRewardsGivenLP := sdk.ZeroInt()
	// Process to increase uncomitted token amount of Eden & Eden boost
	k.cmk.IterateCommitments(
		ctx, func(commitments ctypes.Commitments) bool {
			// Commitment owner
			creator := commitments.Creator
			_, err := sdk.AccAddressFromBech32(creator)
			if err != nil {
				// This could be validator address
				return false
			}

			rewardsByElysStaking := sdk.NewCoins()
			rewardsByEdenCommitted := sdk.NewCoins()
			rewardsByEdenBCommitted := sdk.NewCoins()
			rewardsByUSDCDeposit := sdk.NewCoins()

			// Calculate delegated amount per delegator
			delegatedAmt := k.CalculateDelegatedAmount(ctx, creator)

			// Calculate new unclaimed Eden tokens from Eden & Eden boost committed, Dex rewards distribution
			// Distribute gas fees to stakers

			// Calculate new unclaimed Eden tokens from Elys staked
			newUnclaimedEdenTokens, dexRewards, dexRewardsByStakers := k.CalculateRewardsForStakersByElysStaked(ctx, delegatedAmt, edenAmountPerEpochStakers, dexRevenueStakersAmt)
			totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			newEdenFromElysStaking := newUnclaimedEdenTokens
			newDexRewardFromElysStaking := dexRewardsByStakers.TruncateInt()

			// Calculate new unclaimed Eden tokens from Eden committed
			edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)
			newUnclaimedEdenTokens, dexRewards = k.CalculateRewardsForStakersByCommitted(ctx, edenCommitted, edenAmountPerEpochStakers, dexRevenueStakersAmt)
			totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Sub bucket
			rewardsByEdenCommitted = rewardsByEdenCommitted.Add(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))
			rewardsByEdenCommitted = rewardsByEdenCommitted.Add(sdk.NewCoin(ptypes.BaseCurrency, dexRewards))

			// Calculate new unclaimed Eden tokens from Eden Boost committed
			edenBoostCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
			newUnclaimedEdenTokens, dexRewards = k.CalculateRewardsForStakersByCommitted(ctx, edenBoostCommitted, edenAmountPerEpochStakers, dexRevenueStakersAmt)
			totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Sub bucket
			rewardsByEdenBCommitted = rewardsByEdenBCommitted.Add(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))
			rewardsByEdenBCommitted = rewardsByEdenBCommitted.Add(sdk.NewCoin(ptypes.BaseCurrency, dexRewards))

			// Calculate new unclaimed Eden tokens from LpTokens committed, Dex rewards distribution
			// Distribute gas fees to LPs
			newUnclaimedEdenTokensLp, dexRewardsLp := k.CalculateRewardsForLPs(ctx, totalProxyTVL, commitments, edenAmountPerEpochLPs, gasFeesLPsAmt)
			totalEdenGivenLP = totalEdenGivenLP.Add(newUnclaimedEdenTokensLp)
			totalRewardsGivenLP = totalRewardsGivenLP.Add(dexRewardsLp)

			// Calculate new unclaimed Eden tokens from stable stake LpTokens committed, Dex rewards distribution
			// Distribute gas fees to LPs
			newUnclaimedEdenTokensStableLp, dexRewardsStableLp := k.CalculateRewardsForLPs(ctx, totalProxyTVL, commitments, edenAmountPerEpochLPs, gasFeesLPsAmt)
			totalEdenGivenLP = totalEdenGivenLP.Add(newUnclaimedEdenTokensStableLp)
			totalRewardsGivenLP = totalRewardsGivenLP.Add(dexRewardsStableLp)

			// Sub bucket
			rewardsByUSDCDeposit = rewardsByUSDCDeposit.Add(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokensStableLp))
			rewardsByUSDCDeposit = rewardsByUSDCDeposit.Add(sdk.NewCoin(ptypes.BaseCurrency, dexRewardsStableLp))

			// Calculate the total Eden unclaimed amount
			newUnclaimedEdenTokens = newUnclaimedEdenTokens.Add(newUnclaimedEdenTokensLp)

			// Give commission to validators ( Eden from stakers and Dex rewards from stakers. )
			edenCommissionGiven, dexRewardsCommissionGiven := k.GiveCommissionToValidators(ctx, creator, delegatedAmt, newUnclaimedEdenTokens, dexRewardsByStakers)

			// Minus the given amount and increase with the remains only
			newUnclaimedEdenTokens = newUnclaimedEdenTokens.Sub(edenCommissionGiven)

			// Plus LpDexRewards and minus commission given
			dexRewards = dexRewards.Add(dexRewardsLp).Sub(dexRewardsCommissionGiven)

			// We should deduct validator commissions from "reward by elys staking sub bucket"
			newEdenFromElysStaking = newEdenFromElysStaking.Sub(edenCommissionGiven)
			newDexRewardFromElysStaking = newDexRewardFromElysStaking.Sub(dexRewardsCommissionGiven)

			// Add Eden rewards from Elys staking
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(ptypes.Eden, newEdenFromElysStaking))
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(ptypes.BaseCurrency, newDexRewardFromElysStaking))

			// Calculate new unclaimed Eden-Boost tokens for staker and Eden token holders
			newUnclaimedEdenBoostTokens, newUnclaimedEdenBoostFromElysStaking, newUnclaimedEdenBoostFromEdenCommited := k.CalculateEdenBoostRewards(ctx, delegatedAmt, commitments, epochIdentifier, edenBoostAPR)
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(ptypes.EdenB, newUnclaimedEdenBoostFromElysStaking))
			rewardsByEdenCommitted = rewardsByEdenCommitted.Add(sdk.NewCoin(ptypes.EdenB, newUnclaimedEdenBoostFromEdenCommited))

			// Update Commitments with new unclaimed token amounts
			k.UpdateCommitments(ctx, creator, &commitments, newUnclaimedEdenTokens, newUnclaimedEdenBoostTokens, dexRewards)

			// Update sub buckets commitment with new unclaimed token amounts
			k.UpdateCommitmentsSubBuckets(ctx, creator, &commitments, rewardsByElysStaking, rewardsByEdenCommitted, rewardsByEdenBCommitted, rewardsByUSDCDeposit)
			return false
		},
	)

	// Calcualte the remainings
	edenRemained := edenAmountPerEpochStakers.Sub(totalEdenGiven)
	edenRemainedLP := edenAmountPerEpochLPs.Sub(totalEdenGivenLP)
	dexRewardsRemained := dexRevenueStakersAmt.Sub(sdk.NewDecFromInt(totalRewardsGiven))
	dexRewardsRemainedLP := dexRevenueLPsAmt.Add(gasFeesLPsAmt).Sub(sdk.NewDecFromInt(totalRewardsGivenLP))

	// Fund community the remain coins
	// ----------------------------------
	edenRemainedCoin := sdk.NewDecCoin(ptypes.Eden, edenRemained.Add(edenRemainedLP))
	// TODO:
	// Dummy denom for USDC
	dexRewardsRemainedCoin := sdk.NewDecCoinFromDec(ptypes.BaseCurrency, dexRewardsRemained.Add(dexRewardsRemainedLP))

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(edenRemainedCoin)
	feePool.CommunityPool = feePool.CommunityPool.Add(dexRewardsRemainedCoin)
	k.SetFeePool(ctx, feePool)
	// ----------------------------------
}

// Update commitment record
func (k Keeper) UpdateCommitments(ctx sdk.Context, creator string, commitments *ctypes.Commitments, newUnclaimedEdenTokens sdk.Int, newUnclaimedEdenBoostTokens sdk.Int, dexRewards sdk.Int) {
	// Update unclaimed Eden balances in the Commitments structure
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))
	// Update unclaimed Eden-Boost token balances in the Commitments structure
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.EdenB, newUnclaimedEdenBoostTokens))

	// All dex revenue are collected to incentive module in USDC
	// Gas fees (Elys) are also converted into USDC and collected into total dex revenue wallet of incentive module.
	// Update USDC balances in the Commitments structure.
	// TODO:
	// USDC token denom is dummy one for now until we have real usdc brought through bridge.
	// These are the rewards from each pool, margin, gas fee.
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.BaseCurrency, dexRewards))

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Update sub bucket commitment record
func (k Keeper) UpdateCommitmentsSubBuckets(ctx sdk.Context, creator string, commitments *ctypes.Commitments, rewardsByElysStaking sdk.Coins, rewardsByEdenCommitted sdk.Coins, rewardsByEdenBCommitted sdk.Coins, rewardsByUSDCDeposit sdk.Coins) {
	// Add to Elys staking bucket
	commitments.AddSubBucketRewardsByElysUnclaimed(rewardsByElysStaking)
	// Add to Eden committed bucket
	commitments.AddSubBucketRewardsByEdenUnclaimed(rewardsByEdenCommitted)
	// Add to EdenB committed bucket
	commitments.AddSubBucketRewardsByEdenBUnclaimed(rewardsByEdenBCommitted)
	// Add to USDC deposit bucket
	commitments.AddSubBucketRewardsByUsdcUnclaimed(rewardsByUSDCDeposit)

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Calculate Proxy TVL
func (k Keeper) CalculateProxyTVL(ctx sdk.Context) sdk.Dec {
	multipliedShareSum := sdk.ZeroDec()
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, p.GetPoolId())
		if !found {
			return false
		}

		proxyTVL := tvl.Mul(poolInfo.Multiplier)

		// Calculate total pool share by TVL and multiplier
		multipliedShareSum = multipliedShareSum.Add(proxyTVL)

		return false
	})

	//-----------------------------------
	// Handle stable stake pool
	stableStakePoolId := uint64(stabletypes.PoolId)

	// Get pool info from incentive param
	poolInfo, found := k.GetPoolInfo(ctx, stableStakePoolId)
	if !found {
		k.InitPoolMultiplier(ctx, stableStakePoolId)
		poolInfo, _ = k.GetPoolInfo(ctx, stableStakePoolId)
	}
	tvl := stabletypes.TVL(ctx, k.authKeeper, k.bankKeeper)
	proxyTVL := sdk.NewDecFromInt(tvl).Mul(poolInfo.Multiplier)
	multipliedShareSum = multipliedShareSum.Add(proxyTVL)

	// return total sum of TVL share using multiplier of all pools
	return multipliedShareSum
}

// Caculate total TVL
func (k Keeper) CalculateTVL(ctx sdk.Context) sdk.Dec {
	TVL := sdk.ZeroDec()

	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}
		TVL = TVL.Add(tvl)
		return false
	})

	return TVL
}
