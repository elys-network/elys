package amm

import (
	"math/rand"

	simappparams "cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	ammsimulation "github.com/elys-network/elys/x/amm/simulation"
	"github.com/elys-network/elys/x/amm/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = ammsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreatePool = "op_weight_msg_create_pool"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePool int = 100

	opWeightMsgJoinPool = "op_weight_msg_join_pool"
	// TODO: Determine the simulation weight value
	defaultWeightMsgJoinPool int = 100

	opWeightMsgExitPool = "op_weight_msg_exit_pool"
	// TODO: Determine the simulation weight value
	defaultWeightMsgExitPool int = 100

	opWeightMsgSwapExactAmountIn = "op_weight_msg_swap_exact_amount_in"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSwapExactAmountIn int = 100

	opWeightMsgSwapExactAmountOut = "op_weight_msg_swap_exact_amount_out"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSwapExactAmountOut int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	ammGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&ammGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {
	ammParams := types.DefaultParams()
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyPoolCreationFee), func(r *rand.Rand) string {
			return string(types.Amino.MustMarshalJSON(ammParams.PoolCreationFee))
		}),
	}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreatePool int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = defaultWeightMsgCreatePool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePool,
		ammsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgJoinPool int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgJoinPool, &weightMsgJoinPool, nil,
		func(_ *rand.Rand) {
			weightMsgJoinPool = defaultWeightMsgJoinPool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgJoinPool,
		ammsimulation.SimulateMsgJoinPool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgExitPool int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgExitPool, &weightMsgExitPool, nil,
		func(_ *rand.Rand) {
			weightMsgExitPool = defaultWeightMsgExitPool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExitPool,
		ammsimulation.SimulateMsgExitPool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapExactAmountIn int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSwapExactAmountIn, &weightMsgSwapExactAmountIn, nil,
		func(_ *rand.Rand) {
			weightMsgSwapExactAmountIn = defaultWeightMsgSwapExactAmountIn
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapExactAmountIn,
		ammsimulation.SimulateMsgSwapExactAmountIn(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapExactAmountOut int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSwapExactAmountOut, &weightMsgSwapExactAmountOut, nil,
		func(_ *rand.Rand) {
			weightMsgSwapExactAmountOut = defaultWeightMsgSwapExactAmountOut
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapExactAmountOut,
		ammsimulation.SimulateMsgSwapExactAmountOut(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
