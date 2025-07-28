package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/elys-network/elys/v7/x/tradeshield/keeper"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

func SimulateMsgCancelPerpetualOrder(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		//simAccount, _ := simtypes.RandomAcc(r, accs)
		//msg := &types.MsgCancelPerpetualOrder{
		//	OwnerAddress: simAccount.Address.String(),
		//}

		// TODO: Handling the CancelPerpetualOrder simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgCancelPerpetualOrder{}), "CancelPerpetualOrder simulation not implemented"), nil, nil
	}
}
