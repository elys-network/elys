package types_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPoolAddress(t *testing.T) {
	require.Equal(t, authtypes.NewModuleAddress(types.ModuleName), types.PoolAddress())
}
func TestShareDenom(t *testing.T) {
	require.Equal(t, "stablestake/share", types.GetShareDenom())
}

func TestAddLiabilities(t *testing.T) {
	p := types.AmmPool{
		Id:               1,
		TotalLiabilities: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)},
	}

	p.AddLiabilities(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, p.TotalLiabilities.AmountOf(sdk.DefaultBondDenom), math.NewInt(1010))

	p.AddLiabilities(sdk.NewInt64Coin(ptypes.ATOM, 10))
	require.Equal(t, p.TotalLiabilities.AmountOf(ptypes.ATOM), math.NewInt(10))
}

func TestSubLiabilities(t *testing.T) {
	p := types.AmmPool{
		Id:               1,
		TotalLiabilities: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)},
	}

	p.SubLiabilities(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, p.TotalLiabilities.AmountOf(sdk.DefaultBondDenom), math.NewInt(990))

}
