package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v7/app"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	paramtypes "github.com/elys-network/elys/v7/x/parameter/types"
	stablestaketypes "github.com/elys-network/elys/v7/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func TestSetGetPosition(t *testing.T) {
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

	positionCount := leveragelp.GetPositionCount(ctx)
	require.Equal(t, positionCount, (uint64)(2))
}

func TestSetLiquidation(t *testing.T) {
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
	positionCount := leveragelp.GetPositionCount(ctx)
	require.Equal(t, positionCount, (uint64)(2))
}

func TestIteratePoolPosIdsLiquidationSorted(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	leveragelp := app.LeveragelpKeeper
	stablestake := app.StablestakeKeeper
	simapp.SetStakingParam(app, ctx)
	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

	positions := []struct {
		Id                uint64
		LeveragedLpAmount math.Int
		Borrowed          math.Int
		PoolId            uint64
	}{
		{
			LeveragedLpAmount: math.NewInt(1000),
			Id:                7,
			Borrowed:          math.NewInt(7),
			PoolId:            2,
		},
		{
			LeveragedLpAmount: math.NewInt(100),
			Id:                6,
			Borrowed:          math.NewInt(0),
			PoolId:            1,
		},
		{
			LeveragedLpAmount: math.NewInt(1000),
			Id:                5,
			Borrowed:          math.NewInt(7),
			PoolId:            1,
		},
		{
			LeveragedLpAmount: math.NewInt(1000),
			Id:                4,
			Borrowed:          math.NewInt(8),
			PoolId:            1,
		},
		{
			LeveragedLpAmount: math.NewInt(1000),
			Id:                2,
			Borrowed:          math.NewInt(100),
			PoolId:            1,
		},
		{
			LeveragedLpAmount: math.NewInt(100),
			Id:                1,
			Borrowed:          math.NewInt(10),
			PoolId:            1,
		},
		{
			LeveragedLpAmount: math.NewInt(1000),
			Id:                3,
			Borrowed:          math.NewInt(9),
			PoolId:            1,
		},
	}
	for _, info := range positions {
		position := types.Position{
			LeveragedLpAmount: info.LeveragedLpAmount,
			Id:                info.Id,
			Address:           addr[0].String(),
			Collateral:        sdk.NewCoin(paramtypes.BaseCurrency, math.NewInt(0)),
			Liabilities:       math.NewInt(0),
			AmmPoolId:         info.PoolId,
			PositionHealth:    math.LegacyNewDec(0),
		}
		debt := stablestaketypes.Debt{
			Address:               position.GetPositionAddress().String(),
			Borrowed:              info.Borrowed,
			InterestPaid:          math.ZeroInt(),
			InterestStacked:       math.ZeroInt(),
			BorrowTime:            uint64(ctx.BlockTime().Unix()),
			LastInterestCalcTime:  uint64(ctx.BlockTime().Unix()),
			LastInterestCalcBlock: uint64(ctx.BlockHeight()),
		}
		stablestake.SetDebt(ctx, debt)
		leveragelp.SetPosition(ctx, &position)
	}
}

func TestIteratePoolPosIdsStopLossSorted(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	leveragelp := app.LeveragelpKeeper
	simapp.SetStakingParam(app, ctx)
	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

	positions := []struct {
		Id            uint64
		StopLossPrice math.LegacyDec
		PoolId        uint64
	}{
		{
			StopLossPrice: math.LegacyNewDec(100),
			Id:            7,
			PoolId:        2,
		},
		{
			StopLossPrice: math.LegacyNewDec(101),
			Id:            6,
			PoolId:        1,
		},
		{
			StopLossPrice: math.LegacyNewDec(102),
			Id:            5,
			PoolId:        1,
		},
		{
			StopLossPrice: math.LegacyNewDec(103),
			Id:            4,
			PoolId:        1,
		},
		{
			StopLossPrice: math.LegacyNewDec(104),
			Id:            2,
			PoolId:        1,
		},
		{
			StopLossPrice: math.LegacyNewDec(105),
			Id:            1,
			PoolId:        1,
		},
		{
			StopLossPrice: math.LegacyNewDec(106),
			Id:            3,
			PoolId:        1,
		},
	}
	for _, info := range positions {
		position := types.Position{
			LeveragedLpAmount: math.NewInt(1),
			Id:                info.Id,
			Address:           addr[0].String(),
			Collateral:        sdk.NewCoin(paramtypes.BaseCurrency, math.NewInt(0)),
			Liabilities:       math.NewInt(0),
			AmmPoolId:         info.PoolId,
			PositionHealth:    math.LegacyNewDec(0),
			StopLossPrice:     info.StopLossPrice,
		}
		leveragelp.SetPosition(ctx, &position)
	}
}
