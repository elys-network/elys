package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdk.Dec
	GetPriceFeeder(ctx sdk.Context, feeder string) (val oracletypes.PriceFeeder, found bool)
}
