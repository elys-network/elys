package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestMsgServerUnbond() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		moduleInitBalance sdk.Coins
		unbondAmount      sdk.Int
		expSenderBalance  sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful unbonding process",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenom(), 1000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			unbondAmount:      sdk.NewInt(1000000),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}.Sort(),
			expPass:           true,
		},
		{
			desc:              "lack of balance on the module",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenom(), 1000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000)},
			unbondAmount:      sdk.NewInt(1000000),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			expPass:           false,
		},
		{
			desc:              "lack of sender balance",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenom(), 1000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			unbondAmount:      sdk.NewInt(10000000000000),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
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

			shareDenom := types.GetShareDenom()

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

			// Create a commit LP token message
			msgLiquidCommitLPToken := &ctypes.MsgCommitLiquidTokens{
				Creator:   sender.String(),
				Denom:     tc.senderInitBalance[0].Denom,
				Amount:    tc.senderInitBalance[0].Amount,
				LockUntil: uint64(suite.ctx.BlockTime().Unix()),
			}

			// Commit LP token
			cMsgServer := commitmentkeeper.NewMsgServerImpl(suite.app.CommitmentKeeper)
			_, err = cMsgServer.CommitLiquidTokens(sdk.WrapSDKContext(suite.ctx), msgLiquidCommitLPToken)
			suite.Require().NoError(err)

			params := suite.app.StablestakeKeeper.GetParams(suite.ctx)
			params.TotalValue = sdk.NewInt(1000_000_000)
			suite.app.StablestakeKeeper.SetParams(suite.ctx, params)

			msgServer := keeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
			_, err = msgServer.Unbond(
				sdk.WrapSDKContext(suite.ctx),
				&types.MsgUnbond{
					Creator: sender.String(),
					Amount:  tc.unbondAmount,
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
