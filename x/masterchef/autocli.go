package masterchef

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/masterchef"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: masterchef.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "ExternalIncentive",
					Use:            "external-incentive",
					Short:          "shows external incentive",
					Example:        "elysd q masterchef external-incentive [id]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "PoolInfo",
					Use:            "pool-info",
					Short:          "shows pool info",
					Example:        "elysd q masterchef pool-info [id]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
				{
					RpcMethod:      "PoolRewardInfo",
					Use:            "pool-reward-info",
					Short:          "shows pool reward info",
					Example:        "elysd q masterchef pool-reward-info [id] [reward-denom]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "reward_denom"}},
				},
				{
					RpcMethod:      "UserRewardInfo",
					Use:            "user-reward-info",
					Short:          "shows user reward info",
					Example:        "elysd q masterchef user-reward-info [user] [pool-id] [reward-denom]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}, {ProtoField: "pool_id"}, {ProtoField: "reward_denom"}},
				},
				{
					RpcMethod:      "UserPendingReward",
					Use:            "user-pending-reward",
					Short:          "shows user pending reward",
					Example:        "elysd q masterchef user-pending-reward [user]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "StableStakeApr",
					Use:            "stable-stake-apr",
					Short:          "calculate Stable Stake APR",
					Example:        "elysd q masterchef stable-stake-apr [denom]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "PoolAprs",
					Use:            "pool-aprs",
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
					RpcMethod:      "Apr",
					Use:            "apr",
					Short:          "calculate APR",
					Example:        "elysd q masterchef apr [withdraw-type] [denom]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "withdraw_type"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod: "Aprs",
					Use:       "aprs",
					Short:     "Query aprs",
					Example:   "elysd q masterchef aprs",
				},
				{
					RpcMethod:      "PoolRewards",
					Use:            "pool-rewards",
					Short:          "calculate pool rewards",
					Example:        "elysd q masterchef pool-rewards [ids]",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_ids", Varargs: true}},
				},
			},
		},
	}
}
