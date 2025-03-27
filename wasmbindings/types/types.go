package types

import (
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// ModuleQuerier is an interface that all module queriers should implement.
type ModuleQuerier interface {
	// HandleQuery processes the query and returns an error if it cannot handle it.
	HandleQuery(ctx sdk.Context, query ElysQuery) ([]byte, error)
}

// ModuleMsgHandler is an interface that all module messenger should implement.
type ModuleMessenger interface {
	// HandleMsg processes the message and returns an error if it cannot handle it.
	HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg ElysMsg) ([]sdk.Event, [][]byte, error)
}

// ErrCannotHandleQuery is returned by a ModuleQuerier when it cannot handle a query.
var ErrCannotHandleQuery = errors.New("cannot handle query")

// ErrCannotHandleMsg is returned by a ModuleMsgHandler when it cannot handle a message.
var ErrCannotHandleMsg = errors.New("cannot handle message")

type QueryPlugin struct {
	moduleQueriers []ModuleQuerier
	ammKeeper      *ammkeeper.Keeper
	authKeeper     *authkeeper.AccountKeeper
	bankKeeper     *bankkeeper.BaseKeeper
	stakingKeeper  *stakingkeeper.Keeper
}

type ElysQuery struct {
	// amm queriers
	AmmParams  *ammtypes.QueryParamsRequest  `json:"amm_params,omitempty"`
	AmmPool    *ammtypes.QueryGetPoolRequest `json:"amm_pool,omitempty"`
	AmmPoolAll *ammtypes.QueryAllPoolRequest `json:"amm_pool_all,omitempty"`
	// AmmEarnMiningPoolAll     *ammtypes.QueryEarnPoolRequest              `json:"amm_earn_mining_pool_all,omitempty"`
	AmmDenomLiquidity               *ammtypes.QueryGetDenomLiquidityRequest            `json:"amm_denom_liquidity,omitempty"`
	AmmDenomLiquidityAll            *ammtypes.QueryAllDenomLiquidityRequest            `json:"amm_denom_liquidity_all,omitempty"`
	AmmSwapEstimationExactAmountOut *ammtypes.QuerySwapEstimationExactAmountOutRequest `json:"amm_swap_estimation_exact_amount_out,omitempty"`
	AmmSwapEstimation               *ammtypes.QuerySwapEstimationRequest               `json:"amm_swap_estimation,omitempty"`
	AmmSwapEstimationByDenom        *ammtypes.QuerySwapEstimationByDenomRequest        `json:"amm_swap_estimation_by_denom,omitempty"`
	AmmJoinPoolEstimation           *ammtypes.QueryJoinPoolEstimationRequest           `json:"amm_join_pool_estimation,omitempty"`
	AmmExitPoolEstimation           *ammtypes.QueryExitPoolEstimationRequest           `json:"amm_exit_pool_estimation,omitempty"`
	AmmSlippageTrack                *ammtypes.QuerySlippageTrackRequest                `json:"amm_slippage_track,omitempty"`
	AmmSlippageTrackAll             *ammtypes.QuerySlippageTrackAllRequest             `json:"amm_slippage_track_all,omitempty"`
	AmmBalance                      *ammtypes.QueryBalanceRequest                      `json:"amm_balance,omitempty"`
	AmmInRouteByDenom               *ammtypes.QueryInRouteByDenomRequest               `json:"amm_in_route_by_denom,omitempty"`
	AmmOutRouteByDenom              *ammtypes.QueryOutRouteByDenomRequest              `json:"amm_out_route_by_denom,omitempty"`

	// auth queriers
	AuthAddresses *authtypes.QueryAccountsRequest `json:"auth_addresses,omitempty"`
}

type CustomMessenger struct {
	wrapped          wasmkeeper.Messenger
	moduleMessengers []ModuleMessenger
	amm              *ammkeeper.Keeper
	auth             *authkeeper.AccountKeeper
	bank             *bankkeeper.BaseKeeper
	staking          *stakingkeeper.Keeper
}

type ElysMsg struct {
	// amm messages
	AmmCreatePool               *ammtypes.MsgCreatePool               `json:"amm_create_pool,omitempty"`
	AmmJoinPool                 *ammtypes.MsgJoinPool                 `json:"amm_join_pool,omitempty"`
	AmmExitPool                 *ammtypes.MsgExitPool                 `json:"amm_exit_pool,omitempty"`
	AmmUpFrontSwapExactAmountIn *ammtypes.MsgUpFrontSwapExactAmountIn `json:"amm_upfront_swap_exact_amount_in,omitempty"`
	AmmSwapExactAmountIn        *ammtypes.MsgSwapExactAmountIn        `json:"amm_swap_exact_amount_in,omitempty"`
	AmmSwapExactAmountOut       *ammtypes.MsgSwapExactAmountOut       `json:"amm_swap_exact_amount_out,omitempty"`
	AmmSwapByDenom              *ammtypes.MsgSwapByDenom              `json:"amm_swap_by_denom,omitempty"`
}

type RequestResponse struct {
	Code   uint64 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}

type AuthAddressesResponse struct {
	// addresses are the existing accounts’ addresses
	Addresses []string `json:"addresses"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
