package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func (suite *AmmKeeperTestSuite) TestUpFrontSwapExactAmountIn() {
	testCases := []struct {
		name          string
		setup         func() error
		msg           *types.MsgUpFrontSwapExactAmountIn
		expectedError error
	}{
		{
			name: "unauthorized sender",
			setup: func() error {
				params := suite.app.AmmKeeper.GetParams(suite.ctx)
				// params.AllowedUpfrontSwapMakers = []string{"cosmos1differentaddress"}
				suite.app.AmmKeeper.SetParams(suite.ctx, params)
				msg := &types.MsgUpFrontSwapExactAmountIn{
					Sender:            "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
					Routes:            []types.SwapAmountInRoute{},
					TokenIn:           sdk.NewCoin("tokenA", sdkmath.NewInt(100)),
					TokenOutMinAmount: sdkmath.Int(sdkmath.LegacyMustNewDecFromStr("50")),
				}
				_, err := suite.app.AmmKeeper.UpFrontSwapExactAmountIn(suite.ctx, msg)
				return err
			},
			expectedError: types.ErrUnauthorizedUpFrontSwap,
		},
		{
			name: "successful swap",
			setup: func() error {
				suite.SetupTest()

				// bootstrap accounts
				sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

				params := suite.app.AmmKeeper.GetParams(suite.ctx)
				params.AllowedUpfrontSwapMakers = []string{sender.String()}
				suite.app.AmmKeeper.SetParams(suite.ctx, params)

				poolAddr := types.NewPoolAddress(uint64(1))
				treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
				poolAddr2 := types.NewPoolAddress(uint64(2))
				treasuryAddr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
				poolCoins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}
				pool2Coins := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin("uusdt", 1000000)}

				// bootstrap balances
				senderInitBalance := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}
				err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, senderInitBalance)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, senderInitBalance)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolCoins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolCoins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, pool2Coins)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr2, pool2Coins)
				suite.Require().NoError(err)

				// execute function
				suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
					Denom:     ptypes.Elys,
					Liquidity: sdkmath.NewInt(2000000),
				})
				suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
					Denom:     ptypes.BaseCurrency,
					Liquidity: sdkmath.NewInt(1000000),
				})
				suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
					Denom:     "uusdt",
					Liquidity: sdkmath.NewInt(1000000),
				})

				pool := types.Pool{
					PoolId:            1,
					Address:           poolAddr.String(),
					RebalanceTreasury: treasuryAddr.String(),
					PoolParams: types.PoolParams{
						SwapFee:  sdkmath.LegacyNewDecWithPrec(1, 2),
						FeeDenom: ptypes.BaseCurrency,
					},
					TotalShares: sdk.Coin{},
					PoolAssets: []types.PoolAsset{
						{
							Token:  poolCoins[0],
							Weight: sdkmath.NewInt(10),
						},
						{
							Token:  poolCoins[1],
							Weight: sdkmath.NewInt(10),
						},
					},
					TotalWeight: sdkmath.ZeroInt(),
				}
				pool2 := types.Pool{
					PoolId:            2,
					Address:           poolAddr2.String(),
					RebalanceTreasury: treasuryAddr2.String(),
					PoolParams: types.PoolParams{
						SwapFee:  sdkmath.LegacyNewDecWithPrec(1, 2),
						FeeDenom: ptypes.BaseCurrency,
					},
					TotalShares: sdk.Coin{},
					PoolAssets: []types.PoolAsset{
						{
							Token:  pool2Coins[0],
							Weight: sdkmath.NewInt(10),
						},
						{
							Token:  pool2Coins[1],
							Weight: sdkmath.NewInt(10),
						},
					},
					TotalWeight: sdkmath.ZeroInt(),
				}
				suite.app.AmmKeeper.SetPool(suite.ctx, pool)
				suite.app.AmmKeeper.SetPool(suite.ctx, pool2)
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(2))
				msg := &types.MsgUpFrontSwapExactAmountIn{
					Sender: sender.String(),
					Routes: []types.SwapAmountInRoute{
						{
							PoolId:        1,
							TokenOutDenom: ptypes.BaseCurrency,
						},
					},
					TokenIn:           sdk.NewInt64Coin(ptypes.Elys, 10000),
					TokenOutMinAmount: sdkmath.ZeroInt(),
				}
				_, err = suite.app.AmmKeeper.UpFrontSwapExactAmountIn(suite.ctx, msg)
				return err
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.setup()
			// Validate the results
			if tc.expectedError != nil {
				require.ErrorIs(suite.T(), err, tc.expectedError)
			} else {
				require.NoError(suite.T(), err)
			}
		})
	}
}
