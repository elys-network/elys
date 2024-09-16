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

	opWeightMsgDeletePendingSpotOrder = "op_weight_msg_pending_spot_order"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePendingSpotOrder int = 100

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
		PendingSpotOrderList: []types.PendingSpotOrder{
			{
				Id:      0,
				Creator: sample.AccAddress(),
			},
			{
				Id:      1,
				Creator: sample.AccAddress(),
			},
		},
		PendingSpotOrderCount: 2,
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

	var weightMsgCreatePendingSpotOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePendingSpotOrder, &weightMsgCreatePendingSpotOrder, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePendingSpotOrder = defaultWeightMsgCreatePendingSpotOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePendingSpotOrder,
		tradeshieldsimulation.SimulateMsgCreatePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePendingSpotOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePendingSpotOrder, &weightMsgUpdatePendingSpotOrder, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePendingSpotOrder = defaultWeightMsgUpdatePendingSpotOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePendingSpotOrder,
		tradeshieldsimulation.SimulateMsgUpdatePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePendingSpotOrder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeletePendingSpotOrder, &weightMsgDeletePendingSpotOrder, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePendingSpotOrder = defaultWeightMsgDeletePendingSpotOrder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePendingSpotOrder,
		tradeshieldsimulation.SimulateMsgDeletePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper),
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
			opWeightMsgDeletePendingSpotOrder,
			defaultWeightMsgDeletePendingSpotOrder,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tradeshieldsimulation.SimulateMsgDeletePendingSpotOrder(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
