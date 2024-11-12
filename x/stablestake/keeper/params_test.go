package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	testkeeper "github.com/elys-network/elys/testutil/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.StablestakeKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

func (suite *KeeperTestSuite) TestGetRedemptionRate() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		bondAmount        math.Int
		expSenderBalance  sdk.Coins
		expSenderCommit   sdk.Coin
		expPass           bool
	}{
		{
			desc:              "successful bonding process, redemption should be set",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			bondAmount:        math.NewInt(10000),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 990000)}.Sort(),
			expSenderCommit:   sdk.NewInt64Coin(types.GetShareDenom(), 10000),
			expPass:           true,
		},
		{
			desc:              "lack of balance, redemption should not be set",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			bondAmount:        math.NewInt(10000000000000),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expSenderCommit:   sdk.Coin{},
			expPass:           false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)

			msgServer := keeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
			_, err = msgServer.Bond(
				sdk.WrapSDKContext(suite.ctx),
				&types.MsgBond{
					Creator: sender.String(),
					Amount:  tc.bondAmount,
				})
			if !tc.expPass {
				suite.Require().Error(err)

				// Check redemption rate
				rate := suite.app.StablestakeKeeper.GetRedemptionRate(suite.ctx)
				suite.Require().Equal(math.LegacyZeroDec(), rate)
			} else {
				suite.Require().NoError(err)

				// Check redemption rate
				rate := suite.app.StablestakeKeeper.GetRedemptionRate(suite.ctx)
				suite.Require().Equal(math.LegacyNewDec(1), rate)
			}
		})
	}
}
