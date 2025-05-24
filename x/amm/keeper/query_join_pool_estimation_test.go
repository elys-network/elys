package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v5/app"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestJoinPoolEstimation(t *testing.T) {
	app := simapp.InitElysTestApp(initChain, t)
	ctx := app.BaseApp.NewContext(initChain)
	k := app.AmmKeeper

	// Setup mock pools and assets
	SetupMockPools(k, ctx)

	// Test single coin join pool
	resp, err := k.JoinPoolEstimation(ctx, &types.QueryJoinPoolEstimationRequest{
		PoolId:    1,
		AmountsIn: sdk.NewCoins(sdk.NewInt64Coin("denom1", 200)),
	})
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(resp.AmountsIn).String(), "200denom1")
	require.Equal(t, resp.ShareAmountOut.String(), "95445115010332227amm/pool/1")

	// Test multiple coins join pool
	resp, err = k.JoinPoolEstimation(ctx, &types.QueryJoinPoolEstimationRequest{
		PoolId:    1,
		AmountsIn: sdk.NewCoins(sdk.NewInt64Coin("denom1", 100), sdk.NewInt64Coin("denom2", 200)),
	})
	require.NoError(t, err)
	// Pool ratio is 1:1, so join pool will accept only 100denom1,100denom2
	require.Equal(t, sdk.Coins(resp.AmountsIn).String(), "100denom1,100denom2")
	require.Equal(t, resp.ShareAmountOut.String(), "100000000000000000amm/pool/1")
}
