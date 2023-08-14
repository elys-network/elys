package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/margin/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestSetGetMTP(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	margin := app.MarginKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	for i := 0; i < 2; i++ {
		mtp := types.MTP{
			Address:                  addr[i].String(),
			CollateralAsset:          paramtypes.USDC,
			CollateralAmount:         sdk.NewInt(0),
			Liabilities:              sdk.NewInt(0),
			InterestPaidCollateral:   sdk.NewInt(0),
			InterestPaidCustody:      sdk.NewInt(0),
			InterestUnpaidCollateral: sdk.NewInt(0),
			CustodyAsset:             "ATOM",
			CustodyAmount:            sdk.NewInt(0),
			Leverage:                 sdk.NewDec(0),
			MtpHealth:                sdk.NewDec(0),
			Position:                 types.Position_LONG,
			Id:                       0,
		}
		nullify.Fill(&mtp)

		err := margin.SetMTP(ctx, &mtp)
		require.NoError(t, err)
	}

	mtpCount := margin.GetMTPCount(ctx)
	require.Equal(t, mtpCount, (uint64)(2))

}

func TestGetAllWhitelistedAddress(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	margin := app.MarginKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	// Set whitelisted addresses
	margin.WhitelistAddress(ctx, addr[0].String())
	margin.WhitelistAddress(ctx, addr[1].String())

	// Get all whitelisted addresses
	whitelisted := margin.GetAllWhitelistedAddress(ctx)

	// length should be 2
	require.Equal(t, len(whitelisted), 2)
	// first result should be matches to the first addr
	require.Equal(t, whitelisted[0], addr[0].String())
	// second result should be matched to the second addr
	require.Equal(t, whitelisted[1], addr[1].String())
}
