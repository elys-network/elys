package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/stretchr/testify/require"
)

func TestPortionCoins(t *testing.T) {
	coins := sdk.Coins{sdk.NewInt64Coin("ueden", 1000), sdk.NewInt64Coin("uelys", 10000)}
	portion := keeper.PortionCoins(coins, sdk.ZeroDec())
	require.Equal(t, portion, sdk.Coins(nil))

	portion = keeper.PortionCoins(coins, sdk.NewDecWithPrec(1, 1))
	require.Equal(t, portion, sdk.Coins{sdk.NewInt64Coin("ueden", 100), sdk.NewInt64Coin("uelys", 1000)})

	portion = keeper.PortionCoins(coins, sdk.NewDec(1))
	require.Equal(t, portion, coins)
}

// TODO: add test OnCollectFee
func TestOnCollectFee(t *testing.T) {

}
