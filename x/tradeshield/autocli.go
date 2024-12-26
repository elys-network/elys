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
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              tradeshield.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "CreateSpotOrder",
					Use:            "create-spot-order [order-type] [order-amount] [order-target-denom] [order-price] [owner-address]",
					Short:          "Create a new spot order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "order_type"}, {ProtoField: "order_amount"}, {ProtoField: "order_target_denom"}, {ProtoField: "order_price"}, {ProtoField: "owner_address"}},
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
					RpcMethod:      "CreatePerpetualOpenOrder",
					Use:            "create-perpetual-open-order [position] [leverage] [pool-id] [trading-asset] [collateral] [trigger-price]",
					Short:          "Create a new perpetual open order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "position"}, {ProtoField: "leverage"}, {ProtoField: "pool_id"}, {ProtoField: "trading_asset"}, {ProtoField: "collateral"}, {ProtoField: "trigger_price"}},
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
					RpcMethod:      "ExecuteOrders",
					Use:            "execute-orders [perpetual-order-ids] [spot-order-ids]",
					Short:          "Verify that submitted orders meet the criteria for execution and process those that do, while skipping those that don't.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "perpetual_order_ids", Varargs: true}, {ProtoField: "spot_order_ids", Varargs: true}},
				},
			},
		},
	}
}
