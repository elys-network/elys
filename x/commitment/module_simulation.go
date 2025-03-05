package commitment

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	commitmentsimulation "github.com/elys-network/elys/x/commitment/simulation"
	"github.com/elys-network/elys/x/commitment/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = commitmentsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCommitClaimedRewards = "op_weight_msg_commit_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCommitClaimedRewards int = 100

	opWeightMsgUncommitTokens = "op_weight_msg_uncommit_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUncommitTokens int = 100

	opWeightMsgClaimReward = "op_weight_msg_withdraw_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimReward int = 100

	opWeightMsgCommitLiquidTokens = "op_weight_msg_commit_liquid_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCommitLiquidTokens int = 100

	opWeightMsgVest = "op_weight_msg_vest"
	// TODO: Determine the simulation weight value
	defaultWeightMsgVest int = 100

	opWeightMsgCancelVest = "op_weight_msg_cancel_vest"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCancelVest int = 100

	opWeightMsgVestNow = "op_weight_msg_vest_now"
	// TODO: Determine the simulation weight value
	defaultWeightMsgVestNow int = 100

	opWeightMsgUpdateVestingInfo = "op_weight_msg_update_vesting_info"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateVestingInfo int = 100

	opWeightMsgVestLiquid = "op_weight_msg_vest_liquid"
	// TODO: Determine the simulation weight value
	defaultWeightMsgVestLiquid int = 100

	opWeightMsgClaimRewards = "op_weight_msg_claim_rewards"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimRewards int = 100

	opWeightMsgStake = "op_weight_msg_stake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgStake int = 100

	opWeightMsgUnstake = "op_weight_msg_unstake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnstake int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	commitmentGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&commitmentGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCommitClaimedRewards int
	simState.AppParams.GetOrGenerate(opWeightMsgCommitClaimedRewards, &weightMsgCommitClaimedRewards, nil,
		func(_ *rand.Rand) {
			weightMsgCommitClaimedRewards = defaultWeightMsgCommitClaimedRewards
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCommitClaimedRewards,
		commitmentsimulation.SimulateMsgCommitClaimedRewards(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUncommitTokens int
	simState.AppParams.GetOrGenerate(opWeightMsgUncommitTokens, &weightMsgUncommitTokens, nil,
		func(_ *rand.Rand) {
			weightMsgUncommitTokens = defaultWeightMsgUncommitTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUncommitTokens,
		commitmentsimulation.SimulateMsgUncommitTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClaimReward int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimReward, &weightMsgClaimReward, nil,
		func(_ *rand.Rand) {
			weightMsgClaimReward = defaultWeightMsgClaimReward
		},
	)

	var weightMsgVest int
	simState.AppParams.GetOrGenerate(opWeightMsgVest, &weightMsgVest, nil,
		func(_ *rand.Rand) {
			weightMsgVest = defaultWeightMsgVest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVest,
		commitmentsimulation.SimulateMsgVest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCancelVest int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelVest, &weightMsgCancelVest, nil,
		func(_ *rand.Rand) {
			weightMsgCancelVest = defaultWeightMsgCancelVest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelVest,
		commitmentsimulation.SimulateMsgCancelVest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgVestNow int
	simState.AppParams.GetOrGenerate(opWeightMsgVestNow, &weightMsgVestNow, nil,
		func(_ *rand.Rand) {
			weightMsgVestNow = defaultWeightMsgVestNow
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVestNow,
		commitmentsimulation.SimulateMsgVestNow(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateVestingInfo int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateVestingInfo, &weightMsgUpdateVestingInfo, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateVestingInfo = defaultWeightMsgUpdateVestingInfo
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateVestingInfo,
		commitmentsimulation.SimulateMsgUpdateVestingInfo(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgVestLiquid int
	simState.AppParams.GetOrGenerate(opWeightMsgVestLiquid, &weightMsgVestLiquid, nil,
		func(_ *rand.Rand) {
			weightMsgVestLiquid = defaultWeightMsgVestLiquid
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVestLiquid,
		commitmentsimulation.SimulateMsgVestLiquid(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClaimRewards int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimRewards, &weightMsgClaimRewards, nil,
		func(_ *rand.Rand) {
			weightMsgClaimRewards = defaultWeightMsgClaimRewards
		},
	)

	var weightMsgStake int
	simState.AppParams.GetOrGenerate(opWeightMsgStake, &weightMsgStake, nil,
		func(_ *rand.Rand) {
			weightMsgStake = defaultWeightMsgStake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStake,
		commitmentsimulation.SimulateMsgStake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnstake int
	simState.AppParams.GetOrGenerate(opWeightMsgUnstake, &weightMsgUnstake, nil,
		func(_ *rand.Rand) {
			weightMsgUnstake = defaultWeightMsgUnstake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnstake,
		commitmentsimulation.SimulateMsgUnstake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
