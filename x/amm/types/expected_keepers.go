package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	NewAccount(context.Context, sdk.AccountI) sdk.AccountI
	GetAccount(goCtx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(goCtx context.Context, acc sdk.AccountI)
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(goCtx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(goCtx context.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(goCtx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(goCtx context.Context, name string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(goCtx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(goCtx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SetDenomMetaData(goCtx context.Context, denomMetaData banktypes.Metadata)
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SendCoins(goCtx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

// OracleKeeper defines the expected interface needed to retrieve price info
//
//go:generate mockery --srcpkg . --name OracleKeeper --structname OracleKeeper --filename oracle_keeper.go --with-expecter
type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdkmath.LegacyDec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
	SetPool(ctx sdk.Context, pool oracletypes.Pool)
	SetAccountedPool(ctx sdk.Context, accountedPool oracletypes.AccountedPool)
}

// AssetProfileKeeper defines the expected interfaces
type AssetProfileKeeper interface {
	// SetEntry set a specific entry in the store from its index
	SetEntry(ctx sdk.Context, entry atypes.Entry)
	// GetEntry returns a entry from its index
	GetEntry(ctx sdk.Context, baseDenom string) (val atypes.Entry, found bool)
	// GetEntryByDenom returns a entry from its denom value
	GetEntryByDenom(ctx sdk.Context, denom string) (val atypes.Entry, found bool)
	// GetUsdcDenom returns USDC denom
	GetUsdcDenom(ctx sdk.Context) (usdcDenom string, found bool)
}

// AccountedPoolKeeper defines the expected interfaces
//
//go:generate mockery --srcpkg . --name AccountedPoolKeeper --structname AccountedPoolKeeper --filename accounted_pool_keeper.go --with-expecter
type AccountedPoolKeeper interface {
	GetAccountedBalance(sdk.Context, uint64, string) sdkmath.Int
}

type TierKeeper interface {
	GetMembershipTier(ctx sdk.Context, user sdk.AccAddress) (total_portfolio sdkmath.LegacyDec, tier string, discount sdkmath.LegacyDec)

	CalculateUSDValue(ctx sdk.Context, denom string, amount sdkmath.Int) sdkmath.LegacyDec
}
