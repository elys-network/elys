package tradeshield

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v6/api/elys/tradeshield"
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
					Use:            "show-pending-perpetual-order [owner-address] [pool-id] [order-id]",
					Short:          "shows a pending-perpetual-order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "owner_address"}, {ProtoField: "pool_id"}, {ProtoField: "order_id"}},
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
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              tradeshield.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "CreateSpotOrder",
					Use:            "create-spot-order [order-type] [order-amount] [order-target-denom] [order-price]",
					Short:          "Create a new spot order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "order_type"}, {ProtoField: "order_amount"}, {ProtoField: "order_target_denom"}, {ProtoField: "order_price"}},
				},
				{
					RpcMethod:      "UpdateSpotOrder",
					Use:            "update-spot-order [order-id] [order-price]",
					Short:          "Update a spot order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "order_id"}, {ProtoField: "order_price"}},
				},
				{
					RpcMethod:      "CancelSpotOrder",
					Use:            "cancel-spot-order [order-id]",
					Short:          "Broadcast message cancel-spot-order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "order_id"}},
				},
				{
					RpcMethod:      "CancelSpotOrders",
					Use:            "cancel-spot-orders [order-ids]",
					Short:          "Cancel spot-orders",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "spot_order_ids", Varargs: true}},
				},
				{
					RpcMethod: "CreatePerpetualOpenOrder",
					Skip:      true, // use custom command
				},
				{
					RpcMethod:      "CreatePerpetualCloseOrder",
					Use:            "create-perpetual-close-order [trigger-price] [position-id]",
					Short:          "Create a new perpetual close order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trigger_price"}, {ProtoField: "position_id"}},
				},
				{
					RpcMethod:      "UpdatePerpetualOrder",
					Use:            "update-perpetual-order [id] [trigger-price]",
					Short:          "Update a perpetual order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "order_id"}, {ProtoField: "trigger_price"}},
				},
				{
					RpcMethod:      "CancelPerpetualOrder",
					Use:            "cancel-perpetual-order [order-id]",
					Short:          "Broadcast message cancel-perpetual-order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "order_id"}},
				},
				{
					RpcMethod:      "CancelPerpetualOrders",
					Use:            "cancel-perpetual-orders [order-ids]",
					Short:          "Cancel a perpetual orders by ids",
					Example:        "elysd tx perpetual cancel-perpetual-orders 1 2 3... --from=bob --yes --gas=1000000",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "order_ids", Varargs: true}},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, //	authority gated
				},
				{
					RpcMethod: "ExecuteOrders",
					Skip:      true, // use custom command
				},
			},
		},
	}
}
