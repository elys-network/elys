package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/stretchr/testify/require"
)

// TODO: v0.50Upgrade - test with detail
func TestCalcDelegationAmount(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount(t)
	ctx := app.BaseApp.NewContext(true)

	simapp.SetStakingParam(app, ctx)

	ek := app.EstakingKeeper
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000))

	// Check with non-delegator
	delegatedAmount := ek.CalcDelegationAmount(ctx, addr[0])
	require.Equal(t, delegatedAmount, math.ZeroInt())

	// Check with genesis account (delegator)
	delegatedAmount = ek.CalcDelegationAmount(ctx, genAccount)
	require.Equal(t, delegatedAmount, sdk.DefaultPowerReduction)
}
