package keeper

import (
	"fmt"
	"math"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	etypes "github.com/elys-network/elys/x/epochs/types"
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

// Move gas fees collected to incentive module
func (k Keeper) CollectGasFeesToIncentiveModule(ctx sdk.Context) sdk.Coins {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	// feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	return feesCollectedInt
}

// Pull DEX revenus collected to incentive module
func (k Keeper) CollectDEXRevenusToIncentiveModule(ctx sdk.Context) sdk.Coins {
	// transfer collected fees to the distribution module account

	return sdk.Coins{}
}

// Fund community pool based on community tax
func (k Keeper) UpdateCommunityPool(ctx sdk.Context, amt sdk.DecCoins) sdk.DecCoins {
	// calculate fraction allocated to validators
	communityTax := k.GetCommunityTax(ctx)
	communityRevenus := amt.MulDecTruncate(communityTax)

	// allocate community funding
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(communityRevenus...)
	k.SetFeePool(ctx, feePool)

	return amt.Sub(communityRevenus)
}

// Update total commitment info
func (k Keeper) UpdateTotalCommitmentInfo(ctx sdk.Context) {
	// Fetch total staked Elys amount again
	k.tci.TotalElysBonded = k.stk.TotalBondedTokens(ctx)
	// Initialize with amount zero
	k.tci.TotalEdenEdenBoostCommitted = sdk.ZeroInt()
	// Initialize with amount zero
	k.tci.TotalFeesCollected = sdk.Coins{}
	// Initialize Lp tokens amount
	k.tci.TotalLpTokensCommitted = make(map[string]sdk.Int)

	// Collect gas fees collected
	fees := k.CollectGasFeesToIncentiveModule(ctx)
	// Collect DEX revenus collected
	dexFees := k.CollectDEXRevenusToIncentiveModule(ctx)
	// Calculate total fees - DEX revenus + Gas fees collected
	k.tci.TotalFeesCollected = k.tci.TotalFeesCollected.Add(dexFees...).Add(fees...)

	// Iterate to calculate total Eden, Eden boost and Lp tokens committed
	k.cmk.IterateCommitments(ctx, func(commitments ctypes.Commitments) bool {
		committedEdenToken := commitments.GetCommittedAmountForDenom(ptypes.Eden)
		committedEdenBoostToken := commitments.GetCommittedAmountForDenom(ptypes.EdenB)

		k.tci.TotalEdenEdenBoostCommitted = k.tci.TotalEdenEdenBoostCommitted.Add(committedEdenToken).Add(committedEdenBoostToken)

		// Iterate to calcaulte total Lp tokens committed
		k.lpk.IterateLiquidityPools(ctx, func(l LiquidityPool) bool {
			committedLpToken := commitments.GetCommittedAmountForDenom(l.lpToken)
			k.tci.TotalLpTokensCommitted[l.lpToken] = k.tci.TotalLpTokensCommitted[l.lpToken].Add(committedLpToken)
			return false
		})
		return false
	})
}

// Calculate total share of staking
func (k Keeper) CalculateTotalShareOfStaking(amount sdk.Int) sdk.Dec {
	// Total statked = Elys staked + Eden Committed + Eden boost Committed
	totalStaked := k.tci.TotalElysBonded.Add(k.tci.TotalEdenEdenBoostCommitted)
	if totalStaked.LTE(sdk.ZeroInt()) {
		return sdk.ZeroDec()
	}

	// Share = Amount / Total Staked
	return sdk.NewDecFromInt(amount).QuoInt(totalStaked)
}

// Calculate the delegated amount
func (k Keeper) CalculateDelegatedAmount(ctx sdk.Context, delegator string) sdk.Int {
	// Derivate bech32 based delegator address
	delAdr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		// This could be validator address
		return sdk.ZeroInt()
	}

	// Get elys delegation for creator address
	delegatedAmt := sdk.ZeroDec()

	// Get all delegations
	delegations := k.stk.GetDelegatorDelegations(ctx, delAdr, math.MaxUint16)
	for _, del := range delegations {
		// Get validator address
		valAddr := del.GetValidatorAddr()
		// Get validator
		val := k.stk.Validator(ctx, valAddr)

		shares := del.GetShares()
		tokens := val.TokensFromSharesTruncated(shares)
		delegatedAmt = delegatedAmt.Add(tokens)
	}

	return delegatedAmt.TruncateInt()
}

