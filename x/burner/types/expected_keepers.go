package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetAllDenomMetaData(ctx sdk.Context) []banktypes.Metadata
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	// Methods imported from bank should be defined here
}
