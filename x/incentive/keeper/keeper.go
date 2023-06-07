package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		cmk        types.CommitmentKeeper
		stk        types.StakingKeeper
		tci        *types.TotalCommitmentInfo
		authKeeper types.AccountKeeper
		bankKeeper types.BankKeeper

		feeCollectorName    string // name of the FeeCollector ModuleAccount
		dexRevCollectorName string // name of the Dex Revenue ModuleAccount

		lpk *LiquidityKeeper
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
		lpk:                 NewLiquidityKeeper(),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Update uncommitted token amount
// Called back through epoch hook
func (k Keeper) UpdateUncommittedTokens(ctx sdk.Context, epochIdentifier string, stakeIncentive types.IncentiveInfo, lpIncentive types.IncentiveInfo) {
	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx)

	// Calculate 65% for LP, 35% for Stakers
	// Collect DEX revenue collected
	dexRevenue := k.CollectDEXRevenusToIncentiveModule(ctx)

	dexRevenueDec := sdk.NewDecCoinsFromCoins(dexRevenue...)
	rewardPercentForLps := k.GetDEXRewardPercentForLPs(ctx)
	dexRevenueForLps := dexRevenueDec.MulDecTruncate(rewardPercentForLps)
	dexRevenueForStakers := dexRevenueDec.Sub(dexRevenueForLps)

	// Fund community pool based on the communtiy tax
	dexRevenueRemainedForStakers := k.UpdateCommunityPool(ctx, dexRevenueForStakers)

	// Elys amount in sdk.Dec type
	dexRevenueLPsAmt := dexRevenueForLps.AmountOf(ptypes.Elys)
	dexRevenueStakersAmt := dexRevenueRemainedForStakers.AmountOf(ptypes.Elys)

	// Calculate eden amount per epoch
	edenAmountPerEpochStakers := stakeIncentive.Amount.Quo(sdk.NewInt(stakeIncentive.NumEpochs))
	edenAmountPerEpochLPs := lpIncentive.Amount.Quo(sdk.NewInt(lpIncentive.NumEpochs))
	edenBoostAPR := stakeIncentive.EdenBoostApr

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have mulitplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL := k.lpk.CalculateProxyTVL()

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

			// Calculate delegated amount per delegator
			delegatedAmt := k.CalculateDelegatedAmount(ctx, creator)

			// Calculate new uncommitted Eden tokens from Eden & Eden boost committed, Dex rewards distribution
			newUncommittedEdenTokens, dexRewards, dexRewardsByStakers := k.CalculateRewardsForStakers(ctx, delegatedAmt, commitments, edenAmountPerEpochStakers, dexRevenueStakersAmt)
			totalEdenGiven = totalEdenGiven.Add(newUncommittedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Calculate new uncommitted Eden tokens from LpTokens committed, Dex rewards distribution
			newUncommittedEdenTokensLp, dexRewardsLp := k.CalculateRewardsForLPs(ctx, totalProxyTVL, commitments, edenAmountPerEpochLPs)
			totalEdenGivenLP = totalEdenGivenLP.Add(newUncommittedEdenTokensLp)
			totalRewardsGivenLP = totalRewardsGivenLP.Add(dexRewardsLp)

			// Calculate the total Eden uncommitted amount
			newUncommittedEdenTokens = newUncommittedEdenTokens.Add(newUncommittedEdenTokensLp)

			// Give commission to validators ( Eden from stakers and Dex rewards from stakers. )
			edenCommissionGiven, dexRewardsCommissionGiven := k.GiveCommissionToValidators(ctx, creator, delegatedAmt, newUncommittedEdenTokens, dexRewardsByStakers)

			// Minus the given amount and increase with the remains only
			newUncommittedEdenTokens = newUncommittedEdenTokens.Sub(edenCommissionGiven)

			// Plus LpDexRewards and minus commission given
			dexRewards = dexRewards.Add(dexRewardsLp).Sub(dexRewardsCommissionGiven)

			// Calculate new uncommitted Eden-Boost tokens for staker and Eden token holders
			newUncommittedEdenBoostTokens := k.CalculateEdenBoostRewards(ctx, delegatedAmt, commitments, epochIdentifier, edenBoostAPR)

			// Update Commitments with new uncommitted token amounts
			k.UpdateCommitments(ctx, creator, &commitments, newUncommittedEdenTokens, newUncommittedEdenBoostTokens, dexRewards)

			return false
		},
	)

	// After give DEX rewards, we should update its record in order to avoid double spend.
	k.lpk.UpdateRewardsAccmulated(ctx)

	// Calcualte the remainings
	edenRemained := edenAmountPerEpochStakers.Sub(totalEdenGiven)
	edenRemainedLP := edenAmountPerEpochLPs.Sub(totalEdenGivenLP)
	dexRewardsRemained := dexRevenueStakersAmt.Sub(sdk.NewDecFromInt(totalRewardsGiven))
	dexRewardsRemainedLP := dexRevenueLPsAmt.Sub(sdk.NewDecFromInt(totalRewardsGivenLP))

	// Fund community the remain coins
	// ----------------------------------
	edenRemainedCoin := sdk.NewDecCoin(ptypes.Eden, edenRemained.Add(edenRemainedLP))
	dexRewardsRemainedCoin := sdk.NewDecCoinFromDec(ptypes.Elys, dexRewardsRemained.Add(dexRewardsRemainedLP))

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(edenRemainedCoin)
	feePool.CommunityPool = feePool.CommunityPool.Add(dexRewardsRemainedCoin)
	k.SetFeePool(ctx, feePool)
	// ----------------------------------
}

func (k Keeper) UpdateCommitments(ctx sdk.Context, creator string, commitments *ctypes.Commitments, newUncommittedEdenTokens sdk.Int, newUncommittedEdenBoostTokens sdk.Int, dexRewards sdk.Int) {
	// Update uncommitted Eden balances in the Commitments structure
	k.UpdateTokensCommitment(commitments, newUncommittedEdenTokens, ptypes.Eden)
	// Update uncommitted Eden-Boost token balances in the Commitments structure
	k.UpdateTokensCommitment(commitments, newUncommittedEdenBoostTokens, ptypes.EdenB)
	// Update Elys balances in the Commitments structure
	k.UpdateTokensCommitment(commitments, dexRewards, ptypes.Elys)

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Update the uncommitted Eden token balance
func (k Keeper) UpdateTokensCommitment(commitments *ctypes.Commitments, new_uncommitted_eden_tokens sdk.Int, denom string) {
	uncommittedEden, found := commitments.GetUncommittedTokensForDenom(denom)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedTokens = append(uncommittedTokens, &ctypes.UncommittedTokens{
			Denom:  denom,
			Amount: new_uncommitted_eden_tokens,
		})
		commitments.UncommittedTokens = uncommittedTokens
	} else {
		uncommittedEden.Amount = uncommittedEden.Amount.Add(new_uncommitted_eden_tokens)
	}
}
