package tokenomics

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	tokenomicssimulation "github.com/elys-network/elys/x/tokenomics/simulation"
	"github.com/elys-network/elys/x/tokenomics/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = tokenomicssimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateAirdrop = "op_weight_msg_airdrop"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateAirdrop int = 100

	opWeightMsgUpdateAirdrop = "op_weight_msg_airdrop"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateAirdrop int = 100

	opWeightMsgDeleteAirdrop = "op_weight_msg_airdrop"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAirdrop int = 100

	opWeightMsgCreateGenesisInflation = "op_weight_msg_genesis_inflation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateGenesisInflation int = 100

	opWeightMsgUpdateGenesisInflation = "op_weight_msg_genesis_inflation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateGenesisInflation int = 100

	opWeightMsgDeleteGenesisInflation = "op_weight_msg_genesis_inflation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteGenesisInflation int = 100

	opWeightMsgCreateTimeBasedInflation = "op_weight_msg_time_based_inflation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateTimeBasedInflation int = 100

	opWeightMsgUpdateTimeBasedInflation = "op_weight_msg_time_based_inflation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateTimeBasedInflation int = 100

	opWeightMsgDeleteTimeBasedInflation = "op_weight_msg_time_based_inflation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteTimeBasedInflation int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenomicsGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		AirdropList: []types.Airdrop{
			{
				Authority: sample.AccAddress(),
				Intent:    "0",
			},
			{
				Authority: sample.AccAddress(),
				Intent:    "1",
			},
		},
		TimeBasedInflationList: []types.TimeBasedInflation{
			{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			},
			{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 1,
				EndBlockHeight:   1,
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenomicsGenesis)
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

	var weightMsgCreateAirdrop int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateAirdrop, &weightMsgCreateAirdrop, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAirdrop = defaultWeightMsgCreateAirdrop
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAirdrop,
		tokenomicssimulation.SimulateMsgCreateAirdrop(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateAirdrop int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAirdrop, &weightMsgUpdateAirdrop, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAirdrop = defaultWeightMsgUpdateAirdrop
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAirdrop,
		tokenomicssimulation.SimulateMsgUpdateAirdrop(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteAirdrop int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAirdrop, &weightMsgDeleteAirdrop, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAirdrop = defaultWeightMsgDeleteAirdrop
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAirdrop,
		tokenomicssimulation.SimulateMsgDeleteAirdrop(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateGenesisInflation int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateGenesisInflation, &weightMsgCreateGenesisInflation, nil,
		func(_ *rand.Rand) {
			weightMsgCreateGenesisInflation = defaultWeightMsgCreateGenesisInflation
		},
	)

	var weightMsgUpdateGenesisInflation int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateGenesisInflation, &weightMsgUpdateGenesisInflation, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateGenesisInflation = defaultWeightMsgUpdateGenesisInflation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateGenesisInflation,
		tokenomicssimulation.SimulateMsgUpdateGenesisInflation(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteGenesisInflation int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteGenesisInflation, &weightMsgDeleteGenesisInflation, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteGenesisInflation = defaultWeightMsgDeleteGenesisInflation
		},
	)

	var weightMsgCreateTimeBasedInflation int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTimeBasedInflation, &weightMsgCreateTimeBasedInflation, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTimeBasedInflation = defaultWeightMsgCreateTimeBasedInflation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTimeBasedInflation,
		tokenomicssimulation.SimulateMsgCreateTimeBasedInflation(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateTimeBasedInflation int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTimeBasedInflation, &weightMsgUpdateTimeBasedInflation, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTimeBasedInflation = defaultWeightMsgUpdateTimeBasedInflation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTimeBasedInflation,
		tokenomicssimulation.SimulateMsgUpdateTimeBasedInflation(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteTimeBasedInflation int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteTimeBasedInflation, &weightMsgDeleteTimeBasedInflation, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteTimeBasedInflation = defaultWeightMsgDeleteTimeBasedInflation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteTimeBasedInflation,
		tokenomicssimulation.SimulateMsgDeleteTimeBasedInflation(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
