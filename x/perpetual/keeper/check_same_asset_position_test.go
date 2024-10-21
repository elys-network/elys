package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestCheckSameAssetPosition_NewPosition(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)

	k := app.PerpetualKeeper
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	mtp := types.NewMTP(addr[0].String(), ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency, ptypes.ATOM, types.Position_LONG, sdkmath.LegacyNewDec(5), sdkmath.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault), 1)
	err := k.SetMTP(ctx, mtp)
	assert.Nil(t, err)

	msg := &types.MsgOpen{
		Creator:       addr[0].String(),
		Position:      types.Position_SHORT,
		Leverage:      sdkmath.LegacyNewDec(1),
		TradingAsset:  ptypes.ATOM,
		Collateral:    sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100)),
		StopLossPrice: sdkmath.LegacyNewDec(100),
	}

	mtp = k.CheckSameAssetPosition(ctx, msg)

	// Expect no error
	assert.Nil(t, mtp)
}
