package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	testkeeper "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestShowFeeInfoQuery(t *testing.T) {
	keeper, ctx := testkeeper.MasterchefKeeper(t)

	newInfo := types.FeeInfo{
		GasLp:        math.NewInt(300),
		GasStakers:   math.NewInt(150),
		GasProtocol:  math.NewInt(75),
		DexLp:        math.NewInt(400),
		DexStakers:   math.NewInt(200),
		DexProtocol:  math.NewInt(100),
		PerpLp:       math.NewInt(500),
		PerpStakers:  math.NewInt(250),
		PerpProtocol: math.NewInt(125),
		EdenLp:       math.NewInt(2000),
	}
	keeper.SetFeeInfo(ctx, newInfo, "2024-05-01")

	response, err := keeper.ShowFeeInfo(ctx, &types.QueryShowFeeInfoRequest{Date: "2024-05-01"})
	require.NoError(t, err)
	require.Equal(t, &types.QueryShowFeeInfoResponse{FeeInfo: newInfo}, response)
}

func TestListFeeInfoQuery(t *testing.T) {
	keeper, ctx := testkeeper.MasterchefKeeper(t)

	newInfo := types.FeeInfo{
		GasLp:        math.NewInt(300),
		GasStakers:   math.NewInt(150),
		GasProtocol:  math.NewInt(75),
		DexLp:        math.NewInt(400),
		DexStakers:   math.NewInt(200),
		DexProtocol:  math.NewInt(100),
		PerpLp:       math.NewInt(500),
		PerpStakers:  math.NewInt(250),
		PerpProtocol: math.NewInt(125),
		EdenLp:       math.NewInt(2000),
	}
	keeper.SetFeeInfo(ctx, newInfo, "2024-05-01")
	keeper.SetFeeInfo(ctx, newInfo, "2024-05-02")

	response, err := keeper.ListFeeInfo(ctx, &types.QueryListFeeInfoRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryListFeeInfoResponse{FeeInfo: []types.FeeInfo{newInfo, newInfo}}, response)
}
