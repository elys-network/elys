package masterchef

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/v4/testutil/sample"
	masterchefsimulation "github.com/elys-network/elys/v4/x/masterchef/simulation"
	"github.com/elys-network/elys/v4/x/masterchef/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = masterchefsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgTogglePoolEdenRewards = "op_weight_msg_toggle_pool_eden_rewards"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTogglePoolEdenRewards int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	masterchefGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&masterchefGenesis)
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

	var weightMsgTogglePoolEdenRewards int
	simState.AppParams.GetOrGenerate(opWeightMsgTogglePoolEdenRewards, &weightMsgTogglePoolEdenRewards, nil,
		func(_ *rand.Rand) {
			weightMsgTogglePoolEdenRewards = defaultWeightMsgTogglePoolEdenRewards
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTogglePoolEdenRewards,
		masterchefsimulation.SimulateMsgTogglePoolEdenRewards(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgTogglePoolEdenRewards,
			defaultWeightMsgTogglePoolEdenRewards,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				masterchefsimulation.SimulateMsgTogglePoolEdenRewards(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
