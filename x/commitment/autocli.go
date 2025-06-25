package commitment

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/v6/api/elys/commitment"
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
				{
					RpcMethod:      "RewardProgram",
					Use:            "show-reward-program [address]",
					Short:          "Query reward program",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "TotalRewardProgramClaimed",
					Use:       "total-reward-program-claimed",
					Short:     "Query total reward program claimed",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              commitment.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "CommitClaimedRewards",
					Use:            "commit-claimed-rewards [amount] [denom]",
					Short:          "Broadcast message commit-claimed-rewards",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod:      "UncommitTokens",
					Use:            "uncommit-tokens [amount] [denom]",
					Short:          "Broadcast message uncommit-tokens",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod:      "Vest",
					Use:            "vest [amount] [denom]",
					Short:          "Broadcast message vest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod:      "VestNow",
					Use:            "vest-now [amount] [denom]",
					Short:          "Broadcast message vest-now",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod:      "VestLiquid",
					Use:            "vest-liquid [amount] [denom]",
					Short:          "Broadcast message vest-liquid",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod:      "CancelVest",
					Use:            "cancel-vest [amount] [denom]",
					Short:          "Broadcast message cancel_vest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "denom"}},
				},
				{
					RpcMethod: "ClaimVesting",
					Use:       "claim-vesting",
					Short:     "Broadcast message claim_vesting",
				},
				{
					RpcMethod: "UpdateVestingInfo",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdateEnableVestNow",
					Skip:      true, // authority gated
				},
				{
					RpcMethod: "UpdateAirdropParams",
					Skip:      true, // authority gated
				},
				{
					RpcMethod:      "Stake",
					Use:            "stake [amount] [asset] [validator-address]",
					Short:          "Stake Elys tokens",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "asset"}, {ProtoField: "validator_address"}},
				},
				{
					RpcMethod:      "Unstake",
					Use:            "unstake [amount] [asset] [validator-address]",
					Short:          "Unstake Elys tokens",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "asset"}, {ProtoField: "validator_address"}},
				},
				{
					RpcMethod: "ClaimAirdrop",
					Use:       "claim-airdrop",
					Short:     "Broadcast message claim_airdrop",
				},
				{
					RpcMethod:      "ClaimKol",
					Use:            "claim-kol [refund]",
					Short:          "Broadcast message claim_kol",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "refund"}},
				},
			},
		},
	}
}