// Find out active incentive params
func (k Keeper) GetProperIncentiveParam(ctx sdk.Context, epochIdentifier string) (bool, types.IncentiveInfo, types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// If we don't have enough params
	if len(params.StakeIncentives) < 1 || len(params.LpIncentives) < 1 {
		return false, types.IncentiveInfo{}, types.IncentiveInfo{}
	}

	// Current block timestamp
	timestamp := ctx.BlockTime().Unix()
	foundIncentive := false

	// Incentive params initialize
	stakeIncentive := params.StakeIncentives[0]
	lpIncentive := params.LpIncentives[0]

	// Consider epochIdentifier and start time
	// Consider epochNumber as well
	if stakeIncentive.EpochIdentifier != epochIdentifier || timestamp < stakeIncentive.StartTime.Unix() {
		return false, types.IncentiveInfo{}, types.IncentiveInfo{}
	}

	// Increase current epoch of Stake incentive param
	stakeIncentive.CurrentEpoch = stakeIncentive.CurrentEpoch + 1
	if stakeIncentive.CurrentEpoch == stakeIncentive.NumEpochs {
		params.StakeIncentives = params.StakeIncentives[1:]
	}

	// Increase current epoch of Lp incentive param
	lpIncentive.CurrentEpoch = lpIncentive.CurrentEpoch + 1
	if lpIncentive.CurrentEpoch == lpIncentive.NumEpochs {
		params.LpIncentives = params.LpIncentives[1:]
	}

	// Update params
	k.SetParams(ctx, params)

	// return found, stake, lp incentive params
	return foundIncentive, stakeIncentive, lpIncentive
}

// Update uncommitted token amount
// Called back through epoch hook
func (k Keeper) UpdateUncommittedTokens(ctx sdk.Context, epochIdentifier string, stakeIncentive types.IncentiveInfo, lpIncentive types.IncentiveInfo) {
	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx)

	// Calculate 65% for LP, 35% for Stakers
	dexRevenue := sdk.NewDecCoinsFromCoins(k.tci.TotalFeesCollected...)
	devRevenue65 := dexRevenue.MulDecTruncate(sdk.NewDecWithPrec(65, 1))
	devRevenue35 := dexRevenue.Sub(devRevenue65)

	// Fund community pool based on the communtiy tax
	devRevenueRemained35 := k.UpdateCommunityPool(ctx, devRevenue35)

	// Elys amount in sdk.Dec type
	devRevenue65Amt := devRevenue65.AmountOf(ptypes.Elys)
	devRevenue35Amt := devRevenueRemained35.AmountOf(ptypes.Elys)

	// Calculate eden amount per epoch
	edenAmountPerEpochStake := stakeIncentive.Amount.Quo(sdk.NewInt(stakeIncentive.NumEpochs))
	edenAmountPerEpochLp := lpIncentive.Amount.Quo(sdk.NewInt(lpIncentive.NumEpochs))
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
			newUncommittedEdenTokens, dexRewards, dexRewardsByStake := k.CalculateNewUncommittedEdenTokensAndDexRewards(ctx, delegatedAmt, commitments, edenAmountPerEpochStake, devRevenue35Amt)
			totalEdenGiven = totalEdenGiven.Add(newUncommittedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Calculate new uncommitted Eden tokens from LpTokens committed, Dex rewards distribution
			newUncommittedEdenTokensLp, dexRewardsLp := k.CalculateNewUncommittedEdenTokensFromLPAndDexRewards(ctx, totalProxyTVL, commitments, edenAmountPerEpochLp, devRevenue65Amt)
			totalEdenGivenLP = totalEdenGivenLP.Add(newUncommittedEdenTokensLp)
			totalRewardsGivenLP = totalRewardsGivenLP.Add(dexRewardsLp)

			// Calculate the total Eden uncommitted amount
			newUncommittedEdenTokens = newUncommittedEdenTokens.Add(newUncommittedEdenTokensLp)

			// Give commission to validators ( Eden from stakers and Dex rewards from stakers. )
			edenCommissionGiven, dexRewardsCommissionGiven := k.GiveCommissionToValidators(ctx, creator, delegatedAmt, newUncommittedEdenTokens, dexRewardsByStake)

			// Minus the given amount and increase with the remains only
			newUncommittedEdenTokens = newUncommittedEdenTokens.Sub(edenCommissionGiven)

			// Plus LpDexRewards and minus commission given
			dexRewards = dexRewards.Add(dexRewardsLp).Sub(dexRewardsCommissionGiven)

			// Calculate new uncommitted Eden-Boost tokens for staker and Eden token holders
			newUncommittedEdenBoostTokens := k.CalculateNewUncommittedEdenBoostTokens(ctx, delegatedAmt, commitments, epochIdentifier, edenBoostAPR)

			// Update Commitments with new uncommitted token amounts
			k.UpdateCommitments(ctx, creator, &commitments, newUncommittedEdenTokens, newUncommittedEdenBoostTokens, dexRewards)

			return false
		},
	)

	// Calcualte the remainings
	edenRemained := edenAmountPerEpochStake.Sub(totalEdenGiven)
	edenRemainedLP := edenAmountPerEpochLp.Sub(totalEdenGivenLP)
	dexRewardsRemained := devRevenue35Amt.Sub(sdk.NewDecFromInt(totalRewardsGiven))
	dexRewardsRemainedLP := devRevenue65Amt.Sub(sdk.NewDecFromInt(totalRewardsGivenLP))

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

