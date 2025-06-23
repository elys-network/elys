package perpetual

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v6/api/elys/perpetual"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:              perpetual.Query_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true,
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
					RpcMethod:      "PerpetualCounter",
					Use:            "counter [id]",
					Short:          "Query total open positions for a pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
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
					Use:            "get-mtp [address] [id] [pool-id]",
					Short:          "Query mtp",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "id"}, {ProtoField: "pool_id"}},
				},
				{
					RpcMethod: "OpenEstimation",
					Skip:      true, // use custom command
				},
				{
					RpcMethod: "OpenEstimationByFinal",
					Skip:      true, // use custom command
				},
				{
					RpcMethod:      "CloseEstimation",
					Use:            "close-estimation [address] [position-id] [closing-amount] [pool-id]",
					Short:          "Query close-estimation",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "position_id"}, {ProtoField: "close_amount"}, {ProtoField: "pool_id"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              perpetual.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Open",
					Skip:      true, // use custom command
				},
				{
					RpcMethod:      "Close",
					Use:            "close [mtp-id] [amount]",
					Short:          "Close perpetual position",
					Example:        `elysd tx perpetual close 1 10000000 --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "Whitelist",
					Use:            "whitelist [address]",
					Short:          "Whitelist the provided address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "whitelisted_address"}},
				},
				{
					RpcMethod:      "Dewhitelist",
					Use:            "dewhitelist [address]",
					Short:          "Dewhitelist the provided address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "whitelisted_address"}},
				},
				{
					RpcMethod: "UpdateStopLoss",
					Skip:      true, // use custom command
				},
				{
					RpcMethod: "ClosePositions",
					Skip:      true, // use custom command
				},
				{
					RpcMethod: "UpdateTakeProfitPrice",
					Skip:      true, // use custom command
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdateMaxLeverageForPool",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdateEnabledPools",
					Skip:      true, // authority gated
				},
			},
		},
	}
}
