package burner

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	burnersimulation "github.com/elys-network/elys/x/burner/simulation"
	"github.com/elys-network/elys/x/burner/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = burnersimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgMsgUpdateParams = "op_weight_msg_msg_update_params"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMsgUpdateParams int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	burnerGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&burnerGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgMsgUpdateParams int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMsgUpdateParams, &weightMsgMsgUpdateParams, nil,
		func(_ *rand.Rand) {
			weightMsgMsgUpdateParams = defaultWeightMsgMsgUpdateParams
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMsgUpdateParams,
		burnersimulation.SimulateMsgUpdateParams(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
