package keeper

import (
	gomath "math"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Calculate the delegated amount
func (k Keeper) CalcDelegationAmount(ctx sdk.Context, delegator string) math.Int {
	// Derivate bech32 based delegator address
	delAddr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		// This could be validator address
		return sdk.ZeroInt()
	}

	// Get elys delegation
	delAmount := sdk.ZeroDec()
	delegations := k.Keeper.GetDelegatorDelegations(ctx, delAddr, gomath.MaxUint16)
	for _, del := range delegations {
		valAddr := del.GetValidatorAddr()
		val := k.Keeper.Validator(ctx, valAddr)

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
	delegations := k.Keeper.GetDelegatorDelegations(ctx, delAddr, gomath.MaxUint16)
	for _, del := range delegations {
		valAddr := del.GetValidatorAddr()
		val := k.Keeper.Validator(ctx, valAddr)

		if !val.IsBonded() {
			continue
		}
		shares := del.GetShares()
		tokens := val.TokensFromSharesTruncated(shares)
		delAmount = delAmount.Add(tokens)
	}

	return delAmount.TruncateInt()
}
