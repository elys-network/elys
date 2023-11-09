package types

import (
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
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
	moduleQueriers   []ModuleQuerier
	ammKeeper        *ammkeeper.Keeper
	oracleKeeper     *oraclekeeper.Keeper
	bankKeeper       *bankkeeper.BaseKeeper
	stakingKeeper    *stakingkeeper.Keeper
	commitmentKeeper *commitmentkeeper.Keeper
	marginKeeper     *marginkeeper.Keeper
	incentiveKeeper  *incentivekeeper.Keeper
}

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
	}
}

type ElysQuery struct {
	PriceAll          *oracletypes.QueryAllPriceRequest       `json:"price_all,omitempty"`
	SwapEstimation    *ammtypes.QuerySwapEstimationRequest    `json:"query_swap_estimation,omitempty"`
	AssetInfo         *oracletypes.QueryGetAssetInfoRequest   `json:"asset_info,omitempty"`
	BalanceOfDenom    *ammtypes.QueryBalanceRequest           `json:"balance_of_denom,omitempty"`
	Params            *ammtypes.QueryParamsRequest            `json:"params,omitempty"`
	Pool              *ammtypes.QueryGetPoolRequest           `json:"pool,omitempty"`
	PoolAll           *ammtypes.QueryAllPoolRequest           `json:"pool_all,omitempty"`
	DenomLiquidity    *ammtypes.QueryGetDenomLiquidityRequest `json:"denom_liquidity,omitempty"`
	DenomLiquidityAll *ammtypes.QueryAllDenomLiquidityRequest `json:"denom_liquidity_all,omitempty"`
	SlippageTrack     *ammtypes.QuerySlippageTrackRequest     `json:"slippage_track,omitempty"`
	SlippageTrackAll  *ammtypes.QuerySlippageTrackAllRequest  `json:"slippage_track_all,omitempty"`
}

type CustomMessenger struct {
	wrapped          wasmkeeper.Messenger
	moduleMessengers []ModuleMessenger
	amm              *ammkeeper.Keeper
	margin           *marginkeeper.Keeper
	staking          *stakingkeeper.Keeper
	commitment       *commitmentkeeper.Keeper
	incentive        *incentivekeeper.Keeper
}

type ElysMsg struct {
	MsgSwapExactAmountIn           *ammtypes.MsgSwapExactAmountIn                 `json:"msg_swap_exact_amount_in,omitempty"`
	MsgOpen                        *margintypes.MsgOpen                           `json:"msg_open,omitempty"`
	MsgClose                       *margintypes.MsgClose                          `json:"msg_close,omitempty"`
	MsgStake                       *commitmenttypes.MsgStake                      `json:"msg_stake,omitempty"`
	MsgUnstake                     *commitmenttypes.MsgUnstake                    `json:"msg_unstake,omitempty"`
	MsgBeginRedelegate             *stakingtypes.MsgBeginRedelegate               `json:"msg_begin_redelegate,omitempty"`
	MsgCancelUnbondingDelegation   *stakingtypes.MsgCancelUnbondingDelegation     `json:"msg_cancel_unbonding_delegation"`
	MsgVest                        *commitmenttypes.MsgVest                       `json:"msg_vest"`
	MsgCancelVest                  *commitmenttypes.MsgCancelVest                 `json:"msg_cancel_vest"`
	MsgWithdrawRewards             *incentivetypes.MsgWithdrawRewards             `json:"msg_withdraw_rewards"`
	MsgWithdrawValidatorCommission *incentivetypes.MsgWithdrawValidatorCommission `json:"msg_withdraw_validator_commission"`
}

type RequestResponse struct {
	Code   uint64 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}
