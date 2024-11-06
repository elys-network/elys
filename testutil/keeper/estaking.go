package keeper

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/estaking/keeper"
)

func EstakingKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, sdk.AccAddress, sdk.ValAddress) {
	app, genAccount, valAddr := simapp.InitElysTestAppWithGenAccount()
	baseCtx := app.BaseApp.NewContext(true, tmproto.Header{})
	return &app.EstakingKeeper, baseCtx, genAccount, valAddr
}
