package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
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
		discount sdkmath.LegacyDec,
		overrideSwapFee sdkmath.LegacyDec,
	) (sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coin, sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coin, sdkmath.LegacyDec, sdkmath.LegacyDec, error)
	Balance(goCtx context.Context, req *ammtypes.QueryBalanceRequest) (*ammtypes.QueryBalanceResponse, error)
	GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) sdkmath.LegacyDec
	CalculateUSDValue(ctx sdk.Context, denom string, amount sdkmath.Int) sdkmath.LegacyDec
	CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) sdkmath.LegacyDec
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx context.Context, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator sdk.AccAddress) commitmenttypes.Commitments
	CommitmentVestingInfo(goCtx context.Context, req *commitmenttypes.QueryCommitmentVestingInfoRequest) (*commitmenttypes.QueryCommitmentVestingInfoResponse, error)
	GetAllCommitmentsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]*commitmenttypes.Commitments, *query.PageResponse, error)
}

type TierKeeper interface {
	GetMembershipTier(ctx sdk.Context, user sdk.AccAddress) (total_portfolio sdkmath.LegacyDec, tier string, discount sdkmath.LegacyDec)
	CalculateUSDValue(ctx sdk.Context, denom string, amount sdkmath.Int) sdkmath.LegacyDec
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
