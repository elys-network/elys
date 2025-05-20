package perpetual

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/v4/testutil/sample"
	perpetualsimulation "github.com/elys-network/elys/v4/x/perpetual/simulation"
	"github.com/elys-network/elys/v4/x/perpetual/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = perpetualsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgOpen = "op_weight_msg_open"
	// TODO: Determine the simulation weight value
	defaultWeightMsgOpen int = 100

	opWeightMsgClosep = "op_weight_msg_closep"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClosep int = 100

	opWeightMsgUpdateParams = "op_weight_msg_update_params"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateParams int = 100

	opWeightMsgWhitelist = "op_weight_msg_whitelist"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWhitelist int = 100

	opWeightMsgDewhitelist = "op_weight_msg_dewhitelist"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDewhitelist int = 100

	opWeightMsgAddCollateral = "op_weight_msg_add_collateral"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddCollateral int = 100

	opWeightMsgClosePositions = "op_weight_msg_close_positions"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClosePositions int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	perpetualGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&perpetualGenesis)
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

	var weightMsgOpen int
	simState.AppParams.GetOrGenerate(opWeightMsgOpen, &weightMsgOpen, nil,
		func(_ *rand.Rand) {
			weightMsgOpen = defaultWeightMsgOpen
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgOpen,
		perpetualsimulation.SimulateMsgOpen(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClosep int
	simState.AppParams.GetOrGenerate(opWeightMsgClosep, &weightMsgClosep, nil,
		func(_ *rand.Rand) {
			weightMsgClosep = defaultWeightMsgClosep
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClosep,
		perpetualsimulation.SimulateMsgClosep(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateParams int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateParams, &weightMsgUpdateParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateParams = defaultWeightMsgUpdateParams
		},
	)

	var weightMsgWhitelist int
	simState.AppParams.GetOrGenerate(opWeightMsgWhitelist, &weightMsgWhitelist, nil,
		func(_ *rand.Rand) {
			weightMsgWhitelist = defaultWeightMsgWhitelist
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgWhitelist,
		perpetualsimulation.SimulateMsgWhitelist(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDewhitelist int
	simState.AppParams.GetOrGenerate(opWeightMsgDewhitelist, &weightMsgDewhitelist, nil,
		func(_ *rand.Rand) {
			weightMsgDewhitelist = defaultWeightMsgDewhitelist
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDewhitelist,
		perpetualsimulation.SimulateMsgDewhitelist(am.accountKeeper, am.bankKeeper, am.keeper),
	))
	var weightMsgClosePositions int
	simState.AppParams.GetOrGenerate(opWeightMsgClosePositions, &weightMsgClosePositions, nil,
		func(_ *rand.Rand) {
			weightMsgClosePositions = defaultWeightMsgClosePositions
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClosePositions,
		perpetualsimulation.SimulateMsgClosePositions(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
