package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	atypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/v6/x/commitment/types"
	mastercheftypes "github.com/elys-network/elys/v6/x/masterchef/types"
	parametertypes "github.com/elys-network/elys/v6/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	NewAccount(context.Context, sdk.AccountI) sdk.AccountI
	GetAccount(goCtx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(goCtx context.Context, acc sdk.AccountI)
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
		discount osmomath.BigDec,
		overrideSwapFee osmomath.BigDec,
	) (osmomath.BigDec, osmomath.BigDec, sdk.Coin, osmomath.BigDec, osmomath.BigDec, sdk.Coin, osmomath.BigDec, osmomath.BigDec, error)
	Balance(goCtx context.Context, req *ammtypes.QueryBalanceRequest) (*ammtypes.QueryBalanceResponse, error)
	GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) osmomath.BigDec
	CalculateUSDValue(ctx sdk.Context, denom string, amount sdkmath.Int) osmomath.BigDec
	CalculateCoinsUSDValue(ctx sdk.Context, coins sdk.Coins) osmomath.BigDec
	CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) osmomath.BigDec
	JoinPoolNoSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareOutAmount sdkmath.Int, tokenInMaxs sdk.Coins) (tokenIn sdk.Coins, sharesOut sdkmath.Int, err error)
	ExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount sdkmath.Int, tokenOutMins sdk.Coins, tokenOutDenom string, isLiquidation, applyWeightBreakingFee bool) (exitCoins sdk.Coins, weightBalanceBonus osmomath.BigDec, slippage osmomath.BigDec, swapFee osmomath.BigDec, takerFeesFinal osmomath.BigDec, err error)
	SwapByDenom(ctx sdk.Context, msg *ammtypes.MsgSwapByDenom) (*ammtypes.MsgSwapByDenomResponse, error)
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
	SendCoins(goCtx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

type CommitmentKeeper interface {
	GetCommitments(ctx sdk.Context, creator sdk.AccAddress) commitmenttypes.Commitments
	CommitmentVestingInfo(goCtx context.Context, req *commitmenttypes.QueryCommitmentVestingInfoRequest) (*commitmenttypes.QueryCommitmentVestingInfoResponse, error)
	GetAllCommitmentsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]*commitmenttypes.Commitments, *query.PageResponse, error)
	CommitClaimedRewards(ctx sdk.Context, msg *commitmenttypes.MsgCommitClaimedRewards) (*commitmenttypes.MsgCommitClaimedRewardsResponse, error)
	UncommitTokens(ctx sdk.Context, creator sdk.AccAddress, denom string, amount sdkmath.Int, isLiquidation bool) error
	ProcessTokenVesting(ctx sdk.Context, denom string, amount sdkmath.Int, creator sdk.AccAddress) error
	CancelVest(ctx sdk.Context, msg *commitmenttypes.MsgCancelVest) (*commitmenttypes.MsgCancelVestResponse, error)
	ClaimVesting(ctx sdk.Context, msg *commitmenttypes.MsgClaimVesting) (*commitmenttypes.MsgClaimVestingResponse, error)
	CommitLiquidTokens(ctx sdk.Context, addr sdk.AccAddress, denom string, amount sdkmath.Int, lockUntil uint64) error
	GetParams(ctx sdk.Context) (params commitmenttypes.Params)
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

type TierKeeper interface {
	GetMembershipTier(ctx sdk.Context, user sdk.AccAddress) (total_portfolio sdkmath.LegacyDec, tier string, discount sdkmath.LegacyDec)
	CalculateUSDValue(ctx sdk.Context, denom string, amount sdkmath.Int) sdkmath.LegacyDec
}

type ParameterKeeper interface {
	GetParams(ctx sdk.Context) (params parametertypes.Params)
}

// MasterchefKeeper defines expected interface for masterchef keeper
type MasterchefKeeper interface {
	ClaimRewards(ctx sdk.Context, sender sdk.AccAddress, poolIds []uint64, recipient sdk.AccAddress) error
	UserPoolPendingReward(ctx sdk.Context, user sdk.AccAddress, poolId uint64) sdk.Coins
	GetParams(ctx sdk.Context) (params mastercheftypes.Params)
}

type AssetProfileKeeper interface {
	GetEntry(ctx sdk.Context, denom string) (atypes.Entry, bool)
	SetEntry(ctx sdk.Context, entry atypes.Entry)
	GetUsdcDenom(ctx sdk.Context) (string, bool)
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

type OracleKeeper interface {
	GetDenomPrice(ctx sdk.Context, denom string) osmomath.BigDec
}

type PerpetualKeeper interface {
	Open(ctx sdk.Context, msg *perpetualtypes.MsgOpen) (*perpetualtypes.MsgOpenResponse, error)
	Close(ctx sdk.Context, msg *perpetualtypes.MsgClose) (*perpetualtypes.MsgCloseResponse, error)
	AddCollateral(ctx sdk.Context, msg *perpetualtypes.MsgAddCollateral) (*perpetualtypes.MsgAddCollateralResponse, error)
	UpdateStopLoss(ctx sdk.Context, msg *perpetualtypes.MsgUpdateStopLoss) (*perpetualtypes.MsgUpdateStopLossResponse, error)
	UpdateTakeProfitPrice(ctx sdk.Context, msg *perpetualtypes.MsgUpdateTakeProfitPrice) (*perpetualtypes.MsgUpdateTakeProfitPriceResponse, error)
}
