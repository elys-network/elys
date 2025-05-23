package masterchef

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v5/api/elys/masterchef"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:              masterchef.Query_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "ExternalIncentive",
					Use:            "external-incentive [id]",
					Short:          "shows external incentive",
					Example:        "elysd q masterchef external-incentive [id]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "PoolInfo",
					Use:            "pool-info [id]",
					Short:          "shows pool info",
					Example:        "elysd q masterchef pool-info [id]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
				{
					RpcMethod: "ListPoolInfos",
					Use:       "list-pool-infos",
					Short:     "shows all pool infos",
					Example:   "elysd q masterchef list-pool-infos",
				},
				{
					RpcMethod:      "PoolRewardInfo",
					Use:            "pool-reward-info [id] [reward-denom]",
					Short:          "shows pool reward info",
					Example:        "elysd q masterchef pool-reward-info [id] [reward-denom]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "reward_denom"}},
				},
				{
					RpcMethod:      "UserRewardInfo",
					Use:            "user-reward-info [user] [pool-id] [reward-denom]",
					Short:          "shows user reward info",
					Example:        "elysd q masterchef user-reward-info [user] [pool-id] [reward-denom]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}, {ProtoField: "pool_id"}, {ProtoField: "reward_denom"}},
				},
				{
					RpcMethod:      "UserPendingReward",
					Use:            "user-pending-reward [user]",
					Short:          "shows user pending reward",
					Example:        "elysd q masterchef user-pending-reward [user]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "StableStakeApr",
					Use:            "stable-stake-apr [denom]",
					Short:          "calculate Stable Stake APR",
					Example:        "elysd q masterchef stable-stake-apr [denom]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "PoolAprs",
					Use:            "pool-aprs [ids]",
					Short:          "calculate pool APRs",
					Example:        "elysd q masterchef pool-aprs [ids]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_ids", Varargs: true}},
				},
				{
					RpcMethod:      "ShowFeeInfo",
					Use:            "show-fee-info [date]",
					Short:          "Query show-fee-info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "date"}},
				},
				{
					RpcMethod: "ListFeeInfo",
					Use:       "list-fee-info",
					Short:     "Query list-fee-info",
				},
				{
					RpcMethod: "Apr",
					Skip:      true, // use custom command
				},
				{
					RpcMethod:      "Aprs",
					Use:            "aprs [days]",
					Short:          "Query aprs",
					Example:        "elysd q masterchef aprs [days]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "days"}},
				},
				{
					RpcMethod: "AllLiquidityPoolTVL",
					Use:       "all-liquidity-pool-tvl",
					Short:     "show all pools tvl",
				},
				{
					RpcMethod: "ChainTVL",
					Use:       "chain-tvl",
					Short:     "show chain tvl",
				},
				{
					RpcMethod:      "PoolRewards",
					Use:            "pool-rewards [pool-ids]",
					Short:          "calculate pool rewards",
					Example:        "elysd q masterchef pool-rewards [ids]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_ids", Varargs: true}},
				},
				{
					RpcMethod: "TotalPendingRewards",
					Use:       "total-pending-rewards",
					Short:     "show total pending rewards",
				},
				{
					RpcMethod: "PendingRewards",
					Use:       "pending-rewards",
					Short:     "show pending rewards",
					Example:   "elysd q masterchef pending-rewards",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              masterchef.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "AddExternalRewardDenom",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdatePoolMultipliers",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "TogglePoolEdenRewards",
					Skip:      true, // authority gated
				},
				{
					RpcMethod:      "AddExternalIncentive",
					Use:            "add-external-incentive [reward-denom] [pool-id] [from-block] [to-block] [amount-per-block]",
					Short:          "Broadcast message add-external-incentive",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "reward_denom"}, {ProtoField: "pool_id"}, {ProtoField: "from_block"}, {ProtoField: "to_block"}, {ProtoField: "amount_per_block"}},
				},
				{
					RpcMethod:      "ClaimRewards",
					Use:            "claim-rewards [pool-ids]",
					Short:          "claim rewards including external incentives",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_ids", Varargs: true}},
				},
			},
		},
	}
}
