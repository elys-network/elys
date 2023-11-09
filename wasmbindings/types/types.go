package types

import (
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ammclientwasmtypes "github.com/elys-network/elys/x/amm/client/wasm/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	commitmentclientwasmtypes "github.com/elys-network/elys/x/commitment/client/wasm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	incentiveclientwasmtypes "github.com/elys-network/elys/x/incentive/client/wasm/types"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	marginclientwasmtypes "github.com/elys-network/elys/x/margin/client/wasm/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	oracleclientwasmtypes "github.com/elys-network/elys/x/oracle/client/wasm/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
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
	PriceAll            *oracleclientwasmtypes.PriceAll                `json:"price_all,omitempty"`
	QuerySwapEstimation *ammclientwasmtypes.QuerySwapEstimationRequest `json:"query_swap_estimation,omitempty"`
	AssetInfo           *oracleclientwasmtypes.AssetInfo               `json:"asset_info,omitempty"`
	BalanceOfDenom      *ammclientwasmtypes.QueryBalanceRequest        `json:"balance_of_denom,omitempty"`
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
	MsgSwapExactAmountIn           *ammclientwasmtypes.MsgSwapExactAmountIn                 `json:"msg_swap_exact_amount_in,omitempty"`
	MsgOpen                        *marginclientwasmtypes.MsgOpen                           `json:"msg_open,omitempty"`
	MsgClose                       *marginclientwasmtypes.MsgClose                          `json:"msg_close,omitempty"`
	MsgStake                       *commitmentclientwasmtypes.MsgStake                      `json:"msg_stake,omitempty"`
	MsgUnstake                     *commitmentclientwasmtypes.MsgUnstake                    `json:"msg_unstake,omitempty"`
	MsgBeginRedelegate             *incentiveclientwasmtypes.MsgBeginRedelegate             `json:"msg_begin_redelegate,omitempty"`
	MsgCancelUnbondingDelegation   *incentiveclientwasmtypes.MsgCancelUnbondingDelegation   `json:"msg_cancel_unbonding_delegation"`
	MsgVest                        *incentiveclientwasmtypes.MsgVest                        `json:"msg_vest"`
	MsgCancelVest                  *incentiveclientwasmtypes.MsgCancelVest                  `json:"msg_cancel_vest"`
	MsgWithdrawRewards             *incentiveclientwasmtypes.MsgWithdrawRewards             `json:"msg_withdraw_rewards"`
	MsgWithdrawValidatorCommission *incentiveclientwasmtypes.MsgWithdrawValidatorCommission `json:"msg_withdraw_validator_commission"`
}
