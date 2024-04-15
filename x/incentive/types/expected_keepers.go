package types

import (
	context "context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
	tokenomictypes "github.com/elys-network/elys/x/tokenomics/types"
)

// CommitmentKeeper
type CommitmentKeeper interface {
	// Initiate commitment according to standard staking
	BeforeDelegationCreated(sdk.Context, string, string) error
	// Iterate all commitment
	IterateCommitments(sdk.Context, func(ctypes.Commitments) (stop bool))
	// Update commitment
	SetCommitments(sdk.Context, ctypes.Commitments)
	// Get commitment
	GetCommitments(sdk.Context, string) ctypes.Commitments
	// Burn eden boost
	BurnEdenBoost(ctx sdk.Context, creator string, denom string, amount math.Int) (ctypes.Commitments, error)
}

// Staking keeper
type StakingKeeper interface {
	TotalBondedTokens(sdk.Context) math.Int // total bonded tokens within the validator set
	// iterate through all delegations from one delegator by validator-AccAddress,
	// execute func for each validator
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	// get a particular validator by operator address
	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
	// GetDelegatorDelegations returns a given amount of all the delegations from a delegator.
	GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation)
	// get a particular validator by consensus address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI
	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

// AmmKeeper defines the expected interface needed to swap tokens
type AmmKeeper interface {
	// UpdatePoolForSwap takes a pool, sender, and tokenIn, tokenOut amounts
	// It then updates the pool's balances to the new reserve amounts, and
	// sends the in tokens from the sender to the pool, and the out tokens from the pool to the sender.
	UpdatePoolForSwap(
		ctx sdk.Context,
		pool ammtypes.Pool,
		sender sdk.AccAddress,
		recipient sdk.AccAddress,
		tokenIn sdk.Coin,
		tokenOut sdk.Coin,
		swapFeeIn sdk.Dec,
		swapFeeOut sdk.Dec,
		weightBalanceBonus sdk.Dec,
	) (math.Int, error)
	GetBestPoolWithDenoms(ctx sdk.Context, denoms []string) (pool ammtypes.Pool, found bool)
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
	GetPoolSnapshotOrSet(ctx sdk.Context, pool ammtypes.Pool) (val ammtypes.Pool)

	SwapOutAmtGivenIn(
		ctx sdk.Context, poolId uint64,
		oracleKeeper ammtypes.OracleKeeper,
		snapshot *ammtypes.Pool,
		tokensIn sdk.Coins,
		tokenOutDenom string,
		swapFee sdk.Dec,
	) (tokenOut sdk.Coin, slippage, slippageAmount sdk.Dec, weightBalanceBonus sdk.Dec, err error)
	CalcOutAmtGivenIn(ctx sdk.Context, poolId uint64, oracle ammtypes.OracleKeeper, snapshot *ammtypes.Pool, tokensIn sdk.Coins, tokenOutDenom string, swapFee sdk.Dec) (sdk.Coin, sdk.Dec, error)
	GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) math.LegacyDec
	GetTokenPrice(ctx sdk.Context, tokenInDenom, baseCurrency string) math.LegacyDec
}

// OracleKeeper defines the expected interface needed to retrieve price info
type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdk.Dec
	GetPriceFeeder(ctx sdk.Context, feeder string) (val oracletypes.PriceFeeder, found bool)
}

// AccountedPoolKeeper
type AccountedPoolKeeper interface {
	GetAccountedBalance(sdk.Context, uint64, string) math.Int
}

// AssetProfileKeeper defines the expected interface needed to retrieve denom info
type AssetProfileKeeper interface {
	GetEntry(ctx sdk.Context, baseDenom string) (val aptypes.Entry, found bool)
}

// EpochsKeeper defines the expected epochs keeper used for simulations (noalias)
type EpochsKeeper interface {
	GetEpochInfo(ctx sdk.Context, identifier string) (epochstypes.EpochInfo, bool)
}

// StableStakeKeeper defines the expected epochs keeper used for simulations (noalias)
type StableStakeKeeper interface {
	GetParams(ctx sdk.Context) (params stabletypes.Params)
	BorrowRatio(goCtx context.Context, req *stabletypes.QueryBorrowRatioRequest) (*stabletypes.QueryBorrowRatioResponse, error)
	TVL(ctx sdk.Context, oracleKeeper stabletypes.OracleKeeper, baseCurrency string) math.LegacyDec
}

// TokenomicsKeeper defines the expected tokenomics keeper used for simulations (noalias)
type TokenomicsKeeper interface {
	GetAllTimeBasedInflation(ctx sdk.Context) (list []tokenomictypes.TimeBasedInflation)
}
