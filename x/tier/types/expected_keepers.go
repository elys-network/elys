package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	tradeshieldtypes "github.com/elys-network/elys/x/tradeshield/types"
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
	GetAssetPrice(ctx sdk.Context, asset string) (oracletypes.Price, bool)
	GetAssetPriceFromDenom(ctx sdk.Context, denom string) math.LegacyDec
	GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val oracletypes.PriceFeeder, found bool)
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator sdk.AccAddress) commitmenttypes.Commitments
	CommitmentVestingInfo(goCtx context.Context, req *commitmenttypes.QueryCommitmentVestingInfoRequest) (*commitmenttypes.QueryCommitmentVestingInfoResponse, error)
	GetAllCommitmentsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]*commitmenttypes.Commitments, *query.PageResponse, error)
}

type PerpetualKeeper interface {
	GetAllMTPsForAddress(ctx sdk.Context, mtpAddress sdk.AccAddress) []*perpetualtypes.MTP
	GetMTPsForAddressWithPagination(ctx sdk.Context, mtpAddress sdk.AccAddress, pagination *query.PageRequest) ([]*perpetualtypes.MtpAndPrice, *query.PageResponse, error)
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
	PoolExtraInfo(ctx sdk.Context, pool ammtypes.Pool) ammtypes.PoolExtraInfo
	InRouteByDenom(goCtx context.Context, req *ammtypes.QueryInRouteByDenomRequest) (*ammtypes.QueryInRouteByDenomResponse, error)
	CalcInRouteSpotPrice(ctx sdk.Context,
		tokenIn sdk.Coin,
		routes []*ammtypes.SwapAmountInRoute,
		discount math.LegacyDec,
		overrideSwapFee math.LegacyDec,
	) (math.LegacyDec, math.LegacyDec, sdk.Coin, math.LegacyDec, math.LegacyDec, sdk.Coin, math.LegacyDec, math.LegacyDec, error)
	Balance(goCtx context.Context, req *ammtypes.QueryBalanceRequest) (*ammtypes.QueryBalanceResponse, error)
	GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) math.LegacyDec
	CalculateUSDValue(ctx sdk.Context, denom string, amount math.Int) math.LegacyDec
	CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) math.LegacyDec
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
}

type TradeshieldKeeper interface {
	GetPendingPerpetualOrdersForAddress(ctx sdk.Context, address string, status *tradeshieldtypes.Status, pagination *query.PageRequest) ([]tradeshieldtypes.PerpetualOrder, *query.PageResponse, error)
	GetPendingSpotOrdersForAddress(ctx sdk.Context, address string, status *tradeshieldtypes.Status, pagination *query.PageRequest) ([]tradeshieldtypes.SpotOrder, *query.PageResponse, error)
}
