package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

//go:generate mockery --srcpkg . --name AuthorizationChecker --structname AuthorizationChecker --filename authorization_checker.go --with-expecter
type AuthorizationChecker interface {
	IsWhitelistingEnabled(ctx sdk.Context) bool
	CheckIfWhitelisted(ctx sdk.Context, creator string) bool
}

//go:generate mockery --srcpkg . --name PositionChecker --structname PositionChecker --filename position_checker.go --with-expecter
type PositionChecker interface {
	GetOpenPositionCount(ctx sdk.Context) uint64
	GetMaxOpenPositions(ctx sdk.Context) uint64
}

//go:generate mockery --srcpkg . --name PoolChecker --structname PoolChecker --filename pool_checker.go --with-expecter
type PoolChecker interface {
	GetPool(ctx sdk.Context, poolId uint64) (Pool, bool)
	IsPoolEnabled(ctx sdk.Context, poolId uint64) bool
	IsPoolClosed(ctx sdk.Context, poolId uint64) bool
	GetPoolOpenThreshold(ctx sdk.Context) math.LegacyDec
}

//go:generate mockery --srcpkg . --name OpenChecker --structname OpenChecker --filename open_checker.go --with-expecter
type OpenChecker interface {
	CheckUserAuthorization(ctx sdk.Context, msg *MsgOpen) error
	CheckMaxOpenPositions(ctx sdk.Context) error
	CheckPoolHealth(ctx sdk.Context, poolId uint64) error
	OpenLong(ctx sdk.Context, poolId uint64, msg *MsgOpen) (*Position, error)
	EmitOpenEvent(ctx sdk.Context, position *Position)
	SetPosition(ctx sdk.Context, position *Position) error
	CheckSamePosition(ctx sdk.Context, msg *MsgOpen) *Position
	GetOpenPositionCount(ctx sdk.Context) uint64
	GetMaxOpenPositions(ctx sdk.Context) uint64
}

//go:generate mockery --srcpkg . --name OpenLongChecker --structname OpenLongChecker --filename open_long_checker.go --with-expecter
type OpenLongChecker interface {
	GetMaxLeverageParam(ctx sdk.Context) sdk.Dec
	GetPool(ctx sdk.Context, poolId uint64) (Pool, bool)
	IsPoolEnabled(ctx sdk.Context, poolId uint64) bool
	GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error)
	HasSufficientPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool, assetDenom string, requiredAmount sdk.Int) bool
	CheckMinLiabilities(ctx sdk.Context, collateralTokenAmt sdk.Coin, eta sdk.Dec, pool Pool, ammPool ammtypes.Pool, borrowAsset string) error
	EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (sdk.Int, error)
	UpdatePoolHealth(ctx sdk.Context, pool *Pool) error
	UpdatePositionHealth(ctx sdk.Context, position Position, ammPool ammtypes.Pool) (sdk.Dec, error)
	GetSafetyFactor(ctx sdk.Context) sdk.Dec
	SetPool(ctx sdk.Context, pool Pool)
	GetAmmPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool, assetDenom string) (sdk.Int, error)
	CheckSamePosition(ctx sdk.Context, msg *MsgOpen) *Position
	SetPosition(ctx sdk.Context, position *Position) error
}

//go:generate mockery --srcpkg . --name CloseLongChecker --structname CloseLongChecker --filename close_long_checker.go --with-expecter
type CloseLongChecker interface {
	GetPosition(ctx sdk.Context, positionAddress string, id uint64) (Position, error)
	GetPool(
		ctx sdk.Context,
		poolId uint64,

	) (val Pool, found bool)
	GetAmmPool(ctx sdk.Context, poolId uint64, tradingAsset string) (ammtypes.Pool, error)
	HandleInterest(ctx sdk.Context, position *Position, pool *Pool, ammPool ammtypes.Pool, collateralAsset string, custodyAsset string) error
	EstimateAndRepay(ctx sdk.Context, position Position, pool Pool, ammPool ammtypes.Pool, collateralAsset string, custodyAsset string) (sdk.Int, error)
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
	GetPoolSnapshotOrSet(ctx sdk.Context, pool ammtypes.Pool) (val ammtypes.Pool)

	CalcOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (sdk.Coin, error)
	CalcInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (tokenIn sdk.Coin, err error)
	JoinPoolNoSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareOutAmount sdk.Int, tokenInMaxs sdk.Coins, noRemaining bool) (tokenIn sdk.Coins, sharesOut sdk.Int, err error)
	ExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount sdk.Int, tokenOutMins sdk.Coins, tokenOutDenom string) (exitCoins sdk.Coins, err error)
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
	GetDebt(ctx sdk.Context, addr sdk.AccAddress) stablestaketypes.Debt
	UpdateInterestStackedByAddress(ctx sdk.Context, addr sdk.AccAddress) stablestaketypes.Debt
	Borrow(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error
	Repay(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator string) (val commitmenttypes.Commitments, found bool)
}
