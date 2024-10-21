package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestFeeInfo(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	// Test AddFeeInfo for gas fees
	app.MasterchefKeeper.AddFeeInfo(ctx, sdkmath.LegacyNewDec(100), sdkmath.LegacyNewDec(50), sdkmath.LegacyNewDec(25), true)
	dateString := ctx.BlockTime().Format("2006-01-02")
	info := app.MasterchefKeeper.GetFeeInfo(ctx, dateString)
	require.Equal(t, sdkmath.NewInt(100), info.GasLp)
	require.Equal(t, sdkmath.NewInt(50), info.GasStakers)
	require.Equal(t, sdkmath.NewInt(25), info.GasProtocol)

	// Test AddFeeInfo for dex fees
	app.MasterchefKeeper.AddFeeInfo(ctx, sdkmath.LegacyNewDec(200), sdkmath.LegacyNewDec(100), sdkmath.LegacyNewDec(50), false)
	info = app.MasterchefKeeper.GetFeeInfo(ctx, dateString)
	require.Equal(t, sdkmath.NewInt(200), info.DexLp)
	require.Equal(t, sdkmath.NewInt(100), info.DexStakers)
	require.Equal(t, sdkmath.NewInt(50), info.DexProtocol)

	// Test AddEdenInfo
	app.MasterchefKeeper.AddEdenInfo(ctx, sdkmath.LegacyNewDec(1000))
	info = app.MasterchefKeeper.GetFeeInfo(ctx, dateString)
	require.Equal(t, sdkmath.NewInt(1000), info.EdenLp)

	// Test SetFeeInfo and GetFeeInfo
	newInfo := types.FeeInfo{
		GasLp:        sdkmath.NewInt(300),
		GasStakers:   sdkmath.NewInt(150),
		GasProtocol:  sdkmath.NewInt(75),
		DexLp:        sdkmath.NewInt(400),
		DexStakers:   sdkmath.NewInt(200),
		DexProtocol:  sdkmath.NewInt(100),
		PerpLp:       sdkmath.NewInt(500),
		PerpStakers:  sdkmath.NewInt(250),
		PerpProtocol: sdkmath.NewInt(125),
		EdenLp:       sdkmath.NewInt(2000),
	}
	newDateString := ctx.BlockTime().AddDate(0, 0, 1).Format("2006-01-02")
	app.MasterchefKeeper.SetFeeInfo(ctx, newInfo, newDateString)
	retrievedInfo := app.MasterchefKeeper.GetFeeInfo(ctx, newDateString)
	require.Equal(t, newInfo, retrievedInfo)

	zeroInfo := types.FeeInfo{
		GasLp:        sdkmath.NewInt(0),
		GasStakers:   sdkmath.NewInt(0),
		GasProtocol:  sdkmath.NewInt(0),
		DexLp:        sdkmath.NewInt(0),
		DexStakers:   sdkmath.NewInt(0),
		DexProtocol:  sdkmath.NewInt(0),
		PerpLp:       sdkmath.NewInt(0),
		PerpStakers:  sdkmath.NewInt(0),
		PerpProtocol: sdkmath.NewInt(0),
		EdenLp:       sdkmath.NewInt(0),
	}
	// Test GetFeeInfo for non-existent date
	nonExistentDate := ctx.BlockTime().AddDate(0, 0, 2).Format("2006-01-02")
	emptyInfo := app.MasterchefKeeper.GetFeeInfo(ctx, nonExistentDate)
	require.Equal(t, zeroInfo, emptyInfo)

	//Test RemoveFeeInfo
	app.MasterchefKeeper.RemoveFeeInfo(ctx, dateString)
	removedInfo := app.MasterchefKeeper.GetFeeInfo(ctx, dateString)
	require.Equal(t, zeroInfo, removedInfo)

	// Test GetAllFeeInfos
	allInfos := app.MasterchefKeeper.GetAllFeeInfos(ctx)
	require.Len(t, allInfos, 1)
	require.Equal(t, newInfo, allInfos[0])
}
