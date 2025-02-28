package assetprofile

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	assetprofilesimulation "github.com/elys-network/elys/x/assetprofile/simulation"
	"github.com/elys-network/elys/x/assetprofile/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = assetprofilesimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateEntry = "op_weight_msg_entry"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateEntry int = 100

	opWeightMsgUpdateEntry = "op_weight_msg_entry"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateEntry int = 100

	opWeightMsgDeleteEntry = "op_weight_msg_entry"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteEntry int = 100

	opWeightMsgAddEntry = "op_weight_msg_add_entry"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddEntry int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	assetprofileGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		EntryList: []types.Entry{
			{
				Authority: sample.AccAddress(),
				BaseDenom: "0",
			},
			{
				Authority: sample.AccAddress(),
				BaseDenom: "1",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&assetprofileGenesis)
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

	var weightMsgCreateEntry int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateEntry, &weightMsgCreateEntry, nil,
		func(_ *rand.Rand) {
			weightMsgCreateEntry = defaultWeightMsgCreateEntry
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateEntry,
		assetprofilesimulation.SimulateMsgCreateEntry(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateEntry int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateEntry, &weightMsgUpdateEntry, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEntry = defaultWeightMsgUpdateEntry
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEntry,
		assetprofilesimulation.SimulateMsgUpdateEntry(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteEntry int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteEntry, &weightMsgDeleteEntry, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteEntry = defaultWeightMsgDeleteEntry
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteEntry,
		assetprofilesimulation.SimulateMsgDeleteEntry(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddEntry int
	simState.AppParams.GetOrGenerate(opWeightMsgAddEntry, &weightMsgAddEntry, nil,
		func(_ *rand.Rand) {
			weightMsgAddEntry = defaultWeightMsgAddEntry
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddEntry,
		assetprofilesimulation.SimulateMsgAddEntry(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
