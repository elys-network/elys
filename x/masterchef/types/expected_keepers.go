package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ctypes "github.com/elys-network/elys/v6/x/commitment/types"
	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
	parametertypes "github.com/elys-network/elys/v6/x/parameter/types"
	stabletypes "github.com/elys-network/elys/v6/x/stablestake/types"
	tokenomictypes "github.com/elys-network/elys/v6/x/tokenomics/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CommitmentKeeper
type CommitmentKeeper interface {
	IterateCommitments(sdk.Context, func(ctypes.Commitments) (stop bool))
	GetCommitments(sdk.Context, sdk.AccAddress) ctypes.Commitments
	SetCommitments(sdk.Context, ctypes.Commitments)
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnEdenBoost(ctx sdk.Context, creator sdk.AccAddress, denom string, amount math.Int) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetParams(sdk.Context) ctypes.Params
}

// Staking keeper
type StakingKeeper interface {
	TotalBondedTokens(context.Context) math.Int // total bonded tokens within the validator set
	// iterate through all delegations from one delegator by validator-AccAddress,
	// execute func for each validator
	IterateDelegations(ctx context.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	// get a particular validator by operator address
	Validator(context.Context, sdk.ValAddress) stakingtypes.ValidatorI
	// GetDelegatorDelegations returns a given amount of all the delegations from a delegator.
	GetDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation)
	// get a particular validator by consensus address
	ValidatorByConsAddr(context.Context, sdk.ConsAddress) stakingtypes.ValidatorI
	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(context.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}

// AmmKeeper defines the expected interface needed to swap tokens
type AmmKeeper interface {
	GetBestPoolWithDenoms(ctx sdk.Context, denoms []string, usesOracle bool) (pool ammtypes.Pool, found bool)
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
	GetPoolWithAccountedBalance(ctx sdk.Context, poolId uint64) (val ammtypes.SnapshotPool)

	GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) osmomath.BigDec
	GetTokenPrice(ctx sdk.Context, tokenInDenom, baseCurrency string) osmomath.BigDec
	InternalSwapExactAmountIn(
		ctx sdk.Context,
		sender sdk.AccAddress,
		recipient sdk.AccAddress,
		pool ammtypes.Pool,
		tokenIn sdk.Coin,
		tokenOutDenom string,
		tokenOutMinAmount math.Int,
		swapFee osmomath.BigDec,
		takersFee osmomath.BigDec,
	) (tokenOutAmount math.Int, weightBalanceReward sdk.Coin, err error)
	SwapByDenom(ctx sdk.Context, msg *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error)
	CalculateCoinsUSDValue(ctx sdk.Context, coins sdk.Coins) osmomath.BigDec
}

// OracleKeeper defines the expected interface needed to retrieve price info
type OracleKeeper interface {
	GetAssetPrice(ctx sdk.Context, asset string) (osmomath.BigDec, bool)
	GetDenomPrice(ctx sdk.Context, denom string) osmomath.BigDec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
	//SetPool(ctx sdk.Context, pool oracletypes.Pool)
	//SetAccountedPool(ctx sdk.Context, accountedPool oracletypes.AccountedPool)
	//CurrencyPairProviders(ctx sdk.Context) oracletypes.CurrencyPairProvidersList
	//SetCurrencyPairProviders(ctx sdk.Context, currencyPairProviders oracletypes.CurrencyPairProvidersList)
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
}

// AccountedPoolKeeper
type AccountedPoolKeeper interface {
	GetAccountedBalance(sdk.Context, uint64, string) math.Int
}

// AssetProfileKeeper defines the expected interface needed to retrieve denom info
type AssetProfileKeeper interface {
	GetEntry(ctx sdk.Context, baseDenom string) (val assetprofiletypes.Entry, found bool)
	// GetUsdcDenom returns USDC denom
	GetUsdcDenom(ctx sdk.Context) (string, bool)
}

// StableStakeKeeper defines the expected stablestake keeper used for simulations (noalias)
type StableStakeKeeper interface {
	GetParams(ctx sdk.Context) (params stabletypes.Params)
	BorrowRatio(goCtx context.Context, req *stabletypes.QueryBorrowRatioRequest) (*stabletypes.QueryBorrowRatioResponse, error)
	TVL(ctx sdk.Context, poolId uint64) osmomath.BigDec
	AllTVL(ctx sdk.Context) osmomath.BigDec
	GetTotalAndPerDenomTVL(ctx sdk.Context) (totalTVL osmomath.BigDec, denomTVL sdk.Coins)
	IterateLiquidityPools(sdk.Context, func(stabletypes.Pool) bool)
	GetPoolByDenom(ctx sdk.Context, denom string) (stabletypes.Pool, bool)
	GetPool(ctx sdk.Context, poolId uint64) (pool stabletypes.Pool, found bool)
}

// TokenomicsKeeper defines the expected tokenomics keeper used for simulations (noalias)
type TokenomicsKeeper interface {
	GetAllTimeBasedInflation(ctx sdk.Context) (list []tokenomictypes.TimeBasedInflation)
}

type ParameterKeeper interface {
	GetParams(ctx sdk.Context) (params parametertypes.Params)
}