// Calculate new Eden token amounts based on the given conditions and user's current uncommitted token balance
func (k Keeper) CalculateNewUncommittedEdenTokensAndDexRewards(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, edenAmountPerEpoch sdk.Int, devRevenue35Amt sdk.Dec) (sdk.Int, sdk.Int, sdk.Dec) {
	// -----------Eden & Eden boost calculation ---------------------
	// --------------------------------------------------------------
	// Get eden commitments and eden boost commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)
	edenBoostCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)

	// compute eden reward based on above and param factors for each
	totalEdenCommittedByStake := delegatedAmt.Add(edenCommitted).Add(edenBoostCommitted)
	stakeShare := k.CalculateTotalShareOfStaking(totalEdenCommittedByStake)

	// Calculate newly creating eden amount by its share
	newEdenAllocated := stakeShare.MulInt(edenAmountPerEpoch)

	// -----------------Fund community Eden token----------------------
	// ----------------------------------------------------------------
	edenCoin := sdk.NewDecCoinFromDec(ptypes.Eden, newEdenAllocated)
	newEdenCoinRemained := k.UpdateCommunityPool(ctx, sdk.DecCoins{edenCoin})

	// Get remained Eden amount
	newEdenAllocated = newEdenCoinRemained.AmountOf(ptypes.Eden)

	// --------------------DEX rewards calculation --------------------
	// ----------------------------------------------------------------
	// Calculate dex rewards
	dexRewards := stakeShare.Mul(devRevenue35Amt).TruncateInt()

	// Calculate only elys staking share
	stakeShareByStakeOnly := k.CalculateTotalShareOfStaking(delegatedAmt)
	dexRewardsByStakeOnly := stakeShareByStakeOnly.Mul(devRevenue35Amt)

	return newEdenAllocated.TruncateInt(), dexRewards, dexRewardsByStakeOnly
}

