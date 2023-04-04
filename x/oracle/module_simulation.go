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

	opWeightMsgCreatePrice = "op_weight_msg_price"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePrice int = 100

	opWeightMsgUpdatePrice = "op_weight_msg_price"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePrice int = 100

	opWeightMsgDeletePrice = "op_weight_msg_price"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePrice int = 100

	opWeightMsgCreatePriceFeeder = "op_weight_msg_price_feeder"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePriceFeeder int = 100

	opWeightMsgUpdatePriceFeeder = "op_weight_msg_price_feeder"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePriceFeeder int = 100

	opWeightMsgDeletePriceFeeder = "op_weight_msg_price_feeder"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePriceFeeder int = 100

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
		PriceList: []types.Price{
			{
				Provider: sample.AccAddress(),
				Asset:    "BTC",
				Price:    sdk.ZeroDec(),
			},
			{
				Provider: sample.AccAddress(),
				Asset:    "BTC",
				Price:    sdk.OneDec(),
			},
		},
		PriceFeederList: []types.PriceFeeder{
			{
				Feeder:   sample.AccAddress(),
				IsActive: true,
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

	var weightMsgCreatePrice int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePrice, &weightMsgCreatePrice, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePrice = defaultWeightMsgCreatePrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePrice,
		oraclesimulation.SimulateMsgCreatePrice(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePrice int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePrice, &weightMsgUpdatePrice, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePrice = defaultWeightMsgUpdatePrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePrice,
		oraclesimulation.SimulateMsgUpdatePrice(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePrice int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeletePrice, &weightMsgDeletePrice, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePrice = defaultWeightMsgDeletePrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePrice,
		oraclesimulation.SimulateMsgDeletePrice(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreatePriceFeeder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePriceFeeder, &weightMsgCreatePriceFeeder, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePriceFeeder = defaultWeightMsgCreatePriceFeeder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePriceFeeder,
		oraclesimulation.SimulateMsgCreatePriceFeeder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePriceFeeder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePriceFeeder, &weightMsgUpdatePriceFeeder, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePriceFeeder = defaultWeightMsgUpdatePriceFeeder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePriceFeeder,
		oraclesimulation.SimulateMsgUpdatePriceFeeder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePriceFeeder int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeletePriceFeeder, &weightMsgDeletePriceFeeder, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePriceFeeder = defaultWeightMsgDeletePriceFeeder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePriceFeeder,
		oraclesimulation.SimulateMsgDeletePriceFeeder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
