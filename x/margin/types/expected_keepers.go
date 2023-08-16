package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

//go:generate mockery --srcpkg . --name AuthorizationChecker --structname AuthorizationChecker --filename authorization_checker.go --with-expecter
type AuthorizationChecker interface {
	IsWhitelistingEnabled(ctx sdk.Context) bool
	CheckIfWhitelisted(ctx sdk.Context, creator string) bool
}

//go:generate mockery --srcpkg . --name PositionChecker --structname PositionChecker --filename position_checker.go --with-expecter
type PositionChecker interface {
	GetOpenMTPCount(ctx sdk.Context) uint64
	GetMaxOpenPositions(ctx sdk.Context) int
}

//go:generate mockery --srcpkg . --name PoolChecker --structname PoolChecker --filename pool_checker.go --with-expecter
type PoolChecker interface {
	GetPool(ctx sdk.Context, poolId uint64) (Pool, bool)
	IsPoolEnabled(ctx sdk.Context, poolId uint64) bool
	IsPoolClosed(ctx sdk.Context, poolId uint64) bool
	GetPoolOpenThreshold(ctx sdk.Context) math.LegacyDec
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
//
//go:generate mockery --srcpkg . --name AccountKeeper --structname AccountKeeper --filename account_keeper.go --with-expecter
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

// AmmKeeper defines the expected interface needed to swap tokens
//
//go:generate mockery --srcpkg . --name AmmKeeper --structname AmmKeeper --filename amm_keeper.go --with-expecter
type AmmKeeper interface {
	// Get pool Ids that contains the denom in pool assets
	GetAllPoolIdsWithDenom(sdk.Context, string) []uint64
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
//
//go:generate mockery --srcpkg . --name BankKeeper --structname BankKeeper --filename bank_keeper.go --with-expecter
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}
