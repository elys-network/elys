package oracle

import (
	"math/rand"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/v5/testutil/sample"
	oraclesimulation "github.com/elys-network/elys/v5/x/oracle/simulation"
	"github.com/elys-network/elys/v5/x/oracle/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = oraclesimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgFeedPrice = "op_weight_msg_price"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFeedPrice int = 100

	opWeightMsgDeletePriceFeeder = "op_weight_msg_price_feeder"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePriceFeeder int = 100

	opWeightMsgFeedMultiplePrices = "op_weight_msg_feed_multiple_prices"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFeedMultiplePrices int = 100

	opWeightMsgCreateAssetInfo = "op_weight_msg_create_asset_info"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateAssetInfo int = 100

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
		AssetInfos: []types.AssetInfo{
			{
				Denom: "satoshi",
			},
			{
				Denom: "wei",
			},
		},
		Prices: []types.Price{
			{
				Provider: sample.AccAddress(),
				Asset:    "BTC",
				Price:    sdkmath.LegacyZeroDec(),
			},
			{
				Provider: sample.AccAddress(),
				Asset:    "BTC",
				Price:    sdkmath.LegacyOneDec(),
			},
		},
		PriceFeeders: []types.PriceFeeder{
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

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgFeedPrice int
	simState.AppParams.GetOrGenerate(opWeightMsgFeedPrice, &weightMsgFeedPrice, nil,
		func(_ *rand.Rand) {
			weightMsgFeedPrice = defaultWeightMsgFeedPrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFeedPrice,
		oraclesimulation.SimulateMsgFeedPrice(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePriceFeeder int
	simState.AppParams.GetOrGenerate(opWeightMsgDeletePriceFeeder, &weightMsgDeletePriceFeeder, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePriceFeeder = defaultWeightMsgDeletePriceFeeder
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePriceFeeder,
		oraclesimulation.SimulateMsgDeletePriceFeeder(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgFeedMultiplePrices int
	simState.AppParams.GetOrGenerate(opWeightMsgFeedMultiplePrices, &weightMsgFeedMultiplePrices, nil,
		func(_ *rand.Rand) {
			weightMsgFeedMultiplePrices = defaultWeightMsgFeedMultiplePrices
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFeedMultiplePrices,
		oraclesimulation.SimulateMsgFeedMultiplePrices(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateAssetInfo int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateAssetInfo, &weightMsgCreateAssetInfo, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAssetInfo = defaultWeightMsgCreateAssetInfo
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAssetInfo,
		oraclesimulation.SimulateMsgCreateAssetInfo(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
