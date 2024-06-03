package membershiptier

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	membershiptiersimulation "github.com/elys-network/elys/x/membershiptier/simulation"
	"github.com/elys-network/elys/x/membershiptier/types"
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
	opWeightMsgCreatePortfolio = "op_weight_msg_portfolio"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePortfolio int = 100

	opWeightMsgUpdatePortfolio = "op_weight_msg_portfolio"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePortfolio int = 100

	opWeightMsgDeletePortfolio = "op_weight_msg_portfolio"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePortfolio int = 100

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
				Creator:  sample.AccAddress(),
				Assetkey: types.LiquidKeyPrefix,
			},
			{
				Creator:  sample.AccAddress(),
				Assetkey: types.LiquidKeyPrefix,
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&membershiptierGenesis)
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

	var weightMsgCreatePortfolio int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePortfolio, &weightMsgCreatePortfolio, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePortfolio = defaultWeightMsgCreatePortfolio
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePortfolio,
		membershiptiersimulation.SimulateMsgCreatePortfolio(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePortfolio int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePortfolio, &weightMsgUpdatePortfolio, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePortfolio = defaultWeightMsgUpdatePortfolio
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePortfolio,
		membershiptiersimulation.SimulateMsgUpdatePortfolio(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePortfolio int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeletePortfolio, &weightMsgDeletePortfolio, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePortfolio = defaultWeightMsgDeletePortfolio
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePortfolio,
		membershiptiersimulation.SimulateMsgDeletePortfolio(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreatePortfolio,
			defaultWeightMsgCreatePortfolio,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				membershiptiersimulation.SimulateMsgCreatePortfolio(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdatePortfolio,
			defaultWeightMsgUpdatePortfolio,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				membershiptiersimulation.SimulateMsgUpdatePortfolio(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeletePortfolio,
			defaultWeightMsgDeletePortfolio,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				membershiptiersimulation.SimulateMsgDeletePortfolio(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
