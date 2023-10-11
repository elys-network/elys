package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) CalcMTPInterestLiabilities(ctx sdk.Context, mtp *types.MTP, interestRate sdk.Dec, epochPosition, epochLength int64, ammPool ammtypes.Pool, collateralAsset string) sdk.Int {
	var interestRational, liabilitiesRational, rate, epochPositionRational, epochLengthRational big.Rat

	rate.SetFloat64(interestRate.MustFloat64())

	collateralIndex, _ := k.GetMTPAssetIndex(mtp, collateralAsset, "")
	unpaidCollaterals := sdk.ZeroInt()
	// Calculate collateral interests in base currency
	if mtp.Collaterals[collateralIndex].Denom == ptypes.BaseCurrency {
		unpaidCollaterals = unpaidCollaterals.Add(mtp.InterestUnpaidCollaterals[collateralIndex])
	} else {
		// Liability is in base currency, so convert it to base currency
		unpaidCollateralIn := sdk.NewCoin(mtp.Collaterals[collateralIndex].Denom, mtp.InterestUnpaidCollaterals[collateralIndex])
		C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, ptypes.BaseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroInt()
		}

		unpaidCollaterals = unpaidCollaterals.Add(C)
	}

	liabilitiesRational.SetInt(mtp.Liabilities.BigInt().Add(mtp.Liabilities.BigInt(), unpaidCollaterals.BigInt()))
	interestRational.Mul(&rate, &liabilitiesRational)

	if epochPosition > 0 { // prorate interest if within epoch
		epochPositionRational.SetInt64(epochPosition)
		epochLengthRational.SetInt64(epochLength)
		epochPositionRational.Quo(&epochPositionRational, &epochLengthRational)
		interestRational.Mul(&interestRational, &epochPositionRational)
	}

	interestNew := interestRational.Num().Quo(interestRational.Num(), interestRational.Denom())

	interestNewInt := sdk.NewIntFromBigInt(interestNew.Add(interestNew, unpaidCollaterals.BigInt()))
	// round up to lowest digit if interest too low and rate not 0
	if interestNewInt.IsZero() && !interestRate.IsZero() {
		interestNewInt = sdk.NewInt(1)
	}

	return interestNewInt
}
