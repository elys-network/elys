package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v7/app"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	paramtypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetGetPositionCounter(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	leveragelp := app.LeveragelpKeeper

	simapp.SetStakingParam(app, ctx)
	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

	for i := 0; i < 2; i++ {
		position := types.Position{
			Address:        addr[i].String(),
			Collateral:     sdk.NewCoin(paramtypes.BaseCurrency, math.NewInt(0)),
			Liabilities:    math.NewInt(0),
			AmmPoolId:      1,
			PositionHealth: math.LegacyNewDec(0),
			Id:             0,
		}
		leveragelp.SetPosition(ctx, &position)
	}

	positionCount := leveragelp.GetPositionCounter(ctx, 1)
	require.Equal(t, positionCount.TotalOpen, (uint64)(2))
}
