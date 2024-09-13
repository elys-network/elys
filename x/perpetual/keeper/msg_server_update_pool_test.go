package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

type updatePoolTableTest struct {
	title string
	msg   types.MsgUpdatePool
	want  error
}

func TestUpdatePool(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	//legacyAmino := app.LegacyAmino()
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	pool := types.Pool{
		AmmPoolId:                            2,
		Health:                               sdk.MustNewDecFromStr("0.000000536043604121"),
		Enabled:                              true,
		Closed:                               false,
		BorrowInterestRate:                   sdk.MustNewDecFromStr("0.000000536043150009"),
		LastHeightBorrowInterestRateComputed: 9856395,
		FundingRate:                          sdk.MustNewDecFromStr("0.001000000000000000"),
	}
	app.PerpetualKeeper.SetPool(ctx, pool)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	table := []updatePoolTableTest{
		{
			title: "Invalid authority",
			msg: types.MsgUpdatePool{
				Authority: sample.AccAddress(),
			},
			want: govtypes.ErrInvalidSigner,
		},
		{
			title: "Pool not found",
			msg: types.MsgUpdatePool{
				Authority: authority,
				PoolId:    1,
			},
			want: types.ErrPoolDoesNotExist,
		},
		{
			title: "Pool updated",
			msg: types.MsgUpdatePool{
				Authority: authority,
				PoolId:    2,
			},
			want: nil,
		},
	}

	for _, test := range table {
		t.Run(test.title, func(t *testing.T) {
			msgServer := keeper.NewMsgServerImpl(app.PerpetualKeeper)
			_, err := msgServer.UpdatePool(ctx, &test.msg)
			if err != nil {
				assert.ErrorIs(t, err, test.want)
			} else {
				assert.Equal(t, err, test.want)
			}
		})
	}
}
