package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func TVL(ctx sdk.Context, accountKeeper AccountKeeper, bankKeeper BankKeeper) sdk.Int {
	stableStakePoolAddress := accountKeeper.GetModuleAddress(ModuleName)

	balance := bankKeeper.GetBalance(ctx, stableStakePoolAddress, ptypes.BaseCurrency)
	return balance.Amount
}
