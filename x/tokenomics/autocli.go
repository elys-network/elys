package tokenomics

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/tokenomics"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: tokenomics.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "Airdrop",
					Use:            "show-airdrop [intent]",
					Short:          "shows a airdrop",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "intent"}},
				},
				{
					RpcMethod: "AirdropAll",
					Use:       "list-airdrop",
					Short:     "list all airdrop",
				},
				{
					RpcMethod: "GenesisInflation",
					Use:       "show-genesis-inflation",
					Short:     "shows genesis-inflation",
				},
				{
					RpcMethod:      "TimeBasedInflation",
					Use:            "show-time-based-inflation [start-block-height] [end-block-height]",
					Short:          "shows a time-based-inflation",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "start_block_height"}, {ProtoField: "end_block_height"}},
				},
				{
					RpcMethod: "TimeBasedInflationAll",
					Use:       "list-time-based-inflation",
					Short:     "list all time-based-inflation",
				},
			},
		},
	}
}
