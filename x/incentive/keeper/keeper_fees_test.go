package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestCollectGasFeesToIncentiveModule(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, bk := app.IncentiveKeeper, app.BankKeeper
	// Collect gas fees
	collectedAmt := ik.CollectGasFeesToIncentiveModule(ctx)

	// rewards should be zero
	require.True(t, collectedAmt.IsZero())

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000))
	transferAmt := sdk.NewCoin(ptypes.Elys, sdk.NewInt(100))

	// Deposit 100stake to FeeCollectorName wallet
	err := bk.SendCoinsFromAccountToModule(ctx, addr[0], authtypes.FeeCollectorName, sdk.NewCoins(transferAmt))
	require.NoError(t, err)

	// Collect gas fees again
	collectedAmt = ik.CollectGasFeesToIncentiveModule(ctx)

	// It should be 100stake
	require.Equal(t, collectedAmt, sdk.Coins{sdk.NewCoin(ptypes.Elys, sdk.NewInt(100))})
}
