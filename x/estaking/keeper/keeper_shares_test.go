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
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	ek := app.EstakingKeeper
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000))

	// Check with non-delegator
	delegatedAmount := ek.CalcDelegationAmount(ctx, addr[0])
	require.Equal(t, delegatedAmount, sdk.ZeroInt())

	// Check with genesis account (delegator)
	delegatedAmount = ek.CalcDelegationAmount(ctx, genAccount)
	require.Equal(t, delegatedAmount, sdk.DefaultPowerReduction)
}
