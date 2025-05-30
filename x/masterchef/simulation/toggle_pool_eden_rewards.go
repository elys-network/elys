package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/elys-network/elys/v6/x/masterchef/keeper"
	"github.com/elys-network/elys/v6/x/masterchef/types"
)

func SimulateMsgTogglePoolEdenRewards(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		//simAccount, _ := simtypes.RandomAcc(r, accs)
		//msg := &types.MsgTogglePoolEdenRewards{
		//	Authority: simAccount.Address.String(),
		//}

		// TODO: Handling the TogglePoolEdenRewards simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgTogglePoolEdenRewards{}), "TogglePoolEdenRewards simulation not implemented"), nil, nil
	}
}
