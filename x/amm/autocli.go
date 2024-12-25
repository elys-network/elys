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
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              amm.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "CreatePool",
					Use:            "create-pool [pool_params] [pool_assets]...",
					Short:          "create a new pool and provide the liquidity to it",
					Example:        `elysd tx amm create-pool '{"swap_fee": "0","use_oracle": true, "fee_denom": "uusdc"}' '{"token": {"denom": "uusdc", "amount": "10000"}, "weight": "100", "external_liquidity_ratio": "0"}' --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_params"}, {ProtoField: "pool_assets", Varargs: true}},
				},
				{
					RpcMethod:      "JoinPool",
					Use:            "join-pool [pool-id] [share-amount-out] [max-amounts-in]",
					Short:          "join a new pool and provide the liquidity to it",
					Example:        `elysd tx amm join-pool 1 1000000000000000 1000000utom 1000000uusdc ... --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "share_amount_out"}, {ProtoField: "max_amounts_in", Varargs: true}},
				},
				{
					RpcMethod:      "ExitPool",
					Use:            "exit-pool [pool-id] [share-amount-in] [min-amounts-out]",
					Short:          "exit a pool and withdraw the liquidity",
					Example:        `elysd tx amm exit-pool 0 200000000000000000 1000uatom 1000uusdc ... --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}, {ProtoField: "share_amount_in"}, {ProtoField: "min_amounts_out", Varargs: true}},
				},
				{
					RpcMethod:      "SwapExactAmountIn",
					Use:            "swap-exact-amount-in [recipient] [token-in] [token-out-min-amount] [routes]",
					Short:          "Swap an exact amount of tokens for a minimum of another token, similar to swapping a token on the trade screen GUI.",
					Example:        `elysd tx amm swap-exact-amount-in elys1l0qs5yedfymc0z3n3r7x95h2dadmsvvgerer6k 1000uusdc 10000 '{"pool_id": 1, "token_out_denom": "uatom"}' --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "recipient"}, {ProtoField: "token_in"}, {ProtoField: "token_out_min_amount"}, {ProtoField: "routes", Varargs: true}},
				},
				{
					RpcMethod:      "SwapExactAmountOut",
					Use:            "swap-exact-amount-out [recipient] [token_out] [token_in_max_amount] [routes]",
					Short:          "Swap a maximum amount of tokens for an exact amount of another token, similar to swapping a token on the trade screen GUI.",
					Example:        `elysd tx amm swap-exact-amount-out elys1mmqs5ys4fymc0z3ngr7xdeh2dadmsvvgerer6k 1000uatom 10000 '{"pool_id": 1, "token_in_denom": "uusdc"}' --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "recipient"}, {ProtoField: "token_out"}, {ProtoField: "token_in_max_amount"}, {ProtoField: "routes", Varargs: true}},
				},
				{
					RpcMethod:      "SwapByDenom",
					Use:            "swap-by-denom [amount] [denom-in] [denom-out]",
					Short:          "Swap an exact amount of tokens for a minimum of another token or a maximum amount of tokens for an exact amount on another token, similar to swapping a token on the trade screen GUI.",
					Example:        `elysd tx amm swap-by-denom 0 1000uatom uusdc --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom_in"}, {ProtoField: "denom_out"}},
				},
				{
					RpcMethod: "FeedMultipleExternalLiquidity",
					Skip:      true,
				},
				{
					RpcMethod: "UpdatePoolParams",
					Skip:      true,
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true,
				},
			},
		},
	}
}
