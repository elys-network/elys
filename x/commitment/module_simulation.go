package commitment

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCommitTokens = "op_weight_msg_commit_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCommitTokens int = 100

	opWeightMsgUncommitTokens = "op_weight_msg_uncommit_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUncommitTokens int = 100

	opWeightMsgWithdrawTokens = "op_weight_msg_withdraw_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWithdrawTokens int = 100

	opWeightMsgDepositTokens = "op_weight_msg_deposit_tokens"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDepositTokens int = 100

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
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCommitTokens int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCommitTokens, &weightMsgCommitTokens, nil,
		func(_ *rand.Rand) {
			weightMsgCommitTokens = defaultWeightMsgCommitTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCommitTokens,
		commitmentsimulation.SimulateMsgCommitTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUncommitTokens int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUncommitTokens, &weightMsgUncommitTokens, nil,
		func(_ *rand.Rand) {
			weightMsgUncommitTokens = defaultWeightMsgUncommitTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUncommitTokens,
		commitmentsimulation.SimulateMsgUncommitTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgWithdrawTokens int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWithdrawTokens, &weightMsgWithdrawTokens, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawTokens = defaultWeightMsgWithdrawTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgWithdrawTokens,
		commitmentsimulation.SimulateMsgWithdrawTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDepositTokens int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDepositTokens, &weightMsgDepositTokens, nil,
		func(_ *rand.Rand) {
			weightMsgDepositTokens = defaultWeightMsgDepositTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDepositTokens,
		commitmentsimulation.SimulateMsgDepositTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
