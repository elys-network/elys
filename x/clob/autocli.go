package clob

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v6/api/elys/clob"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: clob.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current parameters",
				},
				{
					RpcMethod: "AllMarkets",
					Use:       "markets",
					Short:     "Query all markets",
				},
				{
					RpcMethod:      "Market",
					Use:            "market [market_id]",
					Short:          "Query market",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}},
				},
				{
					RpcMethod:      "OrderBook",
					Use:            "order-book [market_id] [is_buy]",
					Short:          "Query all market orders",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}, {ProtoField: "is_buy"}},
				},
				{
					RpcMethod: "AllPerpetualADL",
					Use:       "all_perpetual_adl",
					Short:     "Query All perpetual ADL",
				},
				{
					RpcMethod:      "SubAccounts",
					Use:            "subaccounts [address]",
					Short:          "Query all subaccounts for address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "OwnerPerpetuals",
					Use:            "owner_perpetuals [address] [sub_account_id]",
					Short:          "Query all owner perpetuals",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "sub_account_id"}},
				},
				{
					RpcMethod:      "OwnerPerpetualOrder",
					Use:            "owner_orders [address] [sub_account_id]",
					Short:          "Query perpetual order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "sub_account_id"}},
				},
				{
					RpcMethod: "AllPerpetuals",
					Use:       "all-perpetuals",
					Short:     "Query all perpetuals with liquidation price and health",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: clob.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
