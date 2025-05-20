package types

import (
	"context"

	"github.com/osmosis-labs/osmosis/osmomath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v4/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/v4/x/oracle/types"
)

type AmmKeeper interface {
	GetParams(ctx sdk.Context) (params ammtypes.Params)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx context.Context, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// AssetProfileKeeper defines the expected interfaces
type AssetProfileKeeper interface {
	// SetEntry set a specific entry in the store from its index
	SetEntry(ctx sdk.Context, entry assetprofiletypes.Entry)
	// GetEntry returns a entry from its index
	GetEntry(ctx sdk.Context, baseDenom string) (val assetprofiletypes.Entry, found bool)
	// GetUsdcDenom returns USDC denom
	GetUsdcDenom(ctx sdk.Context) (string, bool)
	GetEntryByDenom(ctx sdk.Context, denom string) (val assetprofiletypes.Entry, found bool)
}

// OracleKeeper defines the expected interface needed to retrieve price info
type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (osmomath.BigDec, bool)
	GetDenomPrice(ctx sdk.Context, denom string) osmomath.BigDec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}
