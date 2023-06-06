package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

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
