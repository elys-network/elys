package estaking

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v6/api/elys/estaking"
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
				{
					RpcMethod: "EdenBBurnAmount",
					Use:       "edenb-burn-amount [address] [token_type] [amount]",
					Short:     "Query the amount of EdenB that will be burned when unstaking",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
						{ProtoField: "token_type"},
						{ProtoField: "amount"},
					},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              estaking.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "WithdrawReward",
					Use:            "withdraw-rewards [validator-address]",
					Short:          "Withdraw rewards",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "validator_address"}},
				},
				{
					RpcMethod: "WithdrawAllRewards",
					Use:       "withdraw-all-rewards",
					Short:     "Withdraw all rewards for delegations and Eden/EdenB commit",
				},
				{
					RpcMethod: "WithdrawElysStakingRewards",
					Use:       "withdraw-elys-staking-rewards",
					Short:     "Withdraw rewards for delegations",
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // authority gated
				},
			},
		},
	}
}