// Calculate new Eden token amounts based on LpElys committed and MElys committed
func (k Keeper) CalculateNewUncommittedEdenTokensFromLPAndDexRewards(ctx sdk.Context, totalProxyTVL sdk.Dec, commitments ctypes.Commitments, edenAmountPerEpochLp sdk.Int, devRevenue65Amt sdk.Dec) (sdk.Int, sdk.Int) {
	// Method 2 - Using Proxy TVL
	totalNewEdenAllocated := sdk.ZeroInt()
	totalDexRewardsAllocated := sdk.ZeroDec()

	// Iterate to calculate total Eden from LpElys, MElys committed
	k.lpk.IterateLiquidityPools(ctx, func(l LiquidityPool) bool {
		// ------------ New Eden calculation -------------------
		// -----------------------------------------------------
		// newEdenAllocated = 80 / ( 80 + 90 + 200 + 0) * 100
		// Pool share = 80
		// edenAmountPerEpochLp = 100

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := sdk.NewDecFromInt(l.TVL).MulInt64(l.multiplier)
		poolShare := proxyTVL.Quo(totalProxyTVL)

		// Calculate new Eden for this pool
		newEdenAllocatedForPool := poolShare.MulInt(edenAmountPerEpochLp)

		// this lp token committed
		commmittedLpToken := commitments.GetCommittedAmountForDenom(l.lpToken)
		// this lp token total committed
		totalCommittedLpToken := k.tci.TotalLpTokensCommitted[l.lpToken]

		// Calculalte lp token share of the pool
		lpShare := sdk.NewDecFromInt(commmittedLpToken).QuoInt(totalCommittedLpToken)

		// Calculate new Eden allocated per LP
		newEdenAllocated := lpShare.Mul(newEdenAllocatedForPool).TruncateInt()

		// Sum the total amount
		totalNewEdenAllocated = totalNewEdenAllocated.Add(newEdenAllocated)
		// -------------------------------------------------------

		// ------------------- DEX rewards calculation -------------------
		// ---------------------------------------------------------------
		// Calculate dex rewards per pool
		dexRewardsAllocatedForPool := poolShare.Mul(devRevenue65Amt)
		// Calculate dex rewards per lp
		dexRewardsForLP := lpShare.Mul(dexRewardsAllocatedForPool)
		// Sum total rewards per commitment
		totalDexRewardsAllocated = totalDexRewardsAllocated.Add(dexRewardsForLP)
		//----------------------------------------------------------------
		return false
	})

	// return
	return totalNewEdenAllocated, totalDexRewardsAllocated.TruncateInt()
}

// Calculate epoch counts per year to be used in APR calculation
func (k Keeper) CalculateEpochCountsPerYear(epochIdentifier string) int64 {
	switch epochIdentifier {
	case etypes.WeekEpochID:
		return ptypes.WeeksPerYear
	case etypes.DayEpochID:
		return ptypes.DaysPerYear
	case etypes.HourEpochID:
		return ptypes.HoursPerYear
	}

	return 0
}

// Calculate new Eden-Boost token amounts based on the given conditions and user's current uncommitted token balance
func (k Keeper) CalculateNewUncommittedEdenBoostTokens(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, epochIdentifier string, edenBoostAPR int64) sdk.Int {
	// Get eden commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Compute eden reward based on above and param factors for each
	totalEden := delegatedAmt.Add(edenCommitted)

	// Calculate edenBoostAPR % APR for eden boost
	epochNumsPerYear := k.CalculateEpochCountsPerYear(epochIdentifier)

	return totalEden.Quo(sdk.NewInt(epochNumsPerYear)).Quo(sdk.NewInt(100)).Mul(sdk.NewInt(edenBoostAPR))
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

// Increase uncommitted token amount for the corresponding validator
func (k Keeper) UpdateTokensForValidator(ctx sdk.Context, validator string, new_uncommitted_eden_tokens sdk.Int, dexRewards sdk.Dec) {
	commitments, bfound := k.cmk.GetCommitments(ctx, validator)
	if !bfound {
		return
	}

	// Update Eden amount
	k.UpdateTokensCommitment(&commitments, new_uncommitted_eden_tokens, ptypes.Eden)

	// Update Elys amount
	k.UpdateTokensCommitment(&commitments, dexRewards.TruncateInt(), ptypes.Elys)

	// Update commmitment
	k.cmk.SetCommitments(ctx, commitments)
}

// Give commissions to validators
func (k Keeper) GiveCommissionToValidators(ctx sdk.Context, delegator string, totalDelegationAmt sdk.Int, newUncommittedAmt sdk.Int, dexRewards sdk.Dec) (sdk.Int, sdk.Int) {
	delAdr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt()
	}

	// If there is no delegation, (not elys staker)
	if totalDelegationAmt.LTE(sdk.ZeroInt()) {
		return sdk.ZeroInt(), sdk.ZeroInt()
	}

	// Total Eden given
	totalEdenGiven := sdk.ZeroInt()
	totalDexRewardsGiven := sdk.ZeroInt()

	// Iterate all delegated validators
	k.stk.IterateDelegations(ctx, delAdr, func(index int64, del stypes.DelegationI) (stop bool) {
		valAddr := del.GetValidatorAddr()
		// Get validator
		val := k.stk.Validator(ctx, valAddr)
		// Get commission rate
		comm_rate := val.GetCommission()
		// Get delegator share
		shares := del.GetShares()
		// Get token amount delegated
		delegatedAmt := val.TokensFromSharesTruncated(shares)

		//-----------------------------
		// Eden commission
		//-----------------------------
		// to give = delegated amount / total delegation * newly minted eden * commission rate
		edenCommission := delegatedAmt.QuoInt(totalDelegationAmt).MulInt(newUncommittedAmt).Mul(comm_rate)
		// Sum total commission given
		totalEdenGiven = totalEdenGiven.Add(edenCommission.TruncateInt())
		//-----------------------------

		//-----------------------------
		// Dex rewards commission
		//-----------------------------
		// to give = delegated amount / total delegation * newly minted eden * commission rate
		dexRewardsCommission := delegatedAmt.QuoInt(totalDelegationAmt).Mul(dexRewards).Mul(comm_rate)
		// Sum total commission given
		totalDexRewardsGiven = totalDexRewardsGiven.Add(dexRewardsCommission.TruncateInt())
		//-----------------------------

		// increase uncomitted token amount of validator's commitment
		k.UpdateTokensForValidator(ctx, valAddr.String(), edenCommission.TruncateInt(), dexRewardsCommission)

		return false
	})

	return totalEdenGiven, totalDexRewardsGiven
}

