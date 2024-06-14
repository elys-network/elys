package parameter

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	parametersimulation "github.com/elys-network/elys/x/parameter/simulation"
	"github.com/elys-network/elys/x/parameter/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = parametersimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgUpdateWasmConfig = "op_weight_msg_update_wasm_config"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateWasmConfig int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	parameterGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&parameterGenesis)
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

	/*var weightMsgUpdateWasmConfig int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateWasmConfig, &weightMsgUpdateWasmConfig, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateWasmConfig = defaultWeightMsgUpdateWasmConfig
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateWasmConfig,
		parametersimulation.SimulateMsgUpdateWasmConfig(am.accountKeeper, am.bankKeeper, am.keeper),
	))*/

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
