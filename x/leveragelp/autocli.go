package leveragelp

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v5/api/elys/leveragelp"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:              leveragelp.Query_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod: "QueryPositions",
					Use:       "get-positions",
					Short:     "Query get-positions",
				},
				{
					RpcMethod:      "QueryPositionsByPool",
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
					RpcMethod:      "QueryPositionsForAddress",
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
					RpcMethod:      "Position",
					Use:            "get-position [address] [id]",
					Short:          "Query position",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "id"}},
				},
				{
					RpcMethod:      "LiquidationPrice",
					Use:            "liquidation-price [address] [position_id]",
					Short:          "Query liquidation price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "position_id"}},
				},
				{
					RpcMethod: "OpenEst",
					Skip:      true, // use custom command
				},
				{
					RpcMethod:      "CloseEst",
					Use:            "close-estimation [owner] [pool_id] [lp_amount]",
					Short:          "Query close estimation",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "owner"}, {ProtoField: "id"}, {ProtoField: "lp_amount"}},
				},
				{
					RpcMethod:      "Rewards",
					Use:            "rewards [address] [position-ids]",
					Short:          "Query rewards for position ids owned by address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "ids", Varargs: true}},
				},
				{
					RpcMethod:      "CommittedTokensLocked",
					Use:            "committed-tokens-locked [address]",
					Short:          "Show locked coins in commitment not unlockable for different leveragelp positions",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              leveragelp.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Open",
					Skip:      true, // use custom command
				},
				{
					RpcMethod:      "Close",
					Use:            "close [position-id] [amount]",
					Short:          "Close leveragelp position",
					Example:        `elysd tx leveragelp close 1 10000000 --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "lp_amount"}},
				},
				{
					RpcMethod:      "ClaimRewards",
					Use:            "claim-rewards [position-ids]",
					Short:          "Claim rewards from leveragelp position",
					Example:        `elysd tx leveragelp claim-rewards 1 2 3 4 ... --from=bob --yes --gas=1000000`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "ids", Varargs: true}},
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
					RpcMethod: "UpdateParams",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "AddPool",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdateEnabledPools",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdateMaxLeverageForPool",
					Skip:      true, // authority gated
				},
			},
		},
	}
}
