package commitment

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/commitment"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: commitment.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "ShowCommitments",
					Use:            "show-commitments [creator]",
					Short:          "Query show-commitments",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}},
				},
				{
					RpcMethod:      "CommittedTokensLocked",
					Use:            "committed-tokens-locked [address]",
					Short:          "Show locked coins in commitment not unlockable",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "NumberOfCommitments",
					Use:       "number-of-commitments",
					Short:     "Query number-of-commitments",
				},
				{
					RpcMethod:      "CommitmentVestingInfo",
					Use:            "commitment-vesting-info [address]",
					Short:          "Query commitment-vesting-info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "AirDrop",
					Use:            "show-airdrop [address]",
					Short:          "Query airdrops",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "TotalAirdropClaimed",
					Use:       "total-airdrop-claimed",
					Short:     "Query total airdrop claimed",
				},
				{
					RpcMethod:      "Kol",
					Use:            "show-kol [address]",
					Short:          "Query kols",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
			},
		},
	}
}
