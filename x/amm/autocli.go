package amm

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/amm"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: amm.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "Pool",
					Use:            "show-pool [pool-id]",
					Short:          "shows a pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
				{
					RpcMethod: "PoolAll",
					Use:       "list-pool",
					Short:     "list all pool",
				},
				{
					RpcMethod:      "DenomLiquidity",
					Use:            "show-denom-liquidity [denom]",
					Short:          "shows a denom-liquidity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod: "DenomLiquidityAll",
					Use:       "list-denom-liquidity",
					Short:     "list all denom-liquidity",
				},
				{
					RpcMethod:      "SwapEstimation",
					Use:            "swap-estimation [token-in] [discount] {pool_id token_out_denom}...",
					Short:          "Query SwapEstimation",
					Example:        "elysd q amm swap-estimation 100token 1 token_out1 2 token_out2 ...",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_in"}, {ProtoField: "discount"}, {ProtoField: "routes", Varargs: true}},
				},
				{
					RpcMethod:      "JoinPoolEstimation",
					Use:            "join-pool-estimation [pool_id] [tokens-in]",
					Short:          "Query JoinPoolEstimation",
					Example:        "elysd q amm join-pool-estimation 1 100token,100token2",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "amounts_in", Varargs: true}},
				},
				{
					RpcMethod:      "ExitPoolEstimation",
					Use:            "exit-pool-estimation [pool_id] [share_amount_in] [token_out_denom]",
					Short:          "Query ExitPoolEstimation",
					Example:        "elysd q amm exit-pool-estimation 1 10000 token",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "share_amount_in"}, {ProtoField: "token_out_denom"}},
				},
				{
					RpcMethod:      "SlippageTrack",
					Use:            "tracked-slippage [pool_id]",
					Short:          "Query Tracked Slippage",
					Example:        "elysd q amm tracked-slippage 1",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
				{
					RpcMethod: "SlippageTrackAll",
					Use:       "tracked-slippage-all",
					Short:     "Query All Tracked Slippage",
					Example:   "elysd q amm tracked-slippage-all",
				},
				{
					RpcMethod:      "Balance",
					Use:            "balance [address] [denom]",
					Short:          "Get balance of denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod:      "InRouteByDenom",
					Use:            "in-route-by-denom [denom-in] [denom-out]",
					Example:        "elysd q amm in-route-by-denom uelys uusdc",
					Short:          "Query in-route-by-denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom_in"}, {ProtoField: "denom_out"}},
				},
				{
					RpcMethod:      "OutRouteByDenom",
					Use:            "out-route-by-denom [denom-out] [denom-in]",
					Example:        "elysd q amm out-route-by-denom uusdc uelys",
					Short:          "Query out-route-by-denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom_out"}, {ProtoField: "denom_in"}},
				},
				{
					RpcMethod:      "SwapEstimationByDenom",
					Use:            "swap-estimation-by-denom [amount] [denom-in] [denom-out]",
					Short:          "Query swap-estimation-by-denom",
					Example:        "elysd q amm swap-estimation-by-denom 100uatom uatom uosmo",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom_in"}, {ProtoField: "denom_out"}},
				},
			},
		},
	}
}
