package keeper_test

import (
	"cosmossdk.io/math"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"

	"github.com/elys-network/elys/x/accountedpool/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func TestAccountedPoolUpdate(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	apk := app.AccountedPoolKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Initiate pool
	ammPool := ammtypes.Pool{
		PoolId:      0,
		Address:     addr[0].String(),
		PoolParams:  ammtypes.PoolParams{},
		TotalShares: sdk.NewCoin("lp-token", sdk.NewInt(100)),
		PoolAssets: []ammtypes.PoolAsset{
			{Token: sdk.NewCoin(ptypes.ATOM, sdk.NewInt(5000))},
			{Token: sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000))},
		},
		TotalWeight:       sdk.NewInt(100),
		RebalanceTreasury: addr[0].String(),
	}
	// Initiate pool
	accountedPool := types.AccountedPool{
		PoolId:           0,
		TotalShares:      ammPool.TotalShares,
		PoolAssets:       []ammtypes.PoolAsset{},
		TotalWeight:      ammPool.TotalWeight,
		NonAmmPoolTokens: []sdk.Coin{},
	}

	for _, asset := range ammPool.PoolAssets {
		accountedPool.PoolAssets = append(accountedPool.PoolAssets, asset)
	}
	// Set accounted pool
	apk.SetAccountedPool(ctx, accountedPool)

	perpetualPool := perpetualtypes.Pool{
		AmmPoolId:          0,
		Health:             sdk.NewDec(1),
		BorrowInterestRate: sdk.NewDec(1),
		PoolAssetsLong: []perpetualtypes.PoolAsset{
			{
				Liabilities:           sdk.NewInt(400),
				Custody:               sdk.NewInt(50),
				TakeProfitCustody:     math.NewInt(10),
				TakeProfitLiabilities: math.NewInt(20),
				AssetDenom:            ptypes.BaseCurrency,
			},
			{
				Liabilities:           sdk.NewInt(0),
				Custody:               sdk.NewInt(50),
				TakeProfitCustody:     math.ZeroInt(),
				TakeProfitLiabilities: math.ZeroInt(),
				AssetDenom:            ptypes.ATOM,
			},
		},
		PoolAssetsShort: []perpetualtypes.PoolAsset{
			{
				Liabilities:           sdk.NewInt(400),
				Custody:               sdk.NewInt(70),
				TakeProfitCustody:     math.ZeroInt(),
				TakeProfitLiabilities: math.ZeroInt(),
				AssetDenom:            ptypes.BaseCurrency,
			},
			{
				Liabilities:           sdk.NewInt(0),
				Custody:               sdk.NewInt(50),
				TakeProfitCustody:     math.ZeroInt(),
				TakeProfitLiabilities: math.ZeroInt(),
				AssetDenom:            ptypes.ATOM,
			},
		},
	}
	// Update accounted pool
	err := apk.PerpetualUpdates(ctx, ammPool, perpetualPool)
	require.NoError(t, err)

	apool, found := apk.GetAccountedPool(ctx, (uint64)(0))
	require.Equal(t, found, true)
	require.Equal(t, apool.PoolId, (uint64)(0))

	usdcBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.BaseCurrency)
	require.Equal(t, usdcBalance, sdk.NewInt(1000+400-50+400-70+10-20))
	atomBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.ATOM)
	require.Equal(t, atomBalance, sdk.NewInt(5000-50-50))
}
