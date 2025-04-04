package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) InterestRateComputationForPool(ctx sdk.Context, pool types.Pool) sdkmath.LegacyDec {
	if pool.TotalValue.IsZero() {
		return pool.InterestRate
	}

	interestRateMax := pool.InterestRateMax
	interestRateMin := pool.InterestRateMin
	interestRateIncrease := pool.InterestRateIncrease
	interestRateDecrease := pool.InterestRateDecrease
	healthGainFactor := pool.HealthGainFactor
	prevInterestRate := pool.InterestRate

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	depositDenom := pool.GetDepositDenom()
	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)

	// rate = minRate + (min(borrowRatio, param * maxAllowed) / (param * maxAllowed)) * (maxRate - minRate)
	borrowRatio := sdkmath.LegacyZeroDec()
	if pool.TotalValue.IsPositive() {
		borrowRatio = (pool.TotalValue.Sub(balance.Amount).ToLegacyDec()).Quo(pool.TotalValue.ToLegacyDec())
	}

	maxAllowed := pool.MaxLeverageRatio.Mul(healthGainFactor)
	if borrowRatio.GT(maxAllowed) {
		borrowRatio = maxAllowed
	}

	if maxAllowed.IsZero() {
		borrowRatio = sdkmath.LegacyZeroDec()
		maxAllowed = sdkmath.LegacyOneDec()
	}

	targetInterestRate := interestRateMin.Add((borrowRatio.Quo(maxAllowed).Mul((interestRateMax.Sub(interestRateMin)))))

	interestRateChange := targetInterestRate.Sub(prevInterestRate)
	interestRate := prevInterestRate
	if interestRateChange.GTE(interestRateDecrease.Mul(sdkmath.LegacyNewDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
		interestRate = targetInterestRate
	} else if interestRateChange.GT(interestRateIncrease) {
		interestRate = prevInterestRate.Add(interestRateIncrease)
	} else if interestRateChange.LT(interestRateDecrease.Mul(sdkmath.LegacyNewDec(-1))) {
		interestRate = prevInterestRate.Sub(interestRateDecrease)
	}

	newInterestRate := interestRate

	if interestRate.GT(interestRateMin) && interestRate.LT(interestRateMax) {
		newInterestRate = interestRate
	} else if interestRate.LTE(interestRateMin) {
		newInterestRate = interestRateMin
	} else if interestRate.GTE(interestRateMax) {
		newInterestRate = interestRateMax
	}

	return newInterestRate
}
