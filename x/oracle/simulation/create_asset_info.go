package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/elys-network/elys/v4/x/oracle/keeper"
	"github.com/elys-network/elys/v4/x/oracle/types"
)

func SimulateMsgCreateAssetInfo(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		//simAccount, _ := simtypes.RandomAcc(r, accs)
		//msg := &types.MsgCreateAssetInfo{
		//	Creator: simAccount.Address.String(),
		//}

		// TODO: Handling the CreateAssetInfo simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgCreateAssetInfo{}), "CreateAssetInfo simulation not implemented"), nil, nil
	}
}
