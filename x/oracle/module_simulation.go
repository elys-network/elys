package oracle

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	oraclesimulation "github.com/elys-network/elys/x/oracle/simulation"
	"github.com/elys-network/elys/x/oracle/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = oraclesimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateAssetInfo = "op_weight_msg_asset_info"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateAssetInfo int = 100

	opWeightMsgUpdateAssetInfo = "op_weight_msg_asset_info"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateAssetInfo int = 100

	opWeightMsgDeleteAssetInfo = "op_weight_msg_asset_info"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAssetInfo int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	oracleGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		AssetInfoList: []types.AssetInfo{
			{
				Denom: "satoshi",
			},
			{
				Denom: "wei",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&oracleGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateAssetInfo int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateAssetInfo, &weightMsgCreateAssetInfo, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAssetInfo = defaultWeightMsgCreateAssetInfo
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAssetInfo,
		oraclesimulation.SimulateMsgCreateAssetInfo(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateAssetInfo int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateAssetInfo, &weightMsgUpdateAssetInfo, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAssetInfo = defaultWeightMsgUpdateAssetInfo
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAssetInfo,
		oraclesimulation.SimulateMsgUpdateAssetInfo(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteAssetInfo int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteAssetInfo, &weightMsgDeleteAssetInfo, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAssetInfo = defaultWeightMsgDeleteAssetInfo
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAssetInfo,
		oraclesimulation.SimulateMsgDeleteAssetInfo(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
