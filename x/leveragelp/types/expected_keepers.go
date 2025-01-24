package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(goCtx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// AmmKeeper defines the expected interface needed to swap tokens
type AmmKeeper interface {
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	JoinPoolEst(ctx sdk.Context, poolId uint64, tokenInMaxs sdk.Coins) (tokensIn sdk.Coins, sharesOut sdkmath.Int, slippage sdkmath.LegacyDec, weightBalanceBonus sdkmath.LegacyDec, err error)
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
	GetAccountedPoolSnapshotOrSet(ctx sdk.Context, pool ammtypes.Pool) (val ammtypes.Pool)

	CalcOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee sdkmath.LegacyDec) (sdk.Coin, sdkmath.LegacyDec, error)
	CalcInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut sdk.Coins, tokenInDenom string, swapFee sdkmath.LegacyDec) (tokenIn sdk.Coin, slippage sdkmath.LegacyDec, err error)
	JoinPoolNoSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareOutAmount sdkmath.Int, tokenInMaxs sdk.Coins) (tokenIn sdk.Coins, sharesOut sdkmath.Int, err error)
	ExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount sdkmath.Int, tokenOutMins sdk.Coins, tokenOutDenom string, isLiquidation, applyWeightBreakingFee bool) (exitCoins sdk.Coins, weightBalanceBonus sdkmath.LegacyDec, err error)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
//
//go:generate mockery --srcpkg . --name BankKeeper --structname BankKeeper --filename bank_keeper.go --with-expecter
type BankKeeper interface {
	GetAllBalances(goCtx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(goCtx context.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(goCtx context.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(goCtx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(goCtx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(goCtx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(goCtx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(goCtx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

// StableStakeKeeper defines the expected interface needed on stablestake
type StableStakeKeeper interface {
	GetParams(ctx sdk.Context) stablestaketypes.Params
	GetDepositDenom(ctx sdk.Context) string
	GetDebt(ctx sdk.Context, addr sdk.AccAddress) stablestaketypes.Debt
	UpdateInterestAndGetDebt(ctx sdk.Context, addr sdk.AccAddress) stablestaketypes.Debt
	Borrow(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error
	Repay(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error
	TVL(ctx sdk.Context, oracleKeeper stablestaketypes.OracleKeeper, baseCurrency string) sdkmath.LegacyDec
	GetInterest(ctx sdk.Context, startBlock uint64, startTime uint64, borrowed sdkmath.LegacyDec) sdkmath.Int
	CloseOnUnableToRepay(ctx sdk.Context, addr sdk.AccAddress) error
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator sdk.AccAddress) commitmenttypes.Commitments
}

// AssetProfileKeeper defines the expected interface needed to retrieve denom info
type AssetProfileKeeper interface {
	GetEntry(ctx sdk.Context, baseDenom string) (val assetprofiletypes.Entry, found bool)
	// GetUsdcDenom returns USDC denom
	GetUsdcDenom(ctx sdk.Context) (string, bool)
}

// MasterchefKeeper defines expected interface for masterchef keeper
type MasterchefKeeper interface {
	ClaimRewards(ctx sdk.Context, sender sdk.AccAddress, poolIds []uint64, recipient sdk.AccAddress) error
	UserPoolPendingReward(ctx sdk.Context, user sdk.AccAddress, poolId uint64) sdk.Coins
}

type AccountedPoolKeeper interface {
	GetAccountedBalance(sdk.Context, uint64, string) sdkmath.Int
}
