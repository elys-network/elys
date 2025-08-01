package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	assetprofiletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/stablestake/keeper"
	"github.com/elys-network/elys/v7/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestUnbond() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		moduleInitBalance sdk.Coins
		unbondAmount      math.Int
		maxWithdrawRatio  math.LegacyDec
		expSenderBalance  sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful unbonding process",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenomForPool(1), 5000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			unbondAmount:      math.NewInt(1000000),
			maxWithdrawRatio:  math.LegacyOneDec(),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}.Sort(),
			expPass:           true,
		},
		{
			desc:              "lack of balance on the module",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenomForPool(1), 5000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000)},
			unbondAmount:      math.NewInt(1000000),
			maxWithdrawRatio:  math.LegacyOneDec(),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:           false,
		},
		{
			desc:              "lack of sender balance",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenomForPool(1), 5000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			unbondAmount:      math.NewInt(10000000000000),
			maxWithdrawRatio:  math.LegacyOneDec(),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:           false,
		},
		{
			desc:              "max withdrawal amount",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenomForPool(1), 5000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			unbondAmount:      math.NewInt(700000), // try to withdraw more than 90%
			maxWithdrawRatio:  math.LegacyMustNewDecFromStr("0.9"),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 500000)},
			expPass:           false,
		},
		{
			desc:              "max withdrawal amount",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenomForPool(1), 5000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			unbondAmount:      math.NewInt(700000), // try to withdraw more than 90%
			maxWithdrawRatio:  math.LegacyMustNewDecFromStr("0.9"),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 500000)},
			expPass:           false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.moduleInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.ModuleName, tc.moduleInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)

			shareDenom := types.GetShareDenomForPool(1)

			// Set an entity to assetprofile
			entry := assetprofiletypes.Entry{
				Authority:       authtypes.NewModuleAddress(types.ModuleName).String(),
				BaseDenom:       shareDenom,
				Decimals:        ptypes.BASE_DECIMAL,
				Denom:           shareDenom,
				DisplayName:     shareDenom,
				CommitEnabled:   true,
				WithdrawEnabled: true,
			}
			suite.app.AssetprofileKeeper.SetEntry(suite.ctx, entry)

			// Commit LP token
			err = suite.app.CommitmentKeeper.CommitLiquidTokens(
				suite.ctx, sender,
				tc.senderInitBalance[0].Denom, tc.senderInitBalance[0].Amount,
				uint64(suite.ctx.BlockTime().Unix()),
			)
			suite.Require().NoError(err)

			pool, _ := suite.app.StablestakeKeeper.GetPool(suite.ctx, 1)
			pool.NetAmount = math.NewInt(5000_000)
			pool.MaxLeverageRatio = math.LegacyMustNewDecFromStr("0.8")
			pool.MaxWithdrawRatio = tc.maxWithdrawRatio
			suite.app.StablestakeKeeper.SetPool(suite.ctx, pool)

			msgServer := keeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
			_, err = msgServer.Unbond(
				suite.ctx,
				&types.MsgUnbond{
					Creator: sender.String(),
					Amount:  tc.unbondAmount,
					PoolId:  1,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())
			}
		})
	}
}
