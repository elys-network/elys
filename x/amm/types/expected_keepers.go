package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authtypes.AccountI) authtypes.AccountI
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

// OracleKeeper defines the expected interface needed to retrieve price info
type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdk.Dec
	GetPriceFeeder(ctx sdk.Context, feeder string) (val oracletypes.PriceFeeder, found bool)
}

// AssetProfileKeeper defines the expected interfaces
type AssetProfileKeeper interface {
	// SetEntry set a specific entry in the store from its index
	SetEntry(ctx sdk.Context, entry atypes.Entry)
	// GetEntry returns a entry from its index
	GetEntry(ctx sdk.Context, baseDenom string) (val atypes.Entry, found bool)
}

// AccountedPoolKeeper defines the expected interfaces
type AccountedPoolKeeper interface {
	GetAccountedBalance(sdk.Context, uint64, string) sdk.Int
}
