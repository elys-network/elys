package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestExitPoolEstimation(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})
	k := app.AmmKeeper

	// Setup mock pools and assets
	SetupMockPools(&k, ctx)

	// Try single coin exit pool on non-oracle pool (it's invalid)
	resp, err := k.ExitPoolEstimation(ctx, &types.QueryExitPoolEstimationRequest{
		PoolId:        1,
		ShareAmountIn: types.OneShare.Quo(sdk.NewInt(10)),
		TokenOutDenom: "denom2",
	})
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(resp.AmountsOut).String(), "100denom1,100denom2")

	// Test multiple coins exit pool
	resp, err = k.ExitPoolEstimation(ctx, &types.QueryExitPoolEstimationRequest{
		PoolId:        1,
		ShareAmountIn: types.OneShare.Quo(sdk.NewInt(10)),
		TokenOutDenom: "denom2",
	})
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(resp.AmountsOut).String(), "100denom1,100denom2")
}
