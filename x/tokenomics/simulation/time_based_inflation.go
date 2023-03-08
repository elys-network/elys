package simulation

import (
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "TimeBasedInflation already exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "timeBasedInflation authority not found"), nil, nil
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
			MsgType:         msg.Type(),
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "timeBasedInflation authority not found"), nil, nil
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
			MsgType:         msg.Type(),
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
