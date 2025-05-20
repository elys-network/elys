package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v4/app"
	"github.com/elys-network/elys/v4/x/estaking/keeper"
)

func EstakingKeeper(t *testing.T) (*keeper.Keeper, sdk.Context, sdk.AccAddress, sdk.ValAddress) {
	app, genAccount, valAddr := simapp.InitElysTestAppWithGenAccount(t)
	baseCtx := app.BaseApp.NewContext(true)
	return app.EstakingKeeper, baseCtx, genAccount, valAddr
}
