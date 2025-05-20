package amm

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/v4/testutil/sample"
	ammsimulation "github.com/elys-network/elys/v4/x/amm/simulation"
	"github.com/elys-network/elys/v4/x/amm/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = ammsimulation.FindAccount
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

	opWeightMsgSwapByDenom = "op_weight_msg_swap_by_denom"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSwapByDenom int = 100

	opWeightMsgUpdatePoolParams = "op_weight_msg_update_pool_params"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePoolParams int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
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
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(registry simtypes.StoreDecoderRegistry) {
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreatePool int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = defaultWeightMsgCreatePool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePool,
		ammsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgJoinPool int
	simState.AppParams.GetOrGenerate(opWeightMsgJoinPool, &weightMsgJoinPool, nil,
		func(_ *rand.Rand) {
			weightMsgJoinPool = defaultWeightMsgJoinPool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgJoinPool,
		ammsimulation.SimulateMsgJoinPool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgExitPool int
	simState.AppParams.GetOrGenerate(opWeightMsgExitPool, &weightMsgExitPool, nil,
		func(_ *rand.Rand) {
			weightMsgExitPool = defaultWeightMsgExitPool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExitPool,
		ammsimulation.SimulateMsgExitPool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapExactAmountIn int
	simState.AppParams.GetOrGenerate(opWeightMsgSwapExactAmountIn, &weightMsgSwapExactAmountIn, nil,
		func(_ *rand.Rand) {
			weightMsgSwapExactAmountIn = defaultWeightMsgSwapExactAmountIn
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapExactAmountIn,
		ammsimulation.SimulateMsgSwapExactAmountIn(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapExactAmountOut int
	simState.AppParams.GetOrGenerate(opWeightMsgSwapExactAmountOut, &weightMsgSwapExactAmountOut, nil,
		func(_ *rand.Rand) {
			weightMsgSwapExactAmountOut = defaultWeightMsgSwapExactAmountOut
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapExactAmountOut,
		ammsimulation.SimulateMsgSwapExactAmountOut(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapByDenom int
	simState.AppParams.GetOrGenerate(opWeightMsgSwapByDenom, &weightMsgSwapByDenom, nil,
		func(_ *rand.Rand) {
			weightMsgSwapByDenom = defaultWeightMsgSwapByDenom
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapByDenom,
		ammsimulation.SimulateMsgSwapByDenom(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePoolParams int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePoolParams, &weightMsgUpdatePoolParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePoolParams = defaultWeightMsgUpdatePoolParams
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePoolParams,
		ammsimulation.SimulateMsgUpdatePoolParams(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
