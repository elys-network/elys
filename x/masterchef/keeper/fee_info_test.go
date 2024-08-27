package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestFeeInfo(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	// Test AddFeeInfo for gas fees
	app.MasterchefKeeper.AddFeeInfo(ctx, sdk.NewDec(100), sdk.NewDec(50), sdk.NewDec(25), true)
	dateString := ctx.BlockTime().Format("2006-01-02")
	info := app.MasterchefKeeper.GetFeeInfo(ctx, dateString)
	require.Equal(t, sdk.NewInt(100), info.GasLp)
	require.Equal(t, sdk.NewInt(50), info.GasStakers)
	require.Equal(t, sdk.NewInt(25), info.GasProtocol)

	// Test AddFeeInfo for dex fees
	app.MasterchefKeeper.AddFeeInfo(ctx, sdk.NewDec(200), sdk.NewDec(100), sdk.NewDec(50), false)
	info = app.MasterchefKeeper.GetFeeInfo(ctx, dateString)
	require.Equal(t, sdk.NewInt(200), info.DexLp)
	require.Equal(t, sdk.NewInt(100), info.DexStakers)
	require.Equal(t, sdk.NewInt(50), info.DexProtocol)

	// Test AddEdenInfo
	app.MasterchefKeeper.AddEdenInfo(ctx, sdk.NewDec(1000))
	info = app.MasterchefKeeper.GetFeeInfo(ctx, dateString)
	require.Equal(t, sdk.NewInt(1000), info.EdenLp)

	// Test SetFeeInfo and GetFeeInfo
	newInfo := types.FeeInfo{
		GasLp:        sdk.NewInt(300),
		GasStakers:   sdk.NewInt(150),
		GasProtocol:  sdk.NewInt(75),
		DexLp:        sdk.NewInt(400),
		DexStakers:   sdk.NewInt(200),
		DexProtocol:  sdk.NewInt(100),
		PerpLp:       sdk.NewInt(500),
		PerpStakers:  sdk.NewInt(250),
		PerpProtocol: sdk.NewInt(125),
		EdenLp:       sdk.NewInt(2000),
	}
	newDateString := ctx.BlockTime().AddDate(0, 0, 1).Format("2006-01-02")
	app.MasterchefKeeper.SetFeeInfo(ctx, newInfo, newDateString)
	retrievedInfo := app.MasterchefKeeper.GetFeeInfo(ctx, newDateString)
	require.Equal(t, newInfo, retrievedInfo)

	zeroInfo := types.FeeInfo{
		GasLp:        sdk.NewInt(0),
		GasStakers:   sdk.NewInt(0),
		GasProtocol:  sdk.NewInt(0),
		DexLp:        sdk.NewInt(0),
		DexStakers:   sdk.NewInt(0),
		DexProtocol:  sdk.NewInt(0),
		PerpLp:       sdk.NewInt(0),
		PerpStakers:  sdk.NewInt(0),
		PerpProtocol: sdk.NewInt(0),
		EdenLp:       sdk.NewInt(0),
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
