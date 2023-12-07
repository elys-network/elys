package keeper

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

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

<<<<<<< HEAD
// Calculate the amm pool share using pool multiplier
func (k Keeper) CalculatePoolShare(ctx sdk.Context, pool *ammtypes.Pool) (sdk.Dec, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroDec(), sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Calculate total proxy tvl
	totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

	// Calculate the tvl of the pool
	tvl, err := pool.TVL(ctx, k.oracleKeeper)
	if err != nil {
		return sdk.ZeroDec(), errors.New("pool tvl is invalid")
=======
// Calculate the amm pool ratio
func (k Keeper) CalculatePoolRatio(ctx sdk.Context, pool *ammtypes.Pool) string {
	weightString := ""
	totalWeight := sdk.ZeroInt()
	for _, asset := range pool.PoolAssets {
		totalWeight = totalWeight.Add(asset.Weight)
>>>>>>> 924d0fa (chore: add leverage type and fix pool ratio)
	}

	if totalWeight.IsZero() {
		return weightString
	}

	for _, asset := range pool.PoolAssets {
		weight := sdk.NewDecFromInt(asset.Weight).QuoInt(totalWeight).MulInt(sdk.NewInt(100)).TruncateInt64()
		weightString = fmt.Sprintf("%s : %d", weightString, weight)
	}

	// remove prefix " : " 3 characters
	if len(weightString) > 1 {
		weightString = weightString[3:]
	}

	// returns pool weight string
	return weightString
}
