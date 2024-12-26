package estaking

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/estaking"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: estaking.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "Rewards",
					Use:            "rewards [address]",
					Short:          "shows the rewards of an account",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "Invariant",
					Use:       "invariant",
					Short:     "Query invariant values",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              estaking.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "WithdrawReward",
					Skip:      true,
				},
				{
					RpcMethod:      "WithdrawAllRewards",
					Use:            "withdraw-all-rewards [delegator-address]",
					Short:          "Withdraw all rewards for delegations and Eden/EdenB commit",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "delegator_address"}},
				},
				{
					RpcMethod:      "WithdrawElysStakingRewards",
					Use:            "withdraw-elys-staking-rewards [delegator-address]",
					Short:          "Withdraw rewards for delegations",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "delegator_address"}},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // authority gated
				},
			},
		},
	}
}
