package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) InterestRateComputationForPool(ctx sdk.Context, pool types.Pool) osmomath.BigDec {
	if pool.TotalValue.IsZero() {
		return pool.GetBigDecInterestRate()
	}

	interestRateMax := pool.GetBigDecInterestRateMax()
	interestRateMin := pool.GetBigDecInterestRateMin()
	interestRateIncrease := pool.GetBigDecInterestRateIncrease()
	interestRateDecrease := pool.GetBigDecInterestRateDecrease()
	healthGainFactor := pool.GetBigDecHealthGainFactor()
	prevInterestRate := pool.GetBigDecInterestRate()

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	depositDenom := pool.GetDepositDenom()
	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)

	// rate = minRate + (min(borrowRatio, param * maxAllowed) / (param * maxAllowed)) * (maxRate - minRate)
	borrowRatio := osmomath.ZeroBigDec()
	if pool.TotalValue.IsPositive() {
		borrowRatio = osmomath.BigDecFromSDKInt(pool.TotalValue.Sub(balance.Amount)).Quo(pool.GetBigDecTotalValue())
	}

	maxAllowed := pool.GetBigDecMaxLeverageRatio().Mul(healthGainFactor)
	if borrowRatio.GT(maxAllowed) {
		borrowRatio = maxAllowed
	}

	if maxAllowed.IsZero() {
		borrowRatio = osmomath.ZeroBigDec()
		maxAllowed = osmomath.OneBigDec()
	}

	targetInterestRate := interestRateMin.Add((borrowRatio.Quo(maxAllowed).Mul((interestRateMax.Sub(interestRateMin)))))

	interestRateChange := targetInterestRate.Sub(prevInterestRate)
	interestRate := prevInterestRate
	if interestRateChange.GTE(interestRateDecrease.Mul(osmomath.NewBigDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
		interestRate = targetInterestRate
	} else if interestRateChange.GT(interestRateIncrease) {
		interestRate = prevInterestRate.Add(interestRateIncrease)
	} else if interestRateChange.LT(interestRateDecrease.Mul(osmomath.NewBigDec(-1))) {
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
