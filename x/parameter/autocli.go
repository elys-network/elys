package parameter

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v5/api/elys/parameter"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: parameter.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current parameters",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: parameter.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateMinCommission",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateMaxVotingPower",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateMinSelfDelegation",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateTotalBlocksPerYear",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateRewardsDataLifetime",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateTakerFees",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
