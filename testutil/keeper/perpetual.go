package keeper

import (
	"testing"

	simapp "github.com/elys-network/elys/v7/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/keeper"
)

func PerpetualKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	app := simapp.InitElysTestApp(true, t)
	baseCtx := app.BaseApp.NewContext(true)
	return app.PerpetualKeeper, baseCtx
}
