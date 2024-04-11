package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/stretchr/testify/require"
)

func TestCalcDelegationAmount(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000))

	// Check with non-delegator
	delegatedAmount := ik.CalcDelegationAmount(ctx, addr[0].String())
	require.Equal(t, delegatedAmount, sdk.ZeroInt())

	// Check with genesis account (delegator)
	delegatedAmount = ik.CalcDelegationAmount(ctx, genAccount.String())
	require.Equal(t, delegatedAmount, sdk.DefaultPowerReduction)
}
