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
	opWeightMsgSetAssetInfo = "op_weight_msg_asset_info"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetAssetInfo int = 100

	opWeightMsgDeleteAssetInfo = "op_weight_msg_asset_info"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAssetInfo int = 100

	opWeightMsgFeedPrice = "op_weight_msg_price"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFeedPrice int = 100

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

	var weightMsgSetAssetInfo int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetAssetInfo, &weightMsgSetAssetInfo, nil,
		func(_ *rand.Rand) {
			weightMsgSetAssetInfo = defaultWeightMsgSetAssetInfo
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetAssetInfo,
		oraclesimulation.SimulateMsgSetAssetInfo(am.accountKeeper, am.bankKeeper, am.keeper),
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

	var weightMsgFeedPrice int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgFeedPrice, &weightMsgFeedPrice, nil,
		func(_ *rand.Rand) {
			weightMsgFeedPrice = defaultWeightMsgFeedPrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFeedPrice,
		oraclesimulation.SimulateMsgFeedPrice(am.accountKeeper, am.bankKeeper, am.keeper),
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
