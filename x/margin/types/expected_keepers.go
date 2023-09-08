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
	GetMaxOpenPositions(ctx sdk.Context) uint64
}

//go:generate mockery --srcpkg . --name PoolChecker --structname PoolChecker --filename pool_checker.go --with-expecter
type PoolChecker interface {
	GetPool(ctx sdk.Context, poolId uint64) (Pool, bool)
	IsPoolEnabled(ctx sdk.Context, poolId uint64) bool
	IsPoolClosed(ctx sdk.Context, poolId uint64) bool
	GetPoolOpenThreshold(ctx sdk.Context) math.LegacyDec
}

//go:generate mockery --srcpkg . --name OpenLongChecker --structname OpenLongChecker --filename open_long_checker.go --with-expecter
type OpenLongChecker interface {
	GetMaxLeverageParam(ctx sdk.Context) sdk.Dec
	GetNonNativeAsset(collateralAsset string, borrowAsset string) string
	GetPool(ctx sdk.Context, poolId uint64) (Pool, bool)
	IsPoolEnabled(ctx sdk.Context, poolId uint64) bool
	GetAmmPool(ctx sdk.Context, poolId uint64, nonNativeAsset string) (ammtypes.Pool, error)
	HasSufficientPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool, assetDenom string, requiredAmount sdk.Int) bool
	CheckMinLiabilities(ctx sdk.Context, collateralTokenAmt sdk.Coin, eta sdk.Dec, pool Pool, ammPool ammtypes.Pool, borrowAsset string) error
	EstimateSwap(ctx sdk.Context, leveragedAmtTokenIn sdk.Coin, borrowAsset string, ammPool ammtypes.Pool) (sdk.Int, error)
	EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (sdk.Int, error)
	Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Int, custodyAmount sdk.Int, mtp *MTP, ammPool *ammtypes.Pool, pool *Pool, eta sdk.Dec) error
	UpdatePoolHealth(ctx sdk.Context, pool *Pool) error
	TakeInCustody(ctx sdk.Context, mtp MTP, pool *Pool) error
	UpdateMTPHealth(ctx sdk.Context, mtp MTP, ammPool ammtypes.Pool) (sdk.Dec, error)
	GetSafetyFactor(ctx sdk.Context) sdk.Dec
	SetPool(ctx sdk.Context, pool Pool)
	GetAmmPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool, assetDenom string) (sdk.Int, error)
	CheckLongingAssets(ctx sdk.Context, collateralAsset string, borrowAsset string) error
}

//go:generate mockery --srcpkg . --name CloseLongChecker --structname CloseLongChecker --filename close_long_checker.go --with-expecter
type CloseLongChecker interface {
	GetMTP(ctx sdk.Context, mtpAddress string, id uint64) (MTP, error)
	GetPool(
		ctx sdk.Context,
		poolId uint64,

	) (val Pool, found bool)
	GetAmmPool(ctx sdk.Context, poolId uint64, nonNativeAsset string) (ammtypes.Pool, error)
	HandleInterest(ctx sdk.Context, mtp *MTP, pool *Pool, ammPool ammtypes.Pool) error
	TakeOutCustody(ctx sdk.Context, mtp MTP, pool *Pool) error
	EstimateAndRepay(ctx sdk.Context, mtp MTP, pool Pool, ammPool ammtypes.Pool) (sdk.Int, error)
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
