package types

import (
	"errors"

	cosmos_sdk_math "cosmossdk.io/math"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtype "github.com/elys-network/elys/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	incentiveclientwasmtypes "github.com/elys-network/elys/x/incentive/client/wasm/types"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
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
	PriceAll            *PriceAll                   `json:"price_all,omitempty"`
	QuerySwapEstimation *QuerySwapEstimationRequest `json:"query_swap_estimation,omitempty"`
	AssetInfo           *AssetInfo                  `json:"asset_info,omitempty"`
	BalanceOfDenom      *QueryBalanceRequest        `json:"balance_of_denom,omitempty"`
}

type PriceAll struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type AllPriceResponse struct {
	Price      []oracletypes.Price `protobuf:"bytes,1,rep,name=price,proto3" json:"price"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type QuerySwapEstimationRequest struct {
	TokenIn sdk.Coin                    `protobuf:"bytes,2,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	Routes  []ammtype.SwapAmountInRoute `protobuf:"bytes,1,rep,name=routes,proto3" json:"routes,omitempty"`
}

type QuerySwapEstimationResponse struct {
	SpotPrice sdk.Dec  `protobuf:"bytes,1,opt,name=SpotPrice,proto3" json:"spot_price,omitempty"`
	TokenOut  sdk.Coin `protobuf:"bytes,2,opt,name=tokenOut,proto3" json:"token_out,omitempty"`
}

type AssetInfo struct {
	Denom string `protobuf:"bytes,1,opt,name=Denom,proto3" json:"Denom,omitempty"`
}

type AssetInfoResponse struct {
	AssetInfo *AssetInfoType `protobuf:"bytes,1,opt,name=AssetInfo,proto3" json:"asset_info,omitempty"`
}

type AssetInfoType struct {
	Denom      string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Display    string `protobuf:"bytes,2,opt,name=display,proto3" json:"display,omitempty"`
	BandTicker string `protobuf:"bytes,3,opt,name=bandTicker,proto3" json:"band_ticker,omitempty"`
	ElysTicker string `protobuf:"bytes,4,opt,name=elysTicker,proto3" json:"elys_ticker,omitempty"`
	Decimal    uint64 `protobuf:"varint,5,opt,name=decimal,proto3" json:"decimal,omitempty"`
}

type QueryBalanceRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Denom   string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
}

type QueryBalanceResponse struct {
	Balance sdk.Coin `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance,omitempty"`
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
	MsgSwapExactAmountIn           *MsgSwapExactAmountIn                                    `json:"msg_swap_exact_amount_in,omitempty"`
	MsgOpen                        *MsgOpen                                                 `json:"msg_open,omitempty"`
	MsgClose                       *MsgClose                                                `json:"msg_close,omitempty"`
	MsgStake                       *MsgStake                                                `json:"msg_stake,omitempty"`
	MsgUnstake                     *MsgUnstake                                              `json:"msg_unstake,omitempty"`
	MsgBeginRedelegate             *incentiveclientwasmtypes.MsgBeginRedelegate             `json:"msg_begin_redelegate,omitempty"`
	MsgCancelUnbondingDelegation   *incentiveclientwasmtypes.MsgCancelUnbondingDelegation   `json:"msg_cancel_unbonding_delegation"`
	MsgVest                        *incentiveclientwasmtypes.MsgVest                        `json:"msg_vest"`
	MsgCancelVest                  *incentiveclientwasmtypes.MsgCancelVest                  `json:"msg_cancel_vest"`
	MsgWithdrawRewards             *incentiveclientwasmtypes.MsgWithdrawRewards             `json:"msg_withdraw_rewards"`
	MsgWithdrawValidatorCommission *incentiveclientwasmtypes.MsgWithdrawValidatorCommission `json:"msg_withdraw_validator_commission"`
}

type MsgSwapExactAmountIn struct {
	Sender            string                      `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Routes            []ammtype.SwapAmountInRoute `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	TokenIn           sdk.Coin                    `protobuf:"bytes,3,opt,name=tokenIn,proto3" json:"token_in,omitempty"`
	TokenOutMinAmount cosmos_sdk_math.Int         `protobuf:"bytes,4,opt,name=tokenOutMinAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_min_amount,omitempty"`
	MetaData          *[]byte                     `protobuf:"bytes,5,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgSwapExactAmountInResponse struct {
	TokenOutAmount cosmos_sdk_math.Int `protobuf:"bytes,1,opt,name=tokenOutAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_out_amount,omitempty"`
	MetaData       *[]byte             `protobuf:"bytes,2,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgOpen struct {
	Creator          string               `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	CollateralAsset  string               `protobuf:"bytes,2,opt,name=collateralAsset,proto3" json:"collateral_asset,omitempty"`
	CollateralAmount sdk.Uint             `protobuf:"bytes,3,opt,name=collateralAmount,proto3" json:"collateral_amount,omitempty"`
	BorrowAsset      string               `protobuf:"bytes,4,opt,name=borrowAsset,proto3" json:"borrow_asset,omitempty"`
	Position         margintypes.Position `protobuf:"bytes,5,opt,name=position,proto3" json:"position,omitempty"`
	Leverage         sdk.Dec              `protobuf:"bytes,6,opt,name=leverage,proto3" json:"leverage,omitempty"`
	TakeProfitPrice  sdk.Dec              `protobuf:"bytes,7,opt,name=takeProfitPrice,proto3" json:"take_profit_price,omitempty"`
	MetaData         *[]byte              `protobuf:"bytes,8,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgClose struct {
	Creator  string  `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Id       int64   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	MetaData *[]byte `protobuf:"bytes,3,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgOpenResponse struct {
	MetaData *[]byte `protobuf:"bytes,1,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}
type MsgCloseResponse struct {
	MetaData *[]byte `protobuf:"bytes,1,opt,name=tokenData,proto3" json:"meta_data,omitempty"`
}

type MsgStake struct {
	Address          string              `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Amount           cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Asset            string              `protobuf:"bytes,3,opt,name=asset,proto3" json:"asset,omitempty"`
	ValidatorAddress string              `protobuf:"bytes,4,opt,name=validator_address,proto3" json:"validator_address,omitempty"`
}

type MsgUnstake struct {
	Address          string              `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Amount           cosmos_sdk_math.Int `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Asset            string              `protobuf:"bytes,3,opt,name=asset,proto3" json:"asset,omitempty"`
	ValidatorAddress string              `protobuf:"bytes,4,opt,name=validator_address,proto3" json:"validator_address,omitempty"`
}

type RequestResponse struct {
	Code   uint64 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}
