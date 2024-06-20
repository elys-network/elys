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
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"
	estakingtypes "github.com/elys-network/elys/x/estaking/types"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
	parametertypes "github.com/elys-network/elys/x/parameter/types"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	tierkeeper "github.com/elys-network/elys/x/tier/keeper"
	tiertypes "github.com/elys-network/elys/x/tier/types"
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
	authKeeper          *authkeeper.AccountKeeper
	bankKeeper          *bankkeeper.BaseKeeper
	burnerKeeper        *burnerkeeper.Keeper
	clockKeeper         *clockkeeper.Keeper
	commitmentKeeper    *commitmentkeeper.Keeper
	epochsKeeper        *epochskeeper.Keeper
	incentiveKeeper     *incentivekeeper.Keeper
	masterchefKeeper    *masterchefkeeper.Keeper
	leveragelpKeeper    *leveragelpkeeper.Keeper
	perpetualKeeper     *perpetualkeeper.Keeper
	oracleKeeper        *oraclekeeper.Keeper
	parameterKeeper     *parameterkeeper.Keeper
	stablestakeKeeper   *stablestakekeeper.Keeper
	stakingKeeper       *stakingkeeper.Keeper
	tokenomicsKeeper    *tokenomicskeeper.Keeper
	transferhookKeeper  *transferhookkeeper.Keeper
	estakingKeeper      *estakingkeeper.Keeper
	tierKeeper          *tierkeeper.Keeper
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
	AccountedPoolAccountedPool    *accountedpooltypes.QueryGetAccountedPoolRequest `json:"accounted_pool_accounted_pool,omitempty"`
	AccountedPoolAccountedPoolAll *accountedpooltypes.QueryAllAccountedPoolRequest `json:"accounted_pool_accounted_pool_all,omitempty"`

	// amm queriers
	AmmParams                *ammtypes.QueryParamsRequest                `json:"amm_params,omitempty"`
	AmmPool                  *ammtypes.QueryGetPoolRequest               `json:"amm_pool,omitempty"`
	AmmPoolAll               *ammtypes.QueryAllPoolRequest               `json:"amm_pool_all,omitempty"`
	AmmEarnMiningPoolAll     *ammtypes.QueryEarnPoolRequest              `json:"amm_earn_mining_pool_all,omitempty"`
	AmmDenomLiquidity        *ammtypes.QueryGetDenomLiquidityRequest     `json:"amm_denom_liquidity,omitempty"`
	AmmDenomLiquidityAll     *ammtypes.QueryAllDenomLiquidityRequest     `json:"amm_denom_liquidity_all,omitempty"`
	AmmSwapEstimation        *ammtypes.QuerySwapEstimationRequest        `json:"amm_swap_estimation,omitempty"`
	AmmSwapEstimationByDenom *ammtypes.QuerySwapEstimationByDenomRequest `json:"amm_swap_estimation_by_denom,omitempty"`
	AmmJoinPoolEstimation    *ammtypes.QueryJoinPoolEstimationRequest    `json:"amm_join_pool_estimation,omitempty"`
	AmmExitPoolEstimation    *ammtypes.QueryExitPoolEstimationRequest    `json:"amm_exit_pool_estimation,omitempty"`
	AmmSlippageTrack         *ammtypes.QuerySlippageTrackRequest         `json:"amm_slippage_track,omitempty"`
	AmmSlippageTrackAll      *ammtypes.QuerySlippageTrackAllRequest      `json:"amm_slippage_track_all,omitempty"`
	AmmBalance               *ammtypes.QueryBalanceRequest               `json:"amm_balance,omitempty"`
	AmmInRouteByDenom        *ammtypes.QueryInRouteByDenomRequest        `json:"amm_in_route_by_denom,omitempty"`
	AmmOutRouteByDenom       *ammtypes.QueryOutRouteByDenomRequest       `json:"amm_out_route_by_denom,omitempty"`
	AmmPriceByDenom          *ammtypes.QueryAMMPriceRequest              `json:"amm_price_by_denom,omitempty"`
	// assetprofile queriers
	AssetProfileParams       *assetprofiletypes.QueryParamsRequest          `json:"asset_profile_params,omitempty"`
	AssetProfileEntry        *assetprofiletypes.QueryGetEntryRequest        `json:"asset_profile_entry,omitempty"`
	AssetProfileEntryByDenom *assetprofiletypes.QueryGetEntryByDenomRequest `json:"asset_profile_entry_by_denom,omitempty"`
	AssetProfileEntryAll     *assetprofiletypes.QueryAllEntryRequest        `json:"asset_profile_entry_all,omitempty"`

	// auth queriers
	AuthAddresses *authtypes.QueryAccountsRequest `json:"auth_addresses,omitempty"`

	// burner queriers
	BurnerParams     *burnertypes.QueryParamsRequest     `json:"burner_params,omitempty"`
	BurnerHistory    *burnertypes.QueryGetHistoryRequest `json:"burner_history,omitempty"`
	BurnerHistoryAll *burnertypes.QueryAllHistoryRequest `json:"burner_history_all,omitempty"`

	// clock queriers
	ClockClockContracts *clocktypes.QueryClockContracts `json:"clock_clock_contracts,omitempty"`
	ClockParams         *clocktypes.QueryParamsRequest  `json:"clock_params,omitempty"`

	// commitment queriers
	CommitmentParams               *commitmenttypes.QueryParamsRequest                        `json:"commitment_params,omitempty"`
	CommitmentShowCommitments      *commitmenttypes.QueryShowCommitmentsRequest               `json:"commitment_show_commitments,omitempty"`
	CommitmentDelegations          *commitmenttypes.QueryDelegatorDelegationsRequest          `json:"commitment_delegations,omitempty"`
	CommitmentUnbondingDelegations *commitmenttypes.QueryDelegatorUnbondingDelegationsRequest `json:"commitment_unbonding_delegations,omitempty"`
	CommitmentStakedBalanceOfDenom *ammtypes.QueryBalanceRequest                              `json:"commitment_staked_balance_of_denom,omitempty"`
	CommitmentAllValidators        *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_all_validators,omitempty"`
	CommitmentDelegatorValidators  *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_delegator_validators,omitempty"`
	CommitmentStakedPositions      *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_staked_positions,omitempty"`
	CommitmentUnStakedPositions    *commitmenttypes.QueryValidatorsRequest                    `json:"commitment_un_staked_positions,omitempty"`
	CommitmentVestingInfo          *commitmenttypes.QueryVestingInfoRequest                   `json:"commitment_vesting_info,omitempty"`
	CommitmentNumberOfCommitments  *commitmenttypes.QueryNumberOfCommitmentsRequest           `json:"commitment_number_of_commitments,omitempty"`

	// epochs queriers
	EpochsEpochInfos   *epochstypes.QueryEpochsInfoRequest   `json:"epochs_epoch_infos,omitempty"`
	EpochsCurrentEpoch *epochstypes.QueryCurrentEpochRequest `json:"epochs_current_epoch,omitempty"`

	// incentive queriers
	IncentiveApr  *incentivetypes.QueryAprRequest  `json:"incentive_apr,omitempty"`
	IncentiveAprs *incentivetypes.QueryAprsRequest `json:"incentive_aprs,omitempty"`

	// masterchef queriers
	MasterchefParams            *mastercheftypes.QueryParamsRequest            `json:"masterchef_params,omitempty"`
	MasterchefExternalIncentive *mastercheftypes.QueryExternalIncentiveRequest `json:"masterchef_external_incentive,omitempty"`
	MasterchefPoolInfo          *mastercheftypes.QueryPoolInfoRequest          `json:"masterchef_pool_info,omitempty"`
	MasterchefPoolRewardInfo    *mastercheftypes.QueryPoolRewardInfoRequest    `json:"masterchef_pool_reward_info,omitempty"`
	MasterchefUserRewardInfo    *mastercheftypes.QueryUserRewardInfoRequest    `json:"masterchef_user_reward_info,omitempty"`
	MasterchefUserPendingReward *mastercheftypes.QueryUserPendingRewardRequest `json:"masterchef_user_pending_reward,omitempty"`
	MasterchefStableStakeApr    *mastercheftypes.QueryStableStakeAprRequest    `json:"masterchef_stable_stake_apr,omitempty"`
	MasterchefPoolAprs          *mastercheftypes.QueryPoolAprsRequest          `json:"masterchef_pool_aprs,omitempty"`

	// estaking queriers
	EstakingParams  *estakingtypes.QueryParamsRequest  `json:"estaking_params,omitempty"`
	EstakingRewards *estakingtypes.QueryRewardsRequest `json:"estaking_rewards,omitempty"`

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
	LeveragelpOpenEst                  *leveragelptypes.QueryOpenEstRequest        `json:"leveragelp_open_estimation,omitempty"`
	LeveragelpCloseEst                 *leveragelptypes.QueryCloseEstRequest       `json:"leveragelp_close_estimation,omitempty"`
	LeveragelpRewards                  *leveragelptypes.QueryRewardsRequest        `json:"leveragelp_rewards,omitempty"`

	// perpetual queriers
	PerpetualParams                 *perpetualtypes.ParamsRequest              `json:"perpetual_params,omitempty"`
	PerpetualQueryPositions         *perpetualtypes.PositionsRequest           `json:"perpetual_query_positions,omitempty"`
	PerpetualQueryPositionsByPool   *perpetualtypes.PositionsByPoolRequest     `json:"perpetual_query_positions_by_pool,omitempty"`
	PerpetualGetStatus              *perpetualtypes.StatusRequest              `json:"perpetual_get_status,omitempty"`
	PerpetualGetPositionsForAddress *perpetualtypes.PositionsForAddressRequest `json:"perpetual_get_positions_for_address,omitempty"`
	PerpetualGetWhitelist           *perpetualtypes.WhitelistRequest           `json:"perpetual_get_whitelist,omitempty"`
	PerpetualIsWhitelisted          *perpetualtypes.IsWhitelistedRequest       `json:"perpetual_is_whitelisted,omitempty"`
	PerpetualPool                   *perpetualtypes.QueryGetPoolRequest        `json:"perpetual_pool,omitempty"`
	PerpetualPools                  *perpetualtypes.QueryAllPoolRequest        `json:"perpetual_pools,omitempty"`
	PerpetualMTP                    *perpetualtypes.MTPRequest                 `json:"perpetual_mtp,omitempty"`
	PerpetualOpenEstimation         *perpetualtypes.QueryOpenEstimationRequest `json:"perpetual_open_estimation,omitempty"`

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
	ParameterParams *parametertypes.QueryParamsRequest `json:"parameter_params,omitempty"`

	// stablestake queriers
	StableStakeParams          *stablestaketypes.QueryParamsRequest      `json:"stable_stake_params,omitempty"`
	StableStakeBalanceOfBorrow *stablestaketypes.QueryBorrowRatioRequest `json:"stable_stake_balance_of_borrow,omitempty"`

	// tokenomics queriers
	TokenomicsParams                *tokenomicstypes.QueryParamsRequest                `json:"tokenomics_params,omitempty"`
	TokenomicsAirdrop               *tokenomicstypes.QueryGetAirdropRequest            `json:"tokenomics_airdrop,omitempty"`
	TokenomicsAirdropAll            *tokenomicstypes.QueryAllAirdropRequest            `json:"tokenomics_airdrop_all,omitempty"`
	TokenomicsGenesisInflation      *tokenomicstypes.QueryGetGenesisInflationRequest   `json:"tokenomics_genesis_inflation,omitempty"`
	TokenomicsTimeBasedInflation    *tokenomicstypes.QueryGetTimeBasedInflationRequest `json:"tokenomics_time_based_inflation,omitempty"`
	TokenomicsTimeBasedInflationAll *tokenomicstypes.QueryAllTimeBasedInflationRequest `json:"tokenomics_time_based_inflation_all,omitempty"`

	// transferhook queriers
	TransferHookParams *transferhooktypes.QueryParamsRequest `json:"transfer_hook_params,omitempty"`

	// tier queries
	TierCalculateDiscount *tiertypes.QueryCalculateDiscountRequest `json:"tier_calculate_discount,omitempty"`
}

