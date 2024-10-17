package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

//go:generate mockery --srcpkg . --name AuthorizationChecker --structname AuthorizationChecker --filename authorization_checker.go --with-expecter
type AuthorizationChecker interface {
	IsWhitelistingEnabled(ctx sdk.Context) bool
	CheckIfWhitelisted(ctx sdk.Context, creator sdk.AccAddress) bool
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

//go:generate mockery --srcpkg . --name OpenChecker --structname OpenChecker --filename open_checker.go --with-expecter
type OpenChecker interface {
	CheckUserAuthorization(ctx sdk.Context, msg *MsgOpen) error
	CheckMaxOpenPositions(ctx sdk.Context) error
	CheckLowPoolHealth(ctx sdk.Context, poolId uint64) error
	OpenDefineAssets(ctx sdk.Context, poolId uint64, msg *MsgOpen, baseCurrency string, isBroker bool) (*MTP, error)
	UpdateOpenPrice(ctx sdk.Context, mtp *MTP, ammPool ammtypes.Pool, baseCurrency string) error
	EmitOpenEvent(ctx sdk.Context, mtp *MTP)
	SetMTP(ctx sdk.Context, mtp *MTP) error
	CheckSameAssetPosition(ctx sdk.Context, msg *MsgOpen) *MTP
	GetOpenMTPCount(ctx sdk.Context) uint64
	GetMaxOpenPositions(ctx sdk.Context) uint64
}

//go:generate mockery --srcpkg . --name OpenDefineAssetsChecker --structname OpenDefineAssetsChecker --filename open_define_assets_checker.go --with-expecter
type OpenDefineAssetsChecker interface {
	GetMaxLeverageParam(ctx sdk.Context) sdk.Dec
	GetPool(ctx sdk.Context, poolId uint64) (Pool, bool)
	IsPoolEnabled(ctx sdk.Context, poolId uint64) bool
	GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error)
	EstimateSwap(ctx sdk.Context, leveragedAmtTokenIn sdk.Coin, borrowAsset string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, error)
	EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, error)
	Borrow(ctx sdk.Context, collateralAmount math.Int, custodyAmount math.Int, mtp *MTP, ammPool *ammtypes.Pool, pool *Pool, eta sdk.Dec, baseCurrency string, isBroker bool) error
	UpdatePoolHealth(ctx sdk.Context, pool *Pool) error
	TakeInCustody(ctx sdk.Context, mtp MTP, pool *Pool) error
	GetMTPHealth(ctx sdk.Context, mtp MTP, ammPool ammtypes.Pool, baseCurrency string) (sdk.Dec, error)
	GetSafetyFactor(ctx sdk.Context) sdk.Dec
	SetPool(ctx sdk.Context, pool Pool)
	CheckSameAssetPosition(ctx sdk.Context, msg *MsgOpen) *MTP
	SetMTP(ctx sdk.Context, mtp *MTP) error
	DestroyMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) error
	UpdateOpenPrice(ctx sdk.Context, mtp *MTP, ammPool ammtypes.Pool, baseCurrency string) error
	EmitOpenEvent(ctx sdk.Context, mtp *MTP)
}

//go:generate mockery --srcpkg . --name ClosePositionChecker --structname ClosePositionChecker --filename close_position_checker.go --with-expecter
type ClosePositionChecker interface {
	GetMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) (MTP, error)
	GetPool(
		ctx sdk.Context,
		poolId uint64,
	) (val Pool, found bool)
	GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error)
	SettleBorrowInterest(ctx sdk.Context, mtp *MTP, pool *Pool, ammPool ammtypes.Pool) (math.Int, error)
	TakeOutCustody(ctx sdk.Context, mtp MTP, pool *Pool, amount math.Int) error
	EstimateAndRepay(ctx sdk.Context, mtp *MTP, pool *Pool, ammPool *ammtypes.Pool, baseCurrency string, closingRatio sdk.Dec) (math.Int, error)
}

//go:generate mockery --srcpkg . --name CloseEstimationChecker --structname CloseEstimationChecker --filename close_estimation_checker.go --with-expecter
type CloseEstimationChecker interface {
	GetMTP(ctx sdk.Context, mtpAddress sdk.AccAddress, id uint64) (MTP, error)
	GetPool(
		ctx sdk.Context,
		poolId uint64,
	) (val Pool, found bool)
	GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error)
	EstimateSwap(ctx sdk.Context, leveragedAmtTokenIn sdk.Coin, borrowAsset string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, error)
	EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, error)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
//
//go:generate mockery --srcpkg . --name AccountKeeper --structname AccountKeeper --filename account_keeper.go --with-expecter
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

//go:generate mockery --srcpkg . --name LeverageLpKeeper --structname LeverageLpKeeper --filename leverageLP_keeper.go --with-expecter
type LeverageLpKeeper interface {
	GetPool(ctx sdk.Context, poolId uint64) (leveragelpmoduletypes.Pool, bool)
	// Methods imported from account should be defined here
}

// AmmKeeper defines the expected interface needed to swap tokens
//
//go:generate mockery --srcpkg . --name AmmKeeper --structname AmmKeeper --filename amm_keeper.go --with-expecter
type AmmKeeper interface {
	// Get first pool id that contains all denoms in pool assets
	GetBestPoolWithDenoms(ctx sdk.Context, denoms []string, usesOracle bool) (pool ammtypes.Pool, found bool)
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
	GetPoolSnapshotOrSet(ctx sdk.Context, pool ammtypes.Pool) (val ammtypes.Pool)

	CalcOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (sdk.Coin, sdk.Dec, error)
	CalcInAmtGivenOut(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (tokenIn sdk.Coin, slippage sdk.Dec, err error)

	AddToPoolBalance(ctx sdk.Context, pool *ammtypes.Pool, addShares math.Int, coins sdk.Coins) error
	RemoveFromPoolBalance(ctx sdk.Context, pool *ammtypes.Pool, removeShares math.Int, coins sdk.Coins) error
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

// AssetProfileKeeper defines the expected interfaces
//
//go:generate mockery --srcpkg . --name AssetProfileKeeper --structname AssetProfileKeeper --filename asset_profile_keeper.go --with-expecter
type AssetProfileKeeper interface {
	// GetEntry returns a entry from its index
	GetEntry(ctx sdk.Context, baseDenom string) (val atypes.Entry, found bool)
	// GetEntryByDenom returns a entry from its denom value
	GetEntryByDenom(ctx sdk.Context, denom string) (val atypes.Entry, found bool)
}

type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdk.Dec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}
