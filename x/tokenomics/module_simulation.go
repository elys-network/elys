package tokenomics

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	_ = simappparams.StakePerAccount
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
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenomicsGenesis)
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

	var weightMsgCreateAirdrop int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateAirdrop, &weightMsgCreateAirdrop, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAirdrop = defaultWeightMsgCreateAirdrop
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAirdrop,
		tokenomicssimulation.SimulateMsgCreateAirdrop(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateAirdrop int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateAirdrop, &weightMsgUpdateAirdrop, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAirdrop = defaultWeightMsgUpdateAirdrop
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAirdrop,
		tokenomicssimulation.SimulateMsgUpdateAirdrop(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteAirdrop int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteAirdrop, &weightMsgDeleteAirdrop, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAirdrop = defaultWeightMsgDeleteAirdrop
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAirdrop,
		tokenomicssimulation.SimulateMsgDeleteAirdrop(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateGenesisInflation int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateGenesisInflation, &weightMsgCreateGenesisInflation, nil,
		func(_ *rand.Rand) {
			weightMsgCreateGenesisInflation = defaultWeightMsgCreateGenesisInflation
		},
	)

	var weightMsgUpdateGenesisInflation int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateGenesisInflation, &weightMsgUpdateGenesisInflation, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateGenesisInflation = defaultWeightMsgUpdateGenesisInflation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateGenesisInflation,
		tokenomicssimulation.SimulateMsgUpdateGenesisInflation(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteGenesisInflation int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteGenesisInflation, &weightMsgDeleteGenesisInflation, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteGenesisInflation = defaultWeightMsgDeleteGenesisInflation
		},
	)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
