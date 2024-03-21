package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestJoinPoolEstimation(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})
	k := app.AmmKeeper

	// Setup mock pools and assets
	SetupMockPools(&k, ctx)

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
	require.Equal(t, sdk.Coins(resp.AmountsIn).String(), "100denom1,200denom2")
	require.Equal(t, resp.ShareAmountOut.String(), "100000000000000000amm/pool/1")
}
