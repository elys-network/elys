package keeper

import (
	simapp "github.com/elys-network/elys/app"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
)

func PerpetualKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	app := simapp.InitElysTestApp(true, t)
	baseCtx := app.BaseApp.NewContext(true)
	return app.PerpetualKeeper, baseCtx
}
