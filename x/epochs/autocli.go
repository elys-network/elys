package epochs

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	epochsv1 "github.com/elys-network/elys/api/elys/epochs/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: epochsv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "EpochInfos",
					Use:       "epoch-infos",
					Short:     "Query running epochInfos",
				},
				{
					RpcMethod:      "CurrentEpoch",
					Use:            "current-epoch [identifier]",
					Short:          "Query current epoch by specified identifier",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "identifier"}},
				},
			},
		},
	}
}
