package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/estaking/keeper"
)

func EstakingKeeper(t *testing.T) (*keeper.Keeper, sdk.Context, sdk.AccAddress, sdk.ValAddress) {
	app, genAccount, valAddr := simapp.InitElysTestAppWithGenAccount(t)
	baseCtx := app.BaseApp.NewContext(true)
	return app.EstakingKeeper, baseCtx, genAccount, valAddr
}
