package e2e

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"

	"github.com/stretchr/testify/require"
)

var COLLECTOR = "elys1v30pe777dj9mgsnlv0j2c5wh05m0ya0nlhuv7t"

// TestBeginBlocker tests the module's begin blocker logic.
func TestBeginBlocker(t *testing.T) {
	t.Parallel()

	var wrapper Wrapper
	ctx, _, _ := Suite(t, &wrapper, false)
	validator := wrapper.chain.Validators[0]

	oldBalance, err := wrapper.chain.BankQueryAllBalances(ctx, wrapper.owner.FormattedAddress())
	require.NoError(t, err)

	err = validator.BankSend(ctx, wrapper.owner.KeyName(), ibc.WalletAmount{
		Address: wrapper.pendingOwner.FormattedAddress(),
		Denom:   "uelys",
		Amount:  math.NewInt(1_000_000),
	})
	require.NoError(t, err)

	balance, err := wrapper.chain.BankQueryAllBalances(ctx, COLLECTOR)
	require.NoError(t, err)
	require.True(t, balance.IsZero())

	newBalance, err := wrapper.chain.BankQueryAllBalances(ctx, wrapper.owner.FormattedAddress())
	require.NoError(t, err)
	require.Equal(t,
		oldBalance.
			Sub(sdk.NewInt64Coin("uelys", 1_000_000)).
			Sub(sdk.NewInt64Coin("uelys", 20_000)),
		newBalance,
	)
}
