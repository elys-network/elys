package types

import (
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	accountedpoolkeeper "github.com/elys-network/elys/x/accountedpool/keeper"
	accountedpooltypes "github.com/elys-network/elys/x/accountedpool/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofilekeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	burnerkeeper "github.com/elys-network/elys/x/burner/keeper"
	burnertypes "github.com/elys-network/elys/x/burner/types"
	clockkeeper "github.com/elys-network/elys/x/clock/keeper"
	clocktypes "github.com/elys-network/elys/x/clock/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	epochskeeper "github.com/elys-network/elys/x/epochs/keeper"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	parametertypes "github.com/elys-network/elys/x/parameter/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
	transferhookkeeper "github.com/elys-network/elys/x/transferhook/keeper"
	transferhooktypes "github.com/elys-network/elys/x/transferhook/types"
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
	moduleQueriers      []ModuleQuerier
	accountedpoolKeeper *accountedpoolkeeper.Keeper
	ammKeeper           *ammkeeper.Keeper
	assetprofileKeeper  *assetprofilekeeper.Keeper
	bankKeeper          *bankkeeper.BaseKeeper
	burnerKeeper        *burnerkeeper.Keeper
	clockKeeper         *clockkeeper.Keeper
	commitmentKeeper    *commitmentkeeper.Keeper
	epochsKeeper        *epochskeeper.Keeper
	incentiveKeeper     *incentivekeeper.Keeper
	leveragelpKeeper    *leveragelpkeeper.Keeper
	marginKeeper        *marginkeeper.Keeper
	oracleKeeper        *oraclekeeper.Keeper
	parameterKeeper     *parameterkeeper.Keeper
	stablestakeKeeper   *stablestakekeeper.Keeper
	stakingKeeper       *stakingkeeper.Keeper
	tokenomicsKeeper    *tokenomicskeeper.Keeper
	transferhookKeeper  *transferhookkeeper.Keeper
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
	// accountedpool queriers
	AccountedPoolParams           *accountedpooltypes.QueryParamsRequest           `json:"accounted_pool_params,omitempty"`
	AccountedPoolAccountedPool    *accountedpooltypes.QueryGetAccountedPoolRequest `json:"accounted_pool_accounted_pool,omitempty"`
	AccountedPoolAccountedPoolAll *accountedpooltypes.QueryAllAccountedPoolRequest `json:"accounted_pool_accounted_pool_all,omitempty"`
	PriceAll                      *PriceAll                                        `json:"price_all,omitempty"`
	QuerySwapEstimation           *QuerySwapEstimationRequest                      `json:"query_swap_estimation,omitempty"`
	AssetInfo                     *AssetInfo                                       `json:"asset_info,omitempty"`
	BalanceOfDenom                *QueryBalanceRequest                             `json:"balance_of_denom,omitempty"`
	Delegations                   *QueryDelegatorDelegationsRequest                `json:"delegations,omitempty"`
	UnbondingDelegations          *QueryDelegatorUnbondingDelegationsRequest       `json:"unbonding_delegations,omitempty"`
	StakedBalanceOfDenom          *QueryBalanceRequest                             `json:"staked_balance_of_denom,omitempty"`
	RewardsBalanceOfDenom         *QueryBalanceRequest                             `json:"rewards_balance_of_denom,omitempty"`
	ShowCommitments               *QueryCommitmentsRequest                         `json:"show_commitments,omitempty"`
	BalanceOfBorrow               *QueryBorrowRequest                              `json:"balance_of_borrow,omitempty"`
	Validators                    *QueryValidatorsRequest                          `json:"validators,omitempty"`

	// amm queriers
	AmmParams            *ammtypes.QueryParamsRequest            `json:"amm_params,omitempty"`
	AmmPool              *ammtypes.QueryGetPoolRequest           `json:"amm_pool,omitempty"`
	AmmPoolAll           *ammtypes.QueryAllPoolRequest           `json:"amm_pool_all,omitempty"`
	AmmDenomLiquidity    *ammtypes.QueryGetDenomLiquidityRequest `json:"amm_denom_liquidity,omitempty"`
	AmmDenomLiquidityAll *ammtypes.QueryAllDenomLiquidityRequest `json:"amm_denom_liquidity_all,omitempty"`
	AmmSwapEstimation    *ammtypes.QuerySwapEstimationRequest    `json:"amm_swap_estimation,omitempty"`
	AmmSlippageTrack     *ammtypes.QuerySlippageTrackRequest     `json:"amm_slippage_track,omitempty"`
	AmmSlippageTrackAll  *ammtypes.QuerySlippageTrackAllRequest  `json:"amm_slippage_track_all,omitempty"`
	AmmBalance           *ammtypes.QueryBalanceRequest           `json:"amm_balance,omitempty"`

	// assetprofile queriers
	AssetProfileParams   *assetprofiletypes.QueryParamsRequest   `json:"asset_profile_params,omitempty"`
	AssetProfileEntry    *assetprofiletypes.QueryGetEntryRequest `json:"asset_profile_entry,omitempty"`
	AssetProfileEntryAll *assetprofiletypes.QueryAllEntryRequest `json:"asset_profile_entry_all,omitempty"`

	// burner queriers
	BurnerParams     *burnertypes.QueryParamsRequest     `json:"burner_params,omitempty"`
	BurnerHistory    *burnertypes.QueryGetHistoryRequest `json:"burner_history,omitempty"`
	BurnerHistoryAll *burnertypes.QueryAllHistoryRequest `json:"burner_history_all,omitempty"`

	// clock queriers
	ClockClockContracts *clocktypes.QueryClockContracts `json:"clock_clock_contracts,omitempty"`
	ClockParams         *clocktypes.QueryParamsRequest  `json:"clock_params,omitempty"`

	// commitment queriers
	CommitmentParams                *commitmenttypes.QueryParamsRequest                        `json:"commitment_params,omitempty"`
	CommitmentShowCommitments       *commitmenttypes.QueryShowCommitmentsRequest               `json:"commitment_show_commitments,omitempty"`
	CommitmentDelegations           *commitmenttypes.QueryDelegatorDelegationsRequest          `json:"commitment_delegations,omitempty"`
	CommitmentUnbondingDelegations  *commitmenttypes.QueryDelegatorUnbondingDelegationsRequest `json:"commitment_unbonding_delegations,omitempty"`
	CommitmentStakedBalanceOfDenom  *ammtypes.QueryBalanceRequest                              `json:"commitment_staked_balance_of_denom,omitempty"`
	CommitmentRewardsBalanceOfDenom *ammtypes.QueryBalanceRequest                              `json:"commitment_rewards_balance_of_denom,omitempty"`
	CommitmentAllValidators         *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_all_validators,omitempty"`
	CommitmentDelegatorValidators   *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_delegator_validators,omitempty"`
	CommitmentStakedPositions       *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_staked_positions,omitempty"`
	CommitmentUnStakedPositions     *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_un_staked_positions,omitempty"`

	// epochs queriers
	EpochsEpochInfos   *epochstypes.QueryEpochsInfoRequest   `json:"epochs_epoch_infos,omitempty"`
	EpochsCurrentEpoch *epochstypes.QueryCurrentEpochRequest `json:"epochs_current_epoch,omitempty"`

	// incentive queriers
	IncentiveParams        *incentivetypes.QueryParamsRequest        `json:"incentive_params,omitempty"`
	IncentiveCommunityPool *incentivetypes.QueryCommunityPoolRequest `json:"incentive_community_pool,omitempty"`

	// leveragelp queriers
	LeveragelpParams                   *leveragelptypes.ParamsRequest              `json:"leveragelp_params,omitempty"`
	LeveragelpQueryPositions           *leveragelptypes.PositionsRequest           `json:"leveragelp_query_positions,omitempty"`
	LeveragelpQueryPositionsByPool     *leveragelptypes.PositionsByPoolRequest     `json:"leveragelp_query_positions_by_pool,omitempty"`
	LeveragelpGetStatus                *leveragelptypes.StatusRequest              `json:"leveragelp_get_status,omitempty"`
	LeveragelpQueryPositionsForAddress *leveragelptypes.PositionsForAddressRequest `json:"leveragelp_query_positions_for_address,omitempty"`
	LeveragelpGetWhitelist             *leveragelptypes.WhitelistRequest           `json:"leveragelp_get_whitelist,omitempty"`
	LeveragelpIsWhitelisted            *leveragelptypes.IsWhitelistedRequest       `json:"leveragelp_is_whitelisted,omitempty"`
	LeveragelpPool                     *leveragelptypes.QueryGetPoolRequest        `json:"leveragelp_pool,omitempty"`
	LeveragelpPools                    *leveragelptypes.QueryAllPoolRequest        `json:"leveragelp_pools,omitempty"`
	LeveragelpPosition                 *leveragelptypes.PositionRequest            `json:"leveragelp_position,omitempty"`

	// margin queriers
	MarginParams                 *margintypes.ParamsRequest              `json:"margin_params,omitempty"`
	MarginQueryPositions         *margintypes.PositionsRequest           `json:"margin_query_positions,omitempty"`
	MarginQueryPositionsByPool   *margintypes.PositionsByPoolRequest     `json:"margin_query_positions_by_pool,omitempty"`
	MarginGetStatus              *margintypes.StatusRequest              `json:"margin_get_status,omitempty"`
	MarginGetPositionsForAddress *margintypes.PositionsForAddressRequest `json:"margin_get_positions_for_address,omitempty"`
	MarginGetWhitelist           *margintypes.WhitelistRequest           `json:"margin_get_whitelist,omitempty"`
	MarginIsWhitelisted          *margintypes.IsWhitelistedRequest       `json:"margin_is_whitelisted,omitempty"`
	MarginPool                   *margintypes.QueryGetPoolRequest        `json:"margin_pool,omitempty"`
	MarginPools                  *margintypes.QueryAllPoolRequest        `json:"margin_pools,omitempty"`
	MarginMTP                    *margintypes.MTPRequest                 `json:"margin_mtp,omitempty"`

	// oracle queriers
	OracleParams            *oracletypes.QueryParamsRequest            `json:"oracle_params,omitempty"`
	OracleBandPriceResult   *oracletypes.QueryBandPriceRequest         `json:"oracle_band_price_result,omitempty"`
	OracleLastBandRequestId *oracletypes.QueryLastBandRequestIdRequest `json:"oracle_last_band_request_id,omitempty"`
	OracleAssetInfo         *oracletypes.QueryGetAssetInfoRequest      `json:"oracle_asset_info,omitempty"`
	OracleAssetInfoAll      *oracletypes.QueryAllAssetInfoRequest      `json:"oracle_asset_info_all,omitempty"`
	OraclePrice             *oracletypes.QueryGetPriceRequest          `json:"oracle_price,omitempty"`
	OraclePriceAll          *oracletypes.QueryAllPriceRequest          `json:"oracle_price_all,omitempty"`
	OraclePriceFeeder       *oracletypes.QueryGetPriceFeederRequest    `json:"oracle_price_feeder,omitempty"`
	OraclePriceFeederAll    *oracletypes.QueryAllPriceFeederRequest    `json:"oracle_price_feeder_all,omitempty"`

	// parameter queriers
	ParameterParams              *parametertypes.QueryParamsRequest              `json:"parameter_params,omitempty"`
	ParameterAnteHandlerParamAll *parametertypes.QueryAllAnteHandlerParamRequest `json:"parameter_ante_handler_param_all,omitempty"`

	// stablestake queriers
	StableStakeParams          *stablestaketypes.QueryParamsRequest      `json:"stable_stake_params,omitempty"`
	StableStakeBorrowRatio     *stablestaketypes.QueryBorrowRatioRequest `json:"stable_stake_borrow_ratio,omitempty"`
	StableStakeBalanceOfBorrow *commitmenttypes.QueryBorrowAmountRequest `json:"stable_stake_balance_of_borrow,omitempty"`

	// tokenomics queriers
	TokenomicsParams                *tokenomicstypes.QueryParamsRequest                `json:"tokenomics_params,omitempty"`
	TokenomicsAirdrop               *tokenomicstypes.QueryGetAirdropRequest            `json:"tokenomics_airdrop,omitempty"`
	TokenomicsAirdropAll            *tokenomicstypes.QueryAllAirdropRequest            `json:"tokenomics_airdrop_all,omitempty"`
	TokenomicsGenesisInflation      *tokenomicstypes.QueryGetGenesisInflationRequest   `json:"tokenomics_genesis_inflation,omitempty"`
	TokenomicsTimeBasedInflation    *tokenomicstypes.QueryGetTimeBasedInflationRequest `json:"tokenomics_time_based_inflation,omitempty"`
	TokenomicsTimeBasedInflationAll *tokenomicstypes.QueryAllTimeBasedInflationRequest `json:"tokenomics_time_based_inflation_all,omitempty"`

	// transferhook queriers
	TransferHookParams *transferhooktypes.QueryParamsRequest `json:"transfer_hook_params,omitempty"`
}

type QueryCommitmentsRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

type QueryBorrowRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

type QueryDelegatorDelegationsRequest struct {
	// delegator_addr defines the delegator address to query for.
	DelegatorAddr string `protobuf:"bytes,1,opt,name=delegator_addr,json=delegatorAddr,proto3" json:"delegator_addr,omitempty"`
}

// QueryDelegatorDelegationsResponse is response type for the
// Query/DelegatorDelegations RPC method.
type QueryDelegatorDelegationsResponse struct {
	// delegation_responses defines all the delegations' info of a delegator.
	DelegationResponses []stakingtypes.DelegationResponse `protobuf:"bytes,1,rep,name=delegation_responses,json=delegationResponses,proto3" json:"delegation_responses"`
}

// QueryDelegatorUnbondingDelegationsRequest is request type for the
type QueryDelegatorUnbondingDelegationsRequest struct {
	// delegator_addr defines the delegator address to query for.
	DelegatorAddr string `protobuf:"bytes,1,opt,name=delegator_addr,json=delegatorAddr,proto3" json:"delegator_addr,omitempty"`
}

// QueryUnbondingDelegatorDelegationsResponse is response type for the
// Query/UnbondingDelegatorDelegations RPC method.
type QueryDelegatorUnbondingDelegationsResponse struct {
	UnbondingResponses []stakingtypes.UnbondingDelegation `protobuf:"bytes,1,rep,name=unbonding_responses,json=unbondingResponses,proto3" json:"unbonding_responses"`
}

