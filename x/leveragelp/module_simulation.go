package leveragelp

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	leveragelpsimulation "github.com/elys-network/elys/x/leveragelp/simulation"
	"github.com/elys-network/elys/x/leveragelp/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = leveragelpsimulation.FindAccount
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

	opWeightMsgUpdatePools = "op_weight_msg_update_pools"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePools int = 100

	opWeightMsgWhitelist = "op_weight_msg_whitelist"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWhitelist int = 100

	opWeightMsgDewhitelist = "op_weight_msg_dewhitelist"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDewhitelist int = 100

	opWeightMsgUpdateStopLoss = "op_weight_msg_update_stop_loss"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateStopLoss int = 100

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
	leveragelpGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&leveragelpGenesis)
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

	var weightMsgOpen int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgOpen, &weightMsgOpen, nil,
		func(_ *rand.Rand) {
			weightMsgOpen = defaultWeightMsgOpen
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgOpen,
		leveragelpsimulation.SimulateMsgOpen(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClosep int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClosep, &weightMsgClosep, nil,
		func(_ *rand.Rand) {
			weightMsgClosep = defaultWeightMsgClosep
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClosep,
		leveragelpsimulation.SimulateMsgClosep(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateParams int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateParams, &weightMsgUpdateParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateParams = defaultWeightMsgUpdateParams
		},
	)

	var weightMsgUpdatePools int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePools, &weightMsgUpdatePools, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePools = defaultWeightMsgUpdatePools
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePools,
		leveragelpsimulation.SimulateMsgUpdatePools(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgWhitelist int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWhitelist, &weightMsgWhitelist, nil,
		func(_ *rand.Rand) {
			weightMsgWhitelist = defaultWeightMsgWhitelist
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgWhitelist,
		leveragelpsimulation.SimulateMsgWhitelist(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDewhitelist int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDewhitelist, &weightMsgDewhitelist, nil,
		func(_ *rand.Rand) {
			weightMsgDewhitelist = defaultWeightMsgDewhitelist
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDewhitelist,
		leveragelpsimulation.SimulateMsgDewhitelist(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateStopLoss int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateStopLoss, &weightMsgUpdateStopLoss, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateStopLoss = defaultWeightMsgUpdateStopLoss
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateStopLoss,
		leveragelpsimulation.SimulateMsgUpdateStopLoss(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddCollateral int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddCollateral, &weightMsgAddCollateral, nil,
		func(_ *rand.Rand) {
			weightMsgAddCollateral = defaultWeightMsgAddCollateral
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddCollateral,
		leveragelpsimulation.SimulateMsgAddCollateral(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClosePositions int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClosePositions, &weightMsgClosePositions, nil,
		func(_ *rand.Rand) {
			weightMsgClosePositions = defaultWeightMsgClosePositions
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClosePositions,
		leveragelpsimulation.SimulateMsgClosePositions(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
