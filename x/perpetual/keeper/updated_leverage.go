package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdatedLeverage(ctx sdk.Context, mtp types.MTP) (sdk.Dec, error) {
	var custody_in_usdc sdk.Dec
	if mtp.CustodyAsset != ptypes.BaseCurrency {
		price := k.amm.EstimatePrice(ctx, mtp.CustodyAsset, ptypes.BaseCurrency)
		custody_in_usdc = math.LegacyDec(mtp.Custody).Mul(price)
	}else {
		custody_in_usdc = math.LegacyDec(mtp.Custody)
	}
	var denominator sdk.Dec
	if mtp.LiabilitiesAsset != ptypes.BaseCurrency {
		price := k.amm.EstimatePrice(ctx, mtp.CustodyAsset, ptypes.BaseCurrency)
		denominator = custody_in_usdc.Sub(math.LegacyDec(mtp.Liabilities).Mul(price))
	}else {
		denominator = custody_in_usdc.Sub(math.LegacyDec(mtp.Liabilities))
	}
	if denominator.IsZero() {
		return sdk.ZeroDec(),  nil
	}
	updated_leverage := custody_in_usdc.Quo(denominator)

	return updated_leverage, nil
}
