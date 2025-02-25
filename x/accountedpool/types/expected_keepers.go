package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// AmmKeeper defines the expected interface needed to swap tokens
//
//go:generate mockery --srcpkg . --name AmmKeeper --structname AmmKeeper --filename amm_keeper.go --with-expecter
type AmmKeeper interface {
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
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

// PerpetualKeeper defines the expected interface needed
//
//go:generate mockery --srcpkg . --name PerpetualKeeper --structname PerpetualKeeper --filename perpetual_keeper.go --with-expecter
type PerpetualKeeper interface {
	GetPool(ctx sdk.Context, poolId uint64) (perpetualtypes.Pool, bool)
	IsPoolEnabled(ctx sdk.Context, poolId uint64) bool
	IsPoolClosed(ctx sdk.Context, poolId uint64) bool
	GetAllMTPs(ctx sdk.Context) []perpetualtypes.MTP
}

type OracleKeeper interface {
	SetAccountedPool(ctx sdk.Context, accountedPool oracletypes.AccountedPool)
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}
