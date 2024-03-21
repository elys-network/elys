package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) InterestRateComputation(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)
	if params.TotalValue.IsZero() {
		return params.InterestRate
	}

	interestRateMax := params.InterestRateMax
	interestRateMin := params.InterestRateMin
	interestRateIncrease := params.InterestRateIncrease
	interestRateDecrease := params.InterestRateDecrease
	healthGainFactor := params.HealthGainFactor
	prevInterestRate := params.InterestRate

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	depositDenom := k.GetDepositDenom(ctx)
	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := params.TotalValue.Sub(balance.Amount)
	targetInterestRate := healthGainFactor.
		Mul(sdk.NewDecFromInt(borrowed)).
		Quo(sdk.NewDecFromInt(params.TotalValue))

	interestRateChange := targetInterestRate.Sub(prevInterestRate)
	interestRate := prevInterestRate
	if interestRateChange.GTE(interestRateDecrease.Mul(sdk.NewDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
		interestRate = targetInterestRate
	} else if interestRateChange.GT(interestRateIncrease) {
		interestRate = prevInterestRate.Add(interestRateIncrease)
	} else if interestRateChange.LT(interestRateDecrease.Mul(sdk.NewDec(-1))) {
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
