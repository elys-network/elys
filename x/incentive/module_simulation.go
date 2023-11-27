package incentive

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/testutil/sample"
	incentivesimulation "github.com/elys-network/elys/x/incentive/simulation"
	"github.com/elys-network/elys/x/incentive/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = incentivesimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgSetWithdrawAddress = "op_weight_msg_set_withdraw_address"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetWithdrawAddress int = 100

	opWeightMsgWithdrawValidatorCommission = "op_weight_msg_withdraw_validator_commission"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWithdrawValidatorCommission int = 100

	opWeightMsgWithdrawDelegatorReward = "op_weight_msg_withdraw_delegator_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWithdrawDelegatorReward int = 100

	opWeightMsgUpdateIncentiveParams = "op_weight_msg_update_incentive_params"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateIncentiveParams int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	incentiveGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&incentiveGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgSetWithdrawAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetWithdrawAddress, &weightMsgSetWithdrawAddress, nil,
		func(_ *rand.Rand) {
			weightMsgSetWithdrawAddress = defaultWeightMsgSetWithdrawAddress
		},
	)

	var weightMsgWithdrawValidatorCommission int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWithdrawValidatorCommission, &weightMsgWithdrawValidatorCommission, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawValidatorCommission = defaultWeightMsgWithdrawValidatorCommission
		},
	)

	var weightMsgWithdrawDelegatorReward int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWithdrawDelegatorReward, &weightMsgWithdrawDelegatorReward, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawDelegatorReward = defaultWeightMsgWithdrawDelegatorReward
		},
	)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
