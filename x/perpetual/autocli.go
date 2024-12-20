package perpetual

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/perpetual"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: perpetual.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod: "GetPositions",
					Use:       "get-positions",
					Short:     "Query get-positions",
				},
				{
					RpcMethod:      "GetPositionsByPool",
					Use:            "get-positions-by-pool [amm_pool_id]",
					Short:          "Query get-positions-by-pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amm_pool_id"}},
				},
				{
					RpcMethod: "GetStatus",
					Use:       "get-status",
					Short:     "Query get-status",
				},
				{
					RpcMethod:      "GetPositionsForAddress",
					Use:            "get-positions-for-address [address]",
					Short:          "Query get-positions-for-address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "GetWhitelist",
					Use:       "get-whitelist",
					Short:     "Query get-whitelist",
				},
				{
					RpcMethod:      "IsWhitelisted",
					Use:            "is-whitelisted [address]",
					Short:          "Query is-whitelisted",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "Pool",
					Use:            "show-pool [index]",
					Short:          "shows a pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod: "Pools",
					Use:       "list-pool",
					Short:     "list all pool",
				},
				{
					RpcMethod:      "MTP",
					Use:            "get-mtp [address] [id]",
					Short:          "Query mtp",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "id"}},
				},
				{
					RpcMethod:      "OpenEstimation",
					Use:            "open-estimation [position] [leverage] [trading-asset] [collateral] [pool-id]",
					Short:          "Query open-estimation",
					Example:        "elysd q perpetual open-estimation long 5 uatom 100000000uusdc 1",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "position"}, {ProtoField: "leverage"}, {ProtoField: "trading_asset"}, {ProtoField: "collateral"}, {ProtoField: "poolId"}},
				},
				{
					RpcMethod:      "CloseEstimation",
					Use:            "close-estimation [address] [position-id] [closing-amount]",
					Short:          "Query close-estimation",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "position_id"}, {ProtoField: "close_amount"}},
				},
			},
		},
	}
}
