package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TVL(ctx sdk.Context, accountKeeper AccountKeeper, bankKeeper BankKeeper, baseCurrency string) math.Int {
	stableStakePoolAddress := authtypes.NewModuleAddress(ModuleName)

	balance := bankKeeper.GetBalance(ctx, stableStakePoolAddress, baseCurrency)
	return balance.Amount
}
