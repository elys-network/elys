package assetprofile

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v5/api/elys/assetprofile"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: assetprofile.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current parameters",
				},
				{
					RpcMethod:      "Entry",
					Use:            "show-entry [base-denom]",
					Short:          "Query asset profile entry by base denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "base_denom"}},
				},
				{
					RpcMethod:      "EntryByDenom",
					Use:            "show-entry-by-denom [denom]",
					Short:          "Query asset profile entry by denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod: "EntryAll",
					Use:       "list-entry",
					Short:     "Query all entries",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: assetprofile.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "AddEntry",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "UpdateEntry",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "DeleteEntry",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
