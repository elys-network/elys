package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TVL(ctx sdk.Context, oracleKeeper OracleKeeper, bankKeeper BankKeeper, baseCurrency string) math.LegacyDec {
	stableStakePoolAddress := authtypes.NewModuleAddress(ModuleName)

	balance := bankKeeper.GetBalance(ctx, stableStakePoolAddress, baseCurrency)
	price := oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	return price.MulInt(balance.Amount)
}
