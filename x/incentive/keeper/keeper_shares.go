package keeper

import (
	"fmt"
	gomath "math"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// Calculate the delegated amount
func (k Keeper) CalcDelegationAmount(ctx sdk.Context, delegator string) math.Int {
	// Derivate bech32 based delegator address
	delAddr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		// This could be validator address
		return sdk.ZeroInt()
	}

	// Get elys delegation for creator address
	delAmount := sdk.ZeroDec()

	// Get all delegations
	delegations := k.stk.GetDelegatorDelegations(ctx, delAddr, gomath.MaxUint16)
	for _, del := range delegations {
		// Get validator address
		valAddr := del.GetValidatorAddr()
		// Get validator
		val := k.stk.Validator(ctx, valAddr)

		shares := del.GetShares()
		tokens := val.TokensFromSharesTruncated(shares)
		delAmount = delAmount.Add(tokens)
	}

	return delAmount.TruncateInt()
}

// Calculate delegation to bonded validators
func (k Keeper) CalcBondedDelegationAmount(ctx sdk.Context, delegator string) math.Int {
	// Derivate bech32 based delegator address
	delAddr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		// This could be validator address
		return sdk.ZeroInt()
	}

	// Get elys delegation for creator address
	delAmount := sdk.ZeroDec()

	// Get all delegations
	delegations := k.stk.GetDelegatorDelegations(ctx, delAddr, gomath.MaxUint16)
	for _, del := range delegations {
		// Get validator address
		valAddr := del.GetValidatorAddr()
		// Get validator
		val := k.stk.Validator(ctx, valAddr)

		if !val.IsBonded() {
			continue
		}
		shares := del.GetShares()
		tokens := val.TokensFromSharesTruncated(shares)
		delAmount = delAmount.Add(tokens)
	}

	return delAmount.TruncateInt()
}

// Calculate the amm pool ratio
func (k Keeper) CalculatePoolRatio(ctx sdk.Context, pool *ammtypes.Pool) string {
	weightString := ""
	totalWeight := sdk.ZeroInt()
	for _, asset := range pool.PoolAssets {
		totalWeight = totalWeight.Add(asset.Weight)
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
