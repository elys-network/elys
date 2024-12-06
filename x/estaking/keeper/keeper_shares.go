package keeper

import (
	gomath "math"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Calculate the delegated amount
func (k Keeper) CalcDelegationAmount(ctx sdk.Context, delegator sdk.AccAddress) math.Int {
	// Get elys delegation
	delAmount := math.LegacyZeroDec()
	delegations, err := k.Keeper.GetDelegatorDelegations(ctx, delegator, gomath.MaxUint16)
	if err != nil {
		panic(err)
	}
	for _, del := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(del.GetValidatorAddr())
		if err != nil {
			panic(err)
		}
		val, err := k.Keeper.Validator(ctx, valAddr)
		if err != nil {
			panic(err)
		}

		shares := del.GetShares()
		tokens := val.TokensFromSharesTruncated(shares)
		delAmount = delAmount.Add(tokens)
	}

	return delAmount.TruncateInt()
}
