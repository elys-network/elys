package membershiptier

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/tier"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: tier.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "Portfolio",
					Use:            "show-portfolio [user]",
					Short:          "shows a portfolio",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod: "PortfolioAll",
					Use:       "list-portfolio",
					Short:     "list all portfolio",
				},
				{
					RpcMethod:      "CalculateDiscount",
					Use:            "calculate-discount [user]",
					Short:          "Query calculate-discount",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "LeverageLpTotal",
					Use:            "leverage-lp-total [user]",
					Short:          "Query leverage-lp-total",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "RewardsTotal",
					Use:            "rewards-total [user]",
					Short:          "Query rewards-total",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "StakedPool",
					Use:            "staked-pool [user]",
					Short:          "Query staked-pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "Perpetual",
					Use:            "perpetual [user]",
					Short:          "Query perpetual",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "LiquidTotal",
					Use:            "liquid-total [user]",
					Short:          "Query liquid-total",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "LockedOrder",
					Use:            "locked-order [user]",
					Short:          "Query locked-order",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod:      "GetAmmPrice",
					Use:            "get-amm-price [denom] [decimal]",
					Short:          "Query get-amm-price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "decimal"}},
				},
				{
					RpcMethod:      "GetConsolidatedPrice",
					Use:            "get-consolidated-price [denom]",
					Short:          "Query get-consolidated-price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "Staked",
					Use:            "staked [user]",
					Short:          "Query staked",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
				{
					RpcMethod: "GetUsersPoolData",
					Use:       "get-users-pool-data",
					Short:     "Query to get users pool data",
				},
				{
					RpcMethod: "GetAllPrices",
					Use:       "get-all-prices",
					Short:     "Query get-all-prices",
				},
				{
					RpcMethod:      "GetOraclePrices",
					Use:            "get-oracle-prices [denoms]",
					Short:          "Query get-oracle-prices",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denoms"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              tier.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "SetPortfolio",
					Use:            "set-portfolio [user]",
					Short:          "Broadcast message set-portfolio",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "user"}},
				},
			},
		},
	}
}