// QueryValidatorsRequest is request type for Query/Validators RPC method.
type QueryValidatorsRequest struct {
	// status enables to query for validators matching a given status.
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

type QueryShowCommitmentsResponse struct {
	Commitments *commitmenttypes.Commitments `protobuf:"bytes,1,opt,name=commitments,proto3" json:"commitments,omitempty"`
}

type CustomMessenger struct {
	wrapped          wasmkeeper.Messenger
	moduleMessengers []ModuleMessenger
	accountedpool    *accountedpoolkeeper.Keeper
	amm              *ammkeeper.Keeper
	assetprofile     *assetprofilekeeper.Keeper
	bank             *bankkeeper.BaseKeeper
	burner           *burnerkeeper.Keeper
	clock            *clockkeeper.Keeper
	commitment       *commitmentkeeper.Keeper
	epochs           *epochskeeper.Keeper
	incentive        *incentivekeeper.Keeper
	leveragelp       *leveragelpkeeper.Keeper
	margin           *marginkeeper.Keeper
	oracle           *oraclekeeper.Keeper
	parameter        *parameterkeeper.Keeper
	stablestake      *stablestakekeeper.Keeper
	staking          *stakingkeeper.Keeper
	tokenomics       *tokenomicskeeper.Keeper
	transferhook     *transferhookkeeper.Keeper
}

type ElysMsg struct {
	// accountedpool messages

	// amm messages
	AmmCreatePool                    *ammtypes.MsgCreatePool                    `json:"amm_create_pool,omitempty"`
	AmmJoinPool                      *ammtypes.MsgJoinPool                      `json:"amm_join_pool,omitempty"`
	AmmExitPool                      *ammtypes.MsgExitPool                      `json:"amm_exit_pool,omitempty"`
	AmmSwapExactAmountIn             *ammtypes.MsgSwapExactAmountIn             `json:"amm_swap_exact_amount_in,omitempty"`
	AmmSwapExactAmountOut            *ammtypes.MsgSwapExactAmountOut            `json:"amm_swap_exact_amount_out,omitempty"`
	AmmFeedMultipleExternalLiquidity *ammtypes.MsgFeedMultipleExternalLiquidity `json:"amm_feed_multiple_external_liquidity,omitempty"`

	// assetprofile messages
	AssetProfileCreateEntry *assetprofiletypes.MsgCreateEntry `json:"asset_profile_create_entry,omitempty"`
	AssetProfileUpdateEntry *assetprofiletypes.MsgUpdateEntry `json:"asset_profile_update_entry,omitempty"`
	AssetProfileDeleteEntry *assetprofiletypes.MsgDeleteEntry `json:"asset_profile_delete_entry,omitempty"`

	// burner messages

	// clock messages
	ClockUpdateParams *clocktypes.MsgUpdateParams `json:"clock_update_params,omitempty"`

	// commitment messages
	CommitmentCommitLiquidTokens     *commitmenttypes.MsgCommitLiquidTokens   `json:"commitment_commit_liquid_tokens,omitempty"`
	CommitmentCommitUnclaimedRewards *commitmenttypes.MsgCommitClaimedRewards `json:"commitment_commit_unclaimed_rewards,omitempty"`
	CommitmentUncommitTokens         *commitmenttypes.MsgUncommitTokens       `json:"commitment_uncommit_tokens,omitempty"`
	CommitmentVest                   *commitmenttypes.MsgVest                 `json:"commitment_vest"`
	CommitmentVestNow                *commitmenttypes.MsgVestNow              `json:"commitment_vest_now"`
	CommitmentVestLiquid             *commitmenttypes.MsgVestLiquid           `json:"commitment_vest_liquid"`
	CommitmentCancelVest             *commitmenttypes.MsgCancelVest           `json:"commitment_cancel_vest"`
	CommitmentUpdateVestingInfo      *commitmenttypes.MsgUpdateVestingInfo    `json:"commitment_update_vesting_info"`
	CommitmentStake                  *commitmenttypes.MsgStake                `json:"commitment_stake,omitempty"`
	CommitmentUnstake                *commitmenttypes.MsgUnstake              `json:"commitment_unstake,omitempty"`

	// epochs messages

	// incentive messages
	IncentiveBeginRedelegate             *stakingtypes.MsgBeginRedelegate               `json:"incentive_begin_redelegate,omitempty"`
	IncentiveCancelUnbondingDelegation   *stakingtypes.MsgCancelUnbondingDelegation     `json:"incentive_cancel_unbonding_delegation"`
	IncentiveWithdrawRewards             *incentivetypes.MsgWithdrawRewards             `json:"incentive_withdraw_rewards"`
	IncentiveWithdrawValidatorCommission *incentivetypes.MsgWithdrawValidatorCommission `json:"incentive_withdraw_validator_commission"`

	// leveragelp messages
	LeveragelpOpen         *leveragelptypes.MsgOpen         `json:"leveragelp_open,omitempty"`
	LeveragelpClose        *leveragelptypes.MsgClose        `json:"leveragelp_close,omitempty"`
	LeveragelpUpdateParams *leveragelptypes.MsgUpdateParams `json:"leveragelp_update_params,omitempty"`
	LeveragelpUpdatePools  *leveragelptypes.MsgUpdatePools  `json:"leveragelp_update_pools,omitempty"`
	LeveragelpWhitelist    *leveragelptypes.MsgWhitelist    `json:"leveragelp_whitelist,omitempty"`
	LeveragelpDewhitelist  *leveragelptypes.MsgDewhitelist  `json:"leveragelp_dewhitelist,omitempty"`

	// margin messages
	MarginOpen         *margintypes.MsgOpen         `json:"margin_open,omitempty"`
	MarginClose        *margintypes.MsgClose        `json:"margin_close,omitempty"`
	MarginUpdateParams *margintypes.MsgUpdateParams `json:"margin_update_params,omitempty"`
	MarginUpdatePools  *margintypes.MsgUpdatePools  `json:"margin_update_pools,omitempty"`
	MarginWhitelist    *margintypes.MsgWhitelist    `json:"margin_whitelist,omitempty"`
	MarginDewhitelist  *margintypes.MsgDewhitelist  `json:"margin_dewhitelist,omitempty"`

	// oracle messages
	OracleFeedPrice          *oracletypes.MsgFeedPrice          `json:"oracle_feed_price,omitempty"`
	OracleFeedMultiplePrices *oracletypes.MsgFeedMultiplePrices `json:"oracle_feed_multiple_price,omitempty"`
	OracleRequestBandPrice   *oracletypes.MsgRequestBandPrice   `json:"oracle_request_band_price,omitempty"`
	OracleSetPriceFeeder     *oracletypes.MsgSetPriceFeeder     `json:"oracle_set_price_feeder,omitempty"`
	OracleDeletePriceFeeder  *oracletypes.MsgDeletePriceFeeder  `json:"oracle_delete_price_feeder,omitempty"`

	// parameter messages

	// stablestake messages
	StakestakeBond   *stablestaketypes.MsgBond   `json:"stablestake_bond,omitempty"`
	StakestakeUnbond *stablestaketypes.MsgUnbond `json:"stablestake_unbond,omitempty"`

	// tokenomics messages
	TokenomicsCreateAirdrop            *tokenomicstypes.MsgCreateAirdrop            `json:"tokenomics_create_airdrop,omitempty"`
	TokenomicsUpdateAirdrop            *tokenomicstypes.MsgUpdateAirdrop            `json:"tokenomics_update_airdrop,omitempty"`
	TokenomicsDeleteAirdrop            *tokenomicstypes.MsgDeleteAirdrop            `json:"tokenomics_delete_airdrop,omitempty"`
	TokenomicsUpdateGenesisInflation   *tokenomicstypes.MsgUpdateGenesisInflation   `json:"tokenomics_update_genesis_inflation,omitempty"`
	TokenomicsCreateTimeBasedInflation *tokenomicstypes.MsgCreateTimeBasedInflation `json:"tokenomics_create_time_based_inflation,omitempty"`
	TokenomicsUpdateTimeBasedInflation *tokenomicstypes.MsgUpdateTimeBasedInflation `json:"tokenomics_update_time_based_inflation,omitempty"`
	TokenomicsDeleteTimeBasedInflation *tokenomicstypes.MsgDeleteTimeBasedInflation `json:"tokenomics_delete_time_based_inflation,omitempty"`

	// transferhook messages
}

type RequestResponse struct {
	Code   uint64 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}
