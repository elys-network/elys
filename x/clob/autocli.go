package clob

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/clob"
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
					RpcMethod:      "SubAccounts",
					Use:            "subaccounts [address]",
					Short:          "Query all subaccounts for address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "MarketOrders",
					Use:            "market_orders [market_id] [is_buy]",
					Short:          "Query all market orders",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}, {ProtoField: "is_buy"}},
				},
				{
					RpcMethod:      "OwnerPerpetuals",
					Use:            "owner_perpetuals [owner]",
					Short:          "Query all owner perpetuals",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "owner"}},
				},
				{
					RpcMethod:      "PerpetualOrder",
					Use:            "perpetual_order [market_id] [order_type] [price] [block_height]",
					Short:          "Query perpetual order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}, {ProtoField: "order_type"}, {ProtoField: "price"}, {ProtoField: "block_height"}},
				},
				{
					RpcMethod: "AllPerpetualOrder",
					Use:       "all_perpetual_order",
					Short:     "Query All perpetual orders",
				},
				{
					RpcMethod:      "CurrentTwapPrice",
					Use:            "current_twap_price [market_id]",
					Short:          "Query current twap price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}},
				},
				{
					RpcMethod: "AllTwapPrices",
					Use:       "all_twap_price",
					Short:     "Query All twap price",
				},
				{
					RpcMethod:      "LastAverageTradePrice",
					Use:            "last_average_trade_price [market_id]",
					Short:          "Query last average trade price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}},
				},
				{
					RpcMethod:      "HighestBuyPrice",
					Use:            "highest_buy_price [market_id]",
					Short:          "Query highest buy price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}},
				},
				{
					RpcMethod:      "LowestSellPrice",
					Use:            "lowest_sell_price [market_id]",
					Short:          "Query lowest sell price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}},
				},
				{
					RpcMethod:      "MidPrice",
					Use:            "mid_price [market_id]",
					Short:          "Query mid price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}},
				},
				{
					RpcMethod:      "PerpetualADL",
					Use:            "perpetual_adl [market_id] [id]",
					Short:          "Query perpetual ADL",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "market_id"}, {ProtoField: "id"}},
				},
				{
					RpcMethod: "AllPerpetualADL",
					Use:       "all_perpetual_adl",
					Short:     "Query All perpetual ADL",
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
