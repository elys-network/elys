package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestMsgServerBond() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		bondAmount        math.Int
		expSenderBalance  sdk.Coins
		expSenderCommit   sdk.Coin
		expPass           bool
	}{
		{
			desc:              "successful bonding process",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			bondAmount:        math.NewInt(10000),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 990000)}.Sort(),
			expSenderCommit:   sdk.NewInt64Coin(types.GetShareDenomForPool(1), 10000),
			expPass:           true,
		},
		{
			desc:              "lack of balance",
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
				suite.ctx,
				&types.MsgBond{
					Creator: sender.String(),
					Amount:  tc.bondAmount,
					PoolId:  1,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				// check committed tokens
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, sender)
				suite.Require().Len(commitments.CommittedTokens, 1)
				suite.Require().Equal(commitments.CommittedTokens[0].Amount.String(), tc.expSenderCommit.Amount.String())
				suite.Require().Equal(commitments.CommittedTokens[0].Denom, tc.expSenderCommit.Denom)

			}
		})
	}
}
