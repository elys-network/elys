package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/elys-network/elys/v7/x/perpetual/keeper"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func SimulateMsgWhitelist(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		//simAccount, _ := simtypes.RandomAcc(r, accs)
		//msg := &types.MsgWhitelist{
		//	Authority: simAccount.Address.String(),
		//}

		// TODO: Handling the Whitelist simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgWhitelist{}), "Whitelist simulation not implemented"), nil, nil
	}
}