// withdraw rewards
// Eden, EdenBoost and Elys to USDC
func (k Keeper) ProcessWithdrawRewards(ctx sdk.Context, delegator string) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return err
	}

	// Get commitments
	commitments, bfound := k.cmk.GetCommitments(ctx, delegator)
	if !bfound {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to find commitment")
	}

	// Eden
	uncommittedEden, _ := commitments.GetUncommittedTokensForDenom(ptypes.Eden)
	// Eden B
	uncommittedEdenB, _ := commitments.GetUncommittedTokensForDenom(ptypes.EdenB)
	// Elys
	uncommittedElys, _ := commitments.GetUncommittedTokensForDenom(ptypes.Elys)

	// Withdraw Eden
	err = k.cmk.ProcessWithdrawTokens(ctx, delegator, ptypes.Eden, uncommittedEden.Amount)

	// Withdraw Eden Boost
	err = k.cmk.ProcessWithdrawTokens(ctx, delegator, ptypes.EdenB, uncommittedEdenB.Amount)

	// Convert Elys to USDC
	// Conversion is done inside commmitment module, later on we will have swap module
	err = k.cmk.ProcessWithdrawElysTokens(ctx, delegator, ptypes.Elys, uncommittedElys.Amount)

	return err
}

// withdraw validator commission
// Eden, EdenBoost and Elys to USDC
func (k Keeper) ProcessWithdrawValidatorCommission(ctx sdk.Context, delegator string, validator string) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return err
	}

	_, err = sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// Get commitments
	commitments, bfound := k.cmk.GetCommitments(ctx, validator)
	if !bfound {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to find commitment")
	}

	// Eden
	uncommittedEden, _ := commitments.GetUncommittedTokensForDenom(ptypes.Eden)
	// Eden B
	uncommittedEdenB, _ := commitments.GetUncommittedTokensForDenom(ptypes.EdenB)
	// Elys
	uncommittedElys, _ := commitments.GetUncommittedTokensForDenom(ptypes.Elys)

	// Withdraw Eden
	err = k.cmk.ProcessWithdrawValidatorCommission(ctx, delegator, validator, ptypes.Eden, uncommittedEden.Amount)

	// Withdraw Eden Boost
	err = k.cmk.ProcessWithdrawValidatorCommission(ctx, delegator, validator, ptypes.EdenB, uncommittedEdenB.Amount)

	// Convert Elys to USDC
	// Conversion is done inside commmitment module, later on we will have swap module
	err = k.cmk.ProcessWithdrawValidatorElysCommission(ctx, delegator, validator, ptypes.Elys, uncommittedElys.Amount)

	return err
}
