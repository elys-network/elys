package stablestake

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/v5/testutil/sample"
	stablestakesimulation "github.com/elys-network/elys/v5/x/stablestake/simulation"
	"github.com/elys-network/elys/v5/x/stablestake/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = stablestakesimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgBond = "op_weight_msg_stake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBond int = 100

	opWeightMsgUnbond = "op_weight_msg_unbond"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnbond int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	stablestakeGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&stablestakeGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgBond int
	simState.AppParams.GetOrGenerate(opWeightMsgBond, &weightMsgBond, nil,
		func(_ *rand.Rand) {
			weightMsgBond = defaultWeightMsgBond
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBond,
		stablestakesimulation.SimulateMsgBond(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnbond int
	simState.AppParams.GetOrGenerate(opWeightMsgUnbond, &weightMsgUnbond, nil,
		func(_ *rand.Rand) {
			weightMsgUnbond = defaultWeightMsgUnbond
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnbond,
		stablestakesimulation.SimulateMsgUnbond(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgBond,
			defaultWeightMsgBond,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				stablestakesimulation.SimulateMsgBond(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUnbond,
			defaultWeightMsgUnbond,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				stablestakesimulation.SimulateMsgUnbond(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
