package membershiptier

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/v7/testutil/sample"
	membershiptiersimulation "github.com/elys-network/elys/v7/x/tier/simulation"
	"github.com/elys-network/elys/v7/x/tier/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = membershiptiersimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgSetPortfolio = "op_weight_msg_set_portfolio"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetPortfolio int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	membershiptierGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortfolioList: []types.Portfolio{
			{
				Creator: sample.AccAddress(),
			},
			{
				Creator: sample.AccAddress(),
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&membershiptierGenesis)
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
	var weightMsgSetPortfolio int
	simState.AppParams.GetOrGenerate(opWeightMsgSetPortfolio, &weightMsgSetPortfolio, nil,
		func(_ *rand.Rand) {
			weightMsgSetPortfolio = defaultWeightMsgSetPortfolio
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetPortfolio,
		membershiptiersimulation.SimulateMsgSetPortfolio(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgSetPortfolio,
			defaultWeightMsgSetPortfolio,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				membershiptiersimulation.SimulateMsgSetPortfolio(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
