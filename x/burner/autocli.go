package burner

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/burner"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: burner.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "History",
					Use:            "show-history [timestamp] [denom]]",
					Short:          "show a history",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "timestamp"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod: "HistoryAll",
					Use:       "list-history",
					Short:     "list all history",
				},
			},
		},
	}
}
