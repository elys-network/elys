package tradeshield

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	tradeshieldsimulation "github.com/elys-network/elys/x/tradeshield/simulation"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = tradeshieldsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgCreatePendingSpotOrder = "op_weight_msg_pending_spot_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePendingSpotOrder int = 100

	opWeightMsgUpdatePendingSpotOrder = "op_weight_msg_pending_spot_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePendingSpotOrder int = 100

	opWeightMsgCreatePendingPerpetualOrder = "op_weight_msg_pending_perpetual_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePendingPerpetualOrder int = 100

	opWeightMsgUpdatePendingPerpetualOrder = "op_weight_msg_pending_perpetual_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePendingPerpetualOrder int = 100

	opWeightMsgCancelPerpetualOrders = "op_weight_msg_pending_perpetual_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCancelPerpetualOrders int = 100

	opWeightMsgUpdateParams = "op_weight_msg_update_params"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateParams int = 100

	opWeightMsgExecuteOrders = "op_weight_msg_execute_orders"
	// TODO: Determine the simulation weight value
	defaultWeightMsgExecuteOrders int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tradeshieldGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PendingSpotOrderList: []types.SpotOrder{
			{
				OrderId:      0,
				OwnerAddress: sample.AccAddress(),
			},
			{
				OrderId:      1,
				OwnerAddress: sample.AccAddress(),
			},
		},
		PendingSpotOrderCount: 2,
		PendingPerpetualOrderList: []types.PerpetualOrder{
			{
				OrderId:      0,
				OwnerAddress: sample.AccAddress(),
			},
			{
				OrderId:      1,
				OwnerAddress: sample.AccAddress(),
			},
		},
		PendingPerpetualOrderCount: 2,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tradeshieldGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreatePendingSpotOrder int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePendingSpotOrder, &weightMsgCreatePendingSpotOrder, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePendingSpotOrder = defaultWeightMsgCreatePendingSpotOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePendingSpotOrder,
		tradeshieldsimulation.SimulateMsgCreatePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePendingSpotOrder int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePendingSpotOrder, &weightMsgUpdatePendingSpotOrder, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePendingSpotOrder = defaultWeightMsgUpdatePendingSpotOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePendingSpotOrder,
		tradeshieldsimulation.SimulateMsgUpdatePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreatePendingPerpetualOrder int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePendingPerpetualOrder, &weightMsgCreatePendingPerpetualOrder, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePendingPerpetualOrder = defaultWeightMsgCreatePendingPerpetualOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePendingPerpetualOrder,
		tradeshieldsimulation.SimulateMsgCreatePendingPerpetualOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePendingPerpetualOrder int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePendingPerpetualOrder, &weightMsgUpdatePendingPerpetualOrder, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePendingPerpetualOrder = defaultWeightMsgUpdatePendingPerpetualOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePendingPerpetualOrder,
		tradeshieldsimulation.SimulateMsgUpdatePendingPerpetualOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCancelPerpetualOrders int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelPerpetualOrders, &weightMsgCancelPerpetualOrders, nil,
		func(_ *rand.Rand) {
			weightMsgCancelPerpetualOrders = defaultWeightMsgCancelPerpetualOrders
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelPerpetualOrders,
		tradeshieldsimulation.SimulateMsgCancelPerpetualOrders(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateParams int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateParams, &weightMsgUpdateParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateParams = defaultWeightMsgUpdateParams
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateParams,
		tradeshieldsimulation.SimulateMsgUpdateParams(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgExecuteOrders int
	simState.AppParams.GetOrGenerate(opWeightMsgExecuteOrders, &weightMsgExecuteOrders, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteOrders = defaultWeightMsgExecuteOrders
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExecuteOrders,
		tradeshieldsimulation.SimulateMsgExecuteOrders(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreatePendingSpotOrder,
			defaultWeightMsgCreatePendingSpotOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgCreatePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdatePendingSpotOrder,
			defaultWeightMsgUpdatePendingSpotOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgUpdatePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreatePendingPerpetualOrder,
			defaultWeightMsgCreatePendingPerpetualOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgCreatePendingPerpetualOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdatePendingPerpetualOrder,
			defaultWeightMsgUpdatePendingPerpetualOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgUpdatePendingPerpetualOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCancelPerpetualOrders,
			defaultWeightMsgCancelPerpetualOrders,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgCancelPerpetualOrders(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateParams,
			defaultWeightMsgUpdateParams,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgUpdateParams(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgExecuteOrders,
			defaultWeightMsgExecuteOrders,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgExecuteOrders(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
