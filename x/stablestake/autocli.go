package stablestake

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v6/api/elys/stablestake"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: stablestake.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current parameters",
				},
				{
					RpcMethod:      "BorrowRatio",
					Use:            "borrow-ratio [pool-id]",
					Short:          "Query the borrow ratio",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
				{
					RpcMethod:      "AmmPool",
					Use:            "amm-pool [pool-id]",
					Short:          "Query the amm pool liabilities",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "AllAmmPools",
					Use:       "all-amm-pools",
					Short:     "Query all amm pools liabilities",
				},
				{
					RpcMethod:      "Pool",
					Use:            "pool [pool-id]",
					Short:          "Query pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
				{
					RpcMethod: "Pools",
					Use:       "pools",
					Short:     "Query all pools",
				},
				{
					RpcMethod:      "Debt",
					Use:            "debt [pool-id] [address]",
					Short:          "Query debt",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "address"}},
				},
				{
					RpcMethod:      "GetInterest",
					Use:            "interest [pool-id] [block-height]",
					Short:          "Query interest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "block_height"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              stablestake.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "Bond",
					Use:            "bond [amount]",
					Short:          "Bond coins to the stablestake module",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}},
				},
				{
					RpcMethod:      "Unbond",
					Use:            "unbond [amount]",
					Short:          "Unbond coins from the stablestake module",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
