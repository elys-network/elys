package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/elys-network/elys/x/parameter/keeper"
	"github.com/elys-network/elys/x/parameter/types"
)

func SimulateMsgUpdateWasmConfig(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUpdateWasmConfig{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the UpdateWasmConfig simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "UpdateWasmConfig simulation not implemented"), nil, nil
	}
}