type CustomMessenger struct {
	wrapped          wasmkeeper.Messenger
	moduleMessengers []ModuleMessenger
	accountedpool    *accountedpoolkeeper.Keeper
	amm              *ammkeeper.Keeper
	assetprofile     *assetprofilekeeper.Keeper
	auth             *authkeeper.AccountKeeper
	bank             *bankkeeper.BaseKeeper
	burner           *burnerkeeper.Keeper
	clock            *clockkeeper.Keeper
	commitment       *commitmentkeeper.Keeper
	epochs           *epochskeeper.Keeper
	incentive        *incentivekeeper.Keeper
	masterchef       *masterchefkeeper.Keeper
	estaking         *estakingkeeper.Keeper
	leveragelp       *leveragelpkeeper.Keeper
	perpetual        *perpetualkeeper.Keeper
	oracle           *oraclekeeper.Keeper
	parameter        *parameterkeeper.Keeper
	stablestake      *stablestakekeeper.Keeper
	staking          *stakingkeeper.Keeper
	tokenomics       *tokenomicskeeper.Keeper
	transferhook     *transferhookkeeper.Keeper
	tier             *tierkeeper.Keeper
}

type ElysMsg struct {
	// accountedpool messages

	// amm messages
	AmmCreatePool         *ammtypes.MsgCreatePool         `json:"amm_create_pool,omitempty"`
	AmmJoinPool           *ammtypes.MsgJoinPool           `json:"amm_join_pool,omitempty"`
	AmmExitPool           *ammtypes.MsgExitPool           `json:"amm_exit_pool,omitempty"`
	AmmSwapExactAmountIn  *ammtypes.MsgSwapExactAmountIn  `json:"amm_swap_exact_amount_in,omitempty"`
	AmmSwapExactAmountOut *ammtypes.MsgSwapExactAmountOut `json:"amm_swap_exact_amount_out,omitempty"`
	AmmSwapByDenom        *ammtypes.MsgSwapByDenom        `json:"amm_swap_by_denom,omitempty"`

	// assetprofile messages
	// auth messages
	// burner messages
	// clock messages

	// commitment messages
	CommitmentCommitUnclaimedRewards *commitmenttypes.MsgCommitClaimedRewards `json:"commitment_commit_unclaimed_rewards,omitempty"`
	CommitmentUncommitTokens         *commitmenttypes.MsgUncommitTokens       `json:"commitment_uncommit_tokens,omitempty"`
	CommitmentClaimVesting           *commitmenttypes.MsgClaimVesting         `json:"commitment_claim_vesting,omitempty"`
	CommitmentVest                   *commitmenttypes.MsgVest                 `json:"commitment_vest"`
	CommitmentVestNow                *commitmenttypes.MsgVestNow              `json:"commitment_vest_now"`
	CommitmentVestLiquid             *commitmenttypes.MsgVestLiquid           `json:"commitment_vest_liquid"`
	CommitmentCancelVest             *commitmenttypes.MsgCancelVest           `json:"commitment_cancel_vest"`
	CommitmentUpdateVestingInfo      *commitmenttypes.MsgUpdateVestingInfo    `json:"commitment_update_vesting_info"`
	CommitmentStake                  *commitmenttypes.MsgStake                `json:"commitment_stake,omitempty"`
	CommitmentUnstake                *commitmenttypes.MsgUnstake              `json:"commitment_unstake,omitempty"`

	// epochs messages

	// incentive messages
	IncentiveBeginRedelegate           *incentivetypes.MsgBeginRedelegate           `json:"incentive_begin_redelegate,omitempty"`
	IncentiveCancelUnbondingDelegation *incentivetypes.MsgCancelUnbondingDelegation `json:"incentive_cancel_unbonding_delegation"`

	// masterchef messages
	MasterchefClaimRewards *mastercheftypes.MsgClaimRewards `json:"masterchef_claim_rewards,omitempty"`

	// estaking messages
	EstakingWithdrawReward             *estakingtypes.MsgWithdrawReward             `json:"estaking_withdraw_reward,omitempty"`
	EstakingWithdrawElysStakingRewards *estakingtypes.MsgWithdrawElysStakingRewards `json:"estaking_withdraw_elys_staking_rewards,omitempty"`

	// leveragelp messages
	LeveragelpOpen           *leveragelptypes.MsgOpen           `json:"leveragelp_open,omitempty"`
	LeveragelpClose          *leveragelptypes.MsgClose          `json:"leveragelp_close,omitempty"`
	LeveragelpClaimRewards   *leveragelptypes.MsgClaimRewards   `json:"leveragelp_claim_rewards,omitempty"`
	LeveragelpUpdateStopLoss *leveragelptypes.MsgUpdateStopLoss `json:"leveragelp_update_stop_loss,omitempty"`

	// perpetual messages
	PerpetualOpen          *perpetualtypes.MsgBrokerOpen          `json:"perpetual_open,omitempty"`
	PerpetualClose         *perpetualtypes.MsgBrokerClose         `json:"perpetual_close,omitempty"`
	PerpetualAddCollateral *perpetualtypes.MsgBrokerAddCollateral `json:"perpetual_add_collateral,omitempty"`

	// oracle messages
	// parameter messages

	// stablestake messages
	StakestakeBond   *stablestaketypes.MsgBond   `json:"stablestake_bond,omitempty"`
	StakestakeUnbond *stablestaketypes.MsgUnbond `json:"stablestake_unbond,omitempty"`

	// tokenomics messages
	// transferhook messages

	// tier messages
	TierSetPortfolio *tiertypes.MsgSetPortfolio `json:"tier_set_portfolio,omitempty"`
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
