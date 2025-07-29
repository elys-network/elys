package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v7/x/amm/keeper"
	"github.com/elys-network/elys/v7/x/amm/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestMsgServerCreatePool() {
	for _, tc := range []struct {
		desc                             string
		senderInitBalance                sdk.Coins
		enableBaseCurrencyPairedPoolOnly bool
		poolParams                       types.PoolParams
		poolAssets                       []types.PoolAsset
		expSenderBalance                 sdk.Coins
		expTotalLiquidity                sdk.Coins
		expLpCommitment                  sdk.Coin
		expPass                          bool
	}{
		{
			desc:                             "sender is not authorized",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10_000_000)},
			enableBaseCurrencyPairedPoolOnly: false,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.Coin{},
			expPass:          false,
		},
		{
			desc:                             "no asset is base asset",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10_000_000)},
			enableBaseCurrencyPairedPoolOnly: false,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.Elys, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.Coin{},
			expPass:          false,
		},
		{
			desc:                             "fee is not base asset",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10_000_000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			enableBaseCurrencyPairedPoolOnly: false,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.Elys,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.Coin{},
			expPass:          false,
		},
		{
			desc:                             "zero tvl pool creation",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10_000_000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			enableBaseCurrencyPairedPoolOnly: false,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.NewCoin("amm/pool/1", sdkmath.NewInt(2000000000000000000)),
			expPass:          true,
		},
		{
			desc:                             "positive tvl pool creation",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			enableBaseCurrencyPairedPoolOnly: false,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.NewCoin("amm/pool/1", sdkmath.NewInt(2).Mul(types.OneShare)),
			expPass:          true,
		},
		{
			desc:                             "not enough balance to create pool",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10000000)},
			enableBaseCurrencyPairedPoolOnly: false,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.Coin{},
			expPass:          false,
		},
		{
			desc:                             "base currency paired pool creation without base currency",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			enableBaseCurrencyPairedPoolOnly: true,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.Elys, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.NewCoin("amm/pool/1", sdkmath.NewInt(2).Mul(types.OneShare)),
			expPass:          false,
		},
		{
			desc:                             "base currency paired pool creation with base currency",
			senderInitBalance:                sdk.Coins{sdk.NewInt64Coin(ptypes.Eden, 1000000), sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			enableBaseCurrencyPairedPoolOnly: true,
			poolParams: types.PoolParams{
				SwapFee:   sdkmath.LegacyZeroDec(),
				UseOracle: false,
				FeeDenom:  ptypes.BaseCurrency,
			},
			poolAssets: []types.PoolAsset{
				{
					Token:  sdk.NewInt64Coin(ptypes.Eden, 1000000),
					Weight: sdkmath.OneInt(),
				},
				{
					Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000),
					Weight: sdkmath.OneInt(),
				},
			},
			expSenderBalance: sdk.Coins{},
			expLpCommitment:  sdk.NewCoin("amm/pool/1", sdkmath.NewInt(2).Mul(types.OneShare)),
			expPass:          true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			suite.SetupStableCoinPrices()

			a := suite.app.OracleKeeper.GetAllAssetInfo(suite.ctx)
			fmt.Println(a)
			b := suite.app.AssetprofileKeeper.GetAllEntry(suite.ctx)
			fmt.Println(b)

			// bootstrap accounts
			sender := authtypes.NewModuleAddress(govtypes.ModuleName)
			if tc.desc == "sender is not authorized" {
				sender = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			}

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)

			// execute function
			msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)

			// set params
			params := suite.app.AmmKeeper.GetParams(suite.ctx)
			params.BaseAssets = []string{ptypes.BaseCurrency}
			suite.app.AmmKeeper.SetParams(suite.ctx, params)

			resp, err := msgServer.CreatePool(
				suite.ctx,
				&types.MsgCreatePool{
					Sender:     sender.String(),
					PoolParams: tc.poolParams,
					PoolAssets: tc.poolAssets,
				})
			if !tc.expPass {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(resp.PoolID, uint64(1))
				suite.Require().True(suite.VerifyPoolAssetWithBalance(1))

				pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)
				suite.Require().Len(pools, 1)
				suite.Require().Equal(pools[0].PoolId, uint64(1))
				suite.Require().Equal(pools[0].PoolParams, tc.poolParams)
				suite.Require().Equal(tc.expLpCommitment.Amount.String(), pools[0].TotalShares.Amount.String())

				totalWeight := sdkmath.ZeroInt()
				for _, poolAsset := range tc.poolAssets {
					totalWeight = totalWeight.Add(poolAsset.Weight)
				}
				suite.Require().Equal(pools[0].TotalWeight.String(), totalWeight.MulRaw(types.GuaranteedWeightPrecision).String())

				// check balance change on sender
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
				suite.Require().Equal(balances.String(), tc.expSenderBalance.String())

				// check lp token commitment
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, sender)
				suite.Require().Len(commitments.CommittedTokens, 1)
				suite.Require().Equal(commitments.CommittedTokens[0].Denom, tc.expLpCommitment.Denom)
				suite.Require().Equal(commitments.CommittedTokens[0].Amount.String(), tc.expLpCommitment.Amount.String())
			}
		})
	}
}
