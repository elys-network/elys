package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/leveragelp/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestSetGetMTP(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	leveragelp := app.LeveragelpKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	for i := 0; i < 2; i++ {
		mtp := types.MTP{
			Address:      addr[i].String(),
			Collateral:   sdk.NewCoin(paramtypes.BaseCurrency, sdk.NewInt(0)),
			Liabilities:  sdk.NewInt(0),
			InterestPaid: sdk.NewInt(0),
			AmmPoolId:    1,
			Leverage:     sdk.NewDec(0),
			MtpHealth:    sdk.NewDec(0),
			Id:           0,
		}
		leveragelp.SetMTP(ctx, &mtp)
	}

	mtpCount := leveragelp.GetMTPCount(ctx)
	require.Equal(t, mtpCount, (uint64)(2))
}
