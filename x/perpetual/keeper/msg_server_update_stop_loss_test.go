package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"

	ptypes "github.com/elys-network/elys/x/parameter/types"
) 

func TestUpdate_Stop_Loss(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	perpetual := app.PerpetualKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	for i := 0; i < 2; i++ {
		// update_stop_loss := types.MsgUpdateStopLoss{
		// 	Creator:  addr[i].String(),
		// 	Id: 	  0,
		// 	Price:    sdk.MustNewDecFromStr("100"),
		// }
		mtp := types.MTP{
			Address:                        addr[i].String(),
			CollateralAsset:                ptypes.BaseCurrency,
			CustodyAsset:                   "ATOM",
			Collateral:                     sdk.NewInt(0),
			Liabilities:                    sdk.NewInt(0),
			BorrowInterestPaidCollateral:   sdk.NewInt(0),
			BorrowInterestPaidCustody:      sdk.NewInt(0),
			BorrowInterestUnpaidCollateral: sdk.NewInt(0),
			Custody:                        sdk.NewInt(0),
			Leverage:                       sdk.NewDec(0),
			MtpHealth:                      sdk.NewDec(0),
			Position:                       types.Position_LONG,
			Id:                             0,
			ConsolidateLeverage:            sdk.ZeroDec(),
			SumCollateral:                  sdk.ZeroInt(),
			StopLossPrice:                  sdk.ZeroDec(),
		}
		err := perpetual.SetMTP(ctx, &mtp)
		require.NoError(t, err)
	}

	// perpetual
}