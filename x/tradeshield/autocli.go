package tradeshield

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/tradeshield"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: tradeshield.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "PendingSpotOrder",
					Use:            "show-pending-spot-order [id]",
					Short:          "shows a pending-spot-order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "PendingSpotOrderAll",
					Use:       "list-pending-spot-order",
					Short:     "list all pending-spot-order",
				},
				{
					RpcMethod:      "PendingPerpetualOrder",
					Use:            "show-pending-perpetual-order [id]",
					Short:          "shows a pending-perpetual-order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "PendingPerpetualOrderAll",
					Use:       "list-pending-perpetual-order",
					Short:     "list all pending-perpetual-order",
				},
				{
					RpcMethod:      "PendingPerpetualOrderForAddress",
					Use:            "pending-perpetual-order-for-address [address]",
					Short:          "Query pending-perpetual-order-for-address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "PendingSpotOrderForAddress",
					Use:            "pending-spot-order-for-address [address]",
					Short:          "Query pending-spot-order-for-address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
			},
		},
	}
}
