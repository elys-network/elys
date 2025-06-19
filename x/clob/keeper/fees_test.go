package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (suite *KeeperTestSuite) TestCollectTradingFees() {
	suite.SetupTest()
	markets := suite.CreateMarketWithFees(BaseDenom)
	market := markets[0]
	subAccounts := suite.SetupSubAccounts(2, sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, 1000_000_000), sdk.NewInt64Coin(QuoteDenom, 1000_000_000)))
	buyerAcc := subAccounts[0]
	sellerAcc := subAccounts[1]
	testCases := []struct {
		name               string
		trade              types.Trade
		expectErr          bool
		expectedBuyerFees  math.Int
		expectedSellerFees math.Int
		expectedTotalFees  math.Int
		setup              func(ctx sdk.Context)
	}{
		{
			name: "Success: Buyer is taker, both pay fees",
			trade: types.Trade{
				BuyerSubAccount:     buyerAcc,
				SellerSubAccount:    sellerAcc,
				Price:               math.LegacyMustNewDecFromStr("100"), // 100 uusdc
				Quantity:            math.LegacyMustNewDecFromStr("10"),  // 10 units
				IsBuyerTaker:        true,
				IsBuyerLiquidation:  false,
				IsSellerLiquidation: false,
			},
			expectErr: false,
			// Trade value = 100 * 10 = 1000
			// Buyer is taker, pays 1000 * 0.002 = 2
			expectedBuyerFees: math.NewInt(2),
			// Seller is maker, pays 1000 * 0.001 = 1
			expectedSellerFees: math.NewInt(1),
			// Total fees collected = 2 + 1 = 3
			expectedTotalFees: math.NewInt(3),
		},
		{
			name: "Success: Seller is taker, both pay fees",
			trade: types.Trade{
				BuyerSubAccount:     buyerAcc,
				SellerSubAccount:    sellerAcc,
				Price:               math.LegacyMustNewDecFromStr("100"),
				Quantity:            math.LegacyMustNewDecFromStr("10"),
				IsBuyerTaker:        false, // Seller is the taker
				IsBuyerLiquidation:  false,
				IsSellerLiquidation: false,
			},
			expectErr: false,
			// Trade value = 100 * 10 = 1000
			// Buyer is maker, pays 1000 * 0.001 = 1
			expectedBuyerFees: math.NewInt(1),
			// Seller is taker, pays 1000 * 0.002 = 2
			expectedSellerFees: math.NewInt(2),
			// Total fees collected = 1 + 2 = 3
			expectedTotalFees: math.NewInt(3),
		},
		{
			name: "Success: Buyer is liquidated, pays no fee",
			trade: types.Trade{
				BuyerSubAccount:     buyerAcc,
				SellerSubAccount:    sellerAcc,
				Price:               math.LegacyMustNewDecFromStr("100"),
				Quantity:            math.LegacyMustNewDecFromStr("10"),
				IsBuyerTaker:        true,
				IsBuyerLiquidation:  true, // Buyer is liquidated
				IsSellerLiquidation: false,
			},
			expectErr: false,
			// Buyer pays 0 fee due to liquidation
			expectedBuyerFees: math.NewInt(0),
			// Seller is maker, pays 1000 * 0.001 = 1
			expectedSellerFees: math.NewInt(1),
			// Total fees collected = 0 + 1 = 1
			expectedTotalFees: math.NewInt(1),
		},
		{
			name: "Success: Seller is liquidated, pays no fee",
			trade: types.Trade{
				BuyerSubAccount:     buyerAcc,
				SellerSubAccount:    sellerAcc,
				Price:               math.LegacyMustNewDecFromStr("100"),
				Quantity:            math.LegacyMustNewDecFromStr("10"),
				IsBuyerTaker:        true,
				IsBuyerLiquidation:  false,
				IsSellerLiquidation: true, // Seller is liquidated
			},
			expectErr: false,
			// Buyer is taker, pays 1000 * 0.002 = 2
			expectedBuyerFees: math.NewInt(2),
			// Seller pays 0 fee due to liquidation
			expectedSellerFees: math.NewInt(0),
			// Total fees collected = 2 + 0 = 2
			expectedTotalFees: math.NewInt(2),
		},
		{
			name: "Failure: Buyer has insufficient funds for fee",
			trade: types.Trade{
				BuyerSubAccount:  buyerAcc,
				SellerSubAccount: sellerAcc,
				Price:            math.LegacyMustNewDecFromStr("100"),
				Quantity:         math.LegacyMustNewDecFromStr("10"),
				IsBuyerTaker:     true,
			},
			expectErr:          true,           // Expect an error from the bank keeper
			expectedBuyerFees:  math.NewInt(0), // No fees should be transferred
			expectedSellerFees: math.NewInt(0),
			expectedTotalFees:  math.NewInt(0),
			setup: func(ctx sdk.Context) {
				buyerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, buyerAcc.GetTradingAccountAddress(), QuoteDenom)
				err := suite.app.ClobKeeper.TransferFromSubAccountToSubAccount(ctx, buyerAcc, sellerAcc, sdk.NewCoins(buyerBalance))
				suite.Require().NoError(err)
			},
		},
		{
			name: "Success: Zero fee rates, no fees collected",
			trade: types.Trade{
				BuyerSubAccount:  buyerAcc,
				SellerSubAccount: sellerAcc,
				Price:            math.LegacyMustNewDecFromStr("100"),
				Quantity:         math.LegacyMustNewDecFromStr("10"),
				IsBuyerTaker:     true,
			},
			expectErr:          false,
			expectedBuyerFees:  math.NewInt(0),
			expectedSellerFees: math.NewInt(0),
			expectedTotalFees:  math.NewInt(0),
			setup: func(ctx sdk.Context) {
				market.MakerFeeRate = math.LegacyZeroDec()
				market.TakerFeeRate = math.LegacyZeroDec()
				suite.app.ClobKeeper.SetPerpetualMarket(ctx, market)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			if tc.setup != nil {
				tc.setup(suite.ctx)
			}

			initialBuyerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, buyerAcc.GetTradingAccountAddress(), QuoteDenom)
			initialSellerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sellerAcc.GetTradingAccountAddress(), QuoteDenom)
			initialInsuranceBalance := suite.app.BankKeeper.GetBalance(suite.ctx, market.GetInsuranceAccount(), QuoteDenom)

			err := suite.app.ClobKeeper.CollectTradingFees(suite.ctx, market, tc.trade)
			suite.Require().Equal(err != nil, tc.expectErr)
			// Check final balances
			finalBuyerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, buyerAcc.GetTradingAccountAddress(), QuoteDenom)
			finalSellerBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sellerAcc.GetTradingAccountAddress(), QuoteDenom)
			finalInsuranceBalance := suite.app.BankKeeper.GetBalance(suite.ctx, market.GetInsuranceAccount(), QuoteDenom)

			// Verify buyer's balance changed correctly
			expectedBuyerBalance := initialBuyerBalance.Amount.Sub(tc.expectedBuyerFees)
			suite.Require().True(expectedBuyerBalance.Equal(finalBuyerBalance.Amount), "unexpected buyer balance: got %s, want %s", finalBuyerBalance.Amount, expectedBuyerBalance)

			// Verify seller's balance changed correctly
			expectedSellerBalance := initialSellerBalance.Amount.Sub(tc.expectedSellerFees)
			suite.Require().True(expectedSellerBalance.Equal(finalSellerBalance.Amount), "unexpected seller balance: got %s, want %s", finalSellerBalance.Amount, expectedSellerBalance)

			// Verify insurance fund's balance changed correctly
			expectedInsuranceBalance := initialInsuranceBalance.Amount.Add(tc.expectedTotalFees)
			suite.Require().True(expectedInsuranceBalance.Equal(finalInsuranceBalance.Amount), "unexpected insurance balance: got %s, want %s", finalInsuranceBalance.Amount, expectedInsuranceBalance)
		})
	}
}
