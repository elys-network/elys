package simulation

import (
	"math/rand"

	simappparams "cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/v5/x/tokenomics/keeper"
	"github.com/elys-network/elys/v5/x/tokenomics/types"
)

func SimulateMsgCreateTimeBasedInflation(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		i := r.Int()
		msg := &types.MsgCreateTimeBasedInflation{
			Authority:        simAccount.Address.String(),
			StartBlockHeight: uint64(i),
			EndBlockHeight:   uint64(i),
		}

		_, found := k.GetTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgCreateTimeBasedInflation{}), "TimeBasedInflation already exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgUpdateTimeBasedInflation(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount            = simtypes.Account{}
			timeBasedInflation    = types.TimeBasedInflation{}
			msg                   = &types.MsgUpdateTimeBasedInflation{}
			allTimeBasedInflation = k.GetAllTimeBasedInflation(ctx)
			found                 = false
		)
		for _, obj := range allTimeBasedInflation {
			simAccount, found = FindAccount(accs, obj.Authority)
			if found {
				timeBasedInflation = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgUpdateTimeBasedInflation{}), "timeBasedInflation authority not found"), nil, nil
		}
		msg.Authority = simAccount.Address.String()

		msg.StartBlockHeight = timeBasedInflation.StartBlockHeight
		msg.EndBlockHeight = timeBasedInflation.EndBlockHeight

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgDeleteTimeBasedInflation(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount            = simtypes.Account{}
			timeBasedInflation    = types.TimeBasedInflation{}
			msg                   = &types.MsgUpdateTimeBasedInflation{}
			allTimeBasedInflation = k.GetAllTimeBasedInflation(ctx)
			found                 = false
		)
		for _, obj := range allTimeBasedInflation {
			simAccount, found = FindAccount(accs, obj.Authority)
			if found {
				timeBasedInflation = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgUpdateTimeBasedInflation{}), "timeBasedInflation authority not found"), nil, nil
		}
		msg.Authority = simAccount.Address.String()

		msg.StartBlockHeight = timeBasedInflation.StartBlockHeight
		msg.EndBlockHeight = timeBasedInflation.EndBlockHeight

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
