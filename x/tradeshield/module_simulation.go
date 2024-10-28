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
	opWeightMsgCreateSpotOrder = "op_weight_msg_spot_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateSpotOrder int = 100

	opWeightMsgUpdateSpotOrder = "op_weight_msg_spot_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateSpotOrder int = 100

	opWeightMsgCreatePerpetualOpenOrder = "op_weight_msg_perpetual_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePerpetualOpenOrder int = 100

	opWeightMsgCreatePerpetualCloseOrder = "op_weight_msg_perpetual_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePerpetualCloseOrder int = 100

	opWeightMsgUpdatePerpetualOrder = "op_weight_msg_perpetual_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePerpetualOrder int = 100

	opWeightMsgCancelPerpetualOrders = "op_weight_msg_perpetual_order"
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
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateSpotOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateSpotOrder, &weightMsgCreateSpotOrder, nil,
		func(_ *rand.Rand) {
			weightMsgCreateSpotOrder = defaultWeightMsgCreateSpotOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateSpotOrder,
		tradeshieldsimulation.SimulateMsgCreateSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateSpotOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateSpotOrder, &weightMsgUpdateSpotOrder, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateSpotOrder = defaultWeightMsgUpdateSpotOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateSpotOrder,
		tradeshieldsimulation.SimulateMsgUpdateSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreatePerpetualOpenOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePerpetualOpenOrder, &weightMsgCreatePerpetualOpenOrder, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePerpetualOpenOrder = defaultWeightMsgCreatePerpetualOpenOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePerpetualOpenOrder,
		tradeshieldsimulation.SimulateMsgCreatePerpetualOpenOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreatePerpetualCloseOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePerpetualCloseOrder, &weightMsgCreatePerpetualCloseOrder, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePerpetualCloseOrder = defaultWeightMsgCreatePerpetualCloseOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePerpetualCloseOrder,
		tradeshieldsimulation.SimulateMsgCreatePerpetualCloseOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePerpetualOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePerpetualOrder, &weightMsgUpdatePerpetualOrder, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePerpetualOrder = defaultWeightMsgUpdatePerpetualOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePerpetualOrder,
		tradeshieldsimulation.SimulateMsgUpdatePerpetualOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCancelPerpetualOrders int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCancelPerpetualOrders, &weightMsgCancelPerpetualOrders, nil,
		func(_ *rand.Rand) {
			weightMsgCancelPerpetualOrders = defaultWeightMsgCancelPerpetualOrders
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelPerpetualOrders,
		tradeshieldsimulation.SimulateMsgCancelPerpetualOrders(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateParams int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateParams, &weightMsgUpdateParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateParams = defaultWeightMsgUpdateParams
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateParams,
		tradeshieldsimulation.SimulateMsgUpdateParams(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgExecuteOrders int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgExecuteOrders, &weightMsgExecuteOrders, nil,
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
			opWeightMsgCreateSpotOrder,
			defaultWeightMsgCreateSpotOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgCreateSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateSpotOrder,
			defaultWeightMsgUpdateSpotOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgUpdateSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreatePerpetualOpenOrder,
			defaultWeightMsgCreatePerpetualOpenOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgCreatePerpetualOpenOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdatePerpetualOrder,
			defaultWeightMsgUpdatePerpetualOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgUpdatePerpetualOrder(am.accountKeeper, am.bankKeeper, am.keeper)
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
