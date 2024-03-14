package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

// AmmKeeper defines the expected interface needed to swap tokens
type AmmKeeper interface {
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
	GetPoolSnapshotOrSet(ctx sdk.Context, pool ammtypes.Pool) (val ammtypes.Pool)

	CalcOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (sdk.Coin, sdk.Dec, error)
	CalcInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (tokenIn sdk.Coin, slippage sdk.Dec, err error)
	JoinPoolNoSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareOutAmount math.Int, tokenInMaxs sdk.Coins) (tokenIn sdk.Coins, sharesOut math.Int, err error)
	ExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount math.Int, tokenOutMins sdk.Coins, tokenOutDenom string) (exitCoins sdk.Coins, err error)
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
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

// StableStakeKeeper defines the expected interface needed on stablestake
type StableStakeKeeper interface {
	GetParams(ctx sdk.Context) stablestaketypes.Params
	GetDepositDenom(ctx sdk.Context) string
	GetDebt(ctx sdk.Context, addr sdk.AccAddress) stablestaketypes.Debt
	UpdateInterestStackedByAddress(ctx sdk.Context, addr sdk.AccAddress) stablestaketypes.Debt
	Borrow(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error
	Repay(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error
	TVL(ctx sdk.Context, oracleKeeper stablestaketypes.OracleKeeper, baseCurrency string) math.LegacyDec
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator string) commitmenttypes.Commitments
}

// AssetProfileKeeper defines the expected interface needed to retrieve denom info
type AssetProfileKeeper interface {
	GetEntry(ctx sdk.Context, baseDenom string) (val aptypes.Entry, found bool)
}
