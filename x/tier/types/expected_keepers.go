package types

import (
	"context"

	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/v6/x/commitment/types"
	estakingtypes "github.com/elys-network/elys/v6/x/estaking/types"
	leveragelptypes "github.com/elys-network/elys/v6/x/leveragelp/types"
	mastercheftypes "github.com/elys-network/elys/v6/x/masterchef/types"
	perpetualtypes "github.com/elys-network/elys/v6/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/v6/x/stablestake/types"
	tradeshieldtypes "github.com/elys-network/elys/v6/x/tradeshield/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(goCtx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

type OracleKeeper interface {
	GetAssetInfo(ctx sdk.Context, denom string) (val oracletypes.AssetInfo, found bool)
	GetAssetPrice(ctx sdk.Context, asset string) (math.LegacyDec, bool)
	GetDenomPrice(ctx sdk.Context, denom string) osmomath.BigDec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator sdk.AccAddress) commitmenttypes.Commitments
	CommitmentVestingInfo(goCtx context.Context, req *commitmenttypes.QueryCommitmentVestingInfoRequest) (*commitmenttypes.QueryCommitmentVestingInfoResponse, error)
	GetAllCommitmentsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]*commitmenttypes.Commitments, *query.PageResponse, error)
}

type PerpetualKeeper interface {
	GetAllMTPForAddress(ctx sdk.Context, mtpAddress sdk.AccAddress) ([]*perpetualtypes.MtpAndPrice, error)
}

// AssetProfileKeeper defines the expected interface needed to retrieve denom info
type AssetProfileKeeper interface {
	GetEntry(ctx sdk.Context, baseDenom string) (val assetprofiletypes.Entry, found bool)
	GetAllEntry(ctx sdk.Context) (list []assetprofiletypes.Entry)
	// GetUsdcDenom returns USDC denom
	GetUsdcDenom(ctx sdk.Context) (string, bool)
	// GetEntryByDenom returns a entry from its denom value
	GetEntryByDenom(ctx sdk.Context, denom string) (val assetprofiletypes.Entry, found bool)
}

type AmmKeeper interface {
	// GetPool returns a pool from its index
	GetPool(sdk.Context, uint64) (ammtypes.Pool, bool)
	// Get all pools
	GetAllPool(sdk.Context) []ammtypes.Pool
	// IterateCommitments iterates over all Commitments and performs a callback.
	IterateLiquidityPools(sdk.Context, func(ammtypes.Pool) bool)
	PoolExtraInfo(ctx sdk.Context, pool ammtypes.Pool, days int) ammtypes.PoolExtraInfo
	InRouteByDenom(goCtx context.Context, req *ammtypes.QueryInRouteByDenomRequest) (*ammtypes.QueryInRouteByDenomResponse, error)
	CalcInRouteSpotPrice(ctx sdk.Context,
		tokenIn sdk.Coin,
		routes []*ammtypes.SwapAmountInRoute,
		discount osmomath.BigDec,
		overrideSwapFee osmomath.BigDec,
	) (osmomath.BigDec, osmomath.BigDec, sdk.Coin, osmomath.BigDec, osmomath.BigDec, sdk.Coin, osmomath.BigDec, osmomath.BigDec, error)
	Balance(goCtx context.Context, req *ammtypes.QueryBalanceRequest) (*ammtypes.QueryBalanceResponse, error)
	GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) osmomath.BigDec
	CalculateUSDValue(ctx sdk.Context, denom string, amount math.Int) osmomath.BigDec
	CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) osmomath.BigDec
}

type EstakingKeeper interface {
	Rewards(goCtx context.Context, req *estakingtypes.QueryRewardsRequest) (*estakingtypes.QueryRewardsResponse, error)
}

type MasterchefKeeper interface {
	UserPendingReward(goCtx context.Context, req *mastercheftypes.QueryUserPendingRewardRequest) (*mastercheftypes.QueryUserPendingRewardResponse, error)
}

type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
	GetUnbondingDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) (unbondingDelegations []stakingtypes.UnbondingDelegation, err error)
	GetDelegatorValidators(ctx context.Context, delegatorAddr sdk.AccAddress, maxRetrieve uint32) (stakingtypes.Validators, error)
	GetAllDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress) ([]stakingtypes.Delegation, error)
}

type LeverageLpKeeper interface {
	GetPositionsForAddress(ctx sdk.Context, positionAddress sdk.AccAddress, pagination *query.PageRequest) ([]*leveragelptypes.Position, *query.PageResponse, error)
}

type StablestakeKeeper interface {
	GetParams(ctx sdk.Context) (params stablestaketypes.Params)
	GetDebt(ctx sdk.Context, addr sdk.AccAddress, poolId uint64) stablestaketypes.Debt
	UpdateInterestAndGetDebt(ctx sdk.Context, addr sdk.AccAddress, poolId uint64, borrowingForPool uint64) stablestaketypes.Debt
	GetPool(ctx sdk.Context, poolId uint64) (pool stablestaketypes.Pool, found bool)
	CalculateRedemptionRateForPool(ctx sdk.Context, pool stablestaketypes.Pool) osmomath.BigDec
}

type TradeshieldKeeper interface {
	GetPendingPerpetualOrdersForAddress(ctx sdk.Context, address string, status *tradeshieldtypes.Status, pagination *query.PageRequest) ([]tradeshieldtypes.PerpetualOrder, *query.PageResponse, error)
	GetPendingSpotOrdersForAddress(ctx sdk.Context, address string, status *tradeshieldtypes.Status, pagination *query.PageRequest) ([]tradeshieldtypes.SpotOrder, *query.PageResponse, error)
}
