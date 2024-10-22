package keeper

import (
	simapp "github.com/elys-network/elys/app"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
)

func PerpetualKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	app := simapp.InitElysTestApp(true)
	baseCtx := app.BaseApp.NewContext(true, tmproto.Header{})
	return app.PerpetualKeeper, baseCtx
}
