package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/elys-network/elys/v6/x/perpetual/keeper"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func SimulateMsgOpen(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		//simAccount, _ := simtypes.RandomAcc(r, accs)
		//msg := &types.MsgOpen{
		//	Creator: simAccount.Address.String(),
		//}

		// TODO: Handling the Open simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgOpen{}), "Open simulation not implemented"), nil, nil
	}
}
