package keeper_test

//func (suite *KeeperTestSuite) TestExchange() {
//	suite.ResetSuite()
//	baseDenom := "uatom"
//
//	marketId := uint64(1)
//	balancePerAccount := sdk.NewCoins(sdk.NewInt64Coin("uusdc", 10_000_000))
//	subAccounts := suite.SetupSubAccounts(4, balancePerAccount)
//
//	buyerSubAccount_1 := subAccounts[0]
//	sellerSubAccount_1 := subAccounts[1]
//	buyerSubAccount_2 := subAccounts[2]
//
//	t1 := types.Trade{
//		BuyerSubAccount:  buyerSubAccount_1,
//		SellerSubAccount: sellerSubAccount_1,
//		MarketId:         1,
//		Price:            math.LegacyNewDecWithPrec(101, 1),
//		Quantity:         math.LegacyNewDec(100),
//	}
//	t2 := types.Trade{
//		BuyerSubAccount:  buyerSubAccount_2,
//		SellerSubAccount: sellerSubAccount_1,
//		MarketId:         1,
//		Price:            math.LegacyNewDecWithPrec(103, 1),
//		Quantity:         math.LegacyNewDec(150),
//	}
//	testCases := []struct {
//		name           string
//		expectedErrMsg string
//		input          types.Trade
//		pre            func()
//		post           func()
//	}{
//		{
//			"no market",
//			types.ErrPerpetualMarketNotFound.Error(),
//			t1,
//			func() {
//			},
//			func() {
//			},
//		},
//		{
//			"buyer lack funds",
//			"insufficient funds",
//			t1,
//			func() {
//				suite.CreateMarket(baseDenom)
//				err := suite.app.ClobKeeper.SendFromSubAccount(suite.ctx, buyerSubAccount_1, authtypes.NewModuleAddress(types.ModuleName), balancePerAccount)
//				suite.Require().NoError(err)
//			},
//			func() {
//			},
//		},
//		{
//			"seller lack funds",
//			"insufficient funds",
//			t1,
//			func() {
//				err := suite.app.ClobKeeper.AddToSubAccount(suite.ctx, authtypes.NewModuleAddress(types.ModuleName), buyerSubAccount_1, balancePerAccount)
//				suite.Require().NoError(err)
//				err = suite.app.ClobKeeper.SendFromSubAccount(suite.ctx, sellerSubAccount_1, authtypes.NewModuleAddress(types.ModuleName), balancePerAccount)
//				suite.Require().NoError(err)
//			},
//			func() {
//			},
//		},
//		{
//			"success: both subAccounts never opened a position before",
//			"",
//			t1,
//			func() {
//				err := suite.app.ClobKeeper.AddToSubAccount(suite.ctx, authtypes.NewModuleAddress(types.ModuleName), sellerSubAccount_1, balancePerAccount)
//				suite.Require().NoError(err)
//			},
//			func() {
//				lastAvgTradePrice := suite.app.ClobKeeper.GetLastAverageTradePrice(suite.ctx, 1)
//				suite.Require().Equal(math.LegacyNewDecWithPrec(1010, 2), lastAvgTradePrice)
//
//				market, err := suite.app.ClobKeeper.GetPerpetualMarket(suite.ctx, marketId)
//				suite.Require().NoError(err)
//
//				// new quantity gets created
//				suite.Require().Equal(t1.Quantity, market.TotalOpen)
//
//				buyerPerpetualOwner, found := suite.app.ClobKeeper.GetPerpetualOwner(suite.ctx, t1.BuyerSubAccount.GetOwnerAccAddress(), t1.MarketId)
//				suite.Require().True(found)
//				buyerPerpetual, err := suite.app.ClobKeeper.GetPerpetual(suite.ctx, t1.MarketId, buyerPerpetualOwner.PerpetualId)
//				suite.Require().NoError(err)
//				suite.Require().Equal(buyerPerpetual.Quantity, t1.Quantity)
//				suite.Require().Equal(buyerPerpetual.EntryPrice, t1.Price)
//				suite.Require().Equal(buyerPerpetual.EntryFundingRate, math.LegacyZeroDec())
//				//suite.Require().Equal(buyerPerpetual.Margin, t1.GetRequiredInitialMargin(market))
//
//				sellerPerpetualOwner, found := suite.app.ClobKeeper.GetPerpetualOwner(suite.ctx, t1.SellerSubAccount.GetOwnerAccAddress(), t1.MarketId)
//				suite.Require().True(found)
//				sellerPerpetual, err := suite.app.ClobKeeper.GetPerpetual(suite.ctx, t1.MarketId, sellerPerpetualOwner.PerpetualId)
//				suite.Require().NoError(err)
//				suite.Require().Equal(sellerPerpetual.Quantity, t1.Quantity.Neg())
//				suite.Require().Equal(sellerPerpetual.EntryPrice, t1.Price)
//				suite.Require().Equal(sellerPerpetual.EntryFundingRate, math.LegacyZeroDec())
//				//suite.Require().Equal(sellerPerpetual.Margin.String(), t1.GetRequiredInitialMargin(market).String())
//			},
//		},
//		{
//			"success: seller already have short position, buyer doesn't",
//			"",
//			t2,
//			func() {
//				fmt.Println("Case 2")
//			},
//			func() {
//				lastAvgTradePrice := suite.app.ClobKeeper.GetLastAverageTradePrice(suite.ctx, 1)
//				suite.Require().Equal(math.LegacyNewDecWithPrec(1030, 2), lastAvgTradePrice)
//
//				market, err := suite.app.ClobKeeper.GetPerpetualMarket(suite.ctx, marketId)
//				suite.Require().NoError(err)
//
//				// new quantity gets created
//				suite.Require().Equal(t1.Quantity.Add(t2.Quantity), market.TotalOpen)
//
//				buyerPerpetualOwner, found := suite.app.ClobKeeper.GetPerpetualOwner(suite.ctx, t2.BuyerSubAccount.GetOwnerAccAddress(), t2.MarketId)
//				suite.Require().True(found)
//				buyerPerpetual, err := suite.app.ClobKeeper.GetPerpetual(suite.ctx, t2.MarketId, buyerPerpetualOwner.PerpetualId)
//				suite.Require().NoError(err)
//				suite.Require().Equal(buyerPerpetual.Quantity, t2.Quantity)
//				suite.Require().Equal(buyerPerpetual.EntryPrice, t2.Price)
//				suite.Require().Equal(buyerPerpetual.EntryFundingRate, math.LegacyZeroDec())
//				//suite.Require().Equal(buyerPerpetual.Margin.String(), t2.GetRequiredInitialMargin(market).String())
//
//				sellerPerpetualOwner, found := suite.app.ClobKeeper.GetPerpetualOwner(suite.ctx, t2.SellerSubAccount.GetOwnerAccAddress(), t2.MarketId)
//				suite.Require().True(found)
//				sellerPerpetual, err := suite.app.ClobKeeper.GetPerpetual(suite.ctx, t2.MarketId, sellerPerpetualOwner.PerpetualId)
//				suite.Require().NoError(err)
//				suite.Require().Equal(sellerPerpetual.Quantity, t2.Quantity.Add(t1.Quantity).Neg())
//				suite.Require().Equal(sellerPerpetual.EntryPrice, (t2.GetTradeValue().Add(t1.GetTradeValue())).Quo(sellerPerpetual.Quantity.Abs()).Neg())
//				suite.Require().Equal(sellerPerpetual.EntryFundingRate, math.LegacyZeroDec())
//				//fmt.Println(t1.GetRequiredInitialMargin(market))
//				//fmt.Println(t1.GetRequiredInitialMargin(market))
//				//suite.Require().Equal(sellerPerpetual.Margin.String(), t1.GetRequiredInitialMargin(market).Add(t2.GetRequiredInitialMargin(market)).String())
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		suite.Run(tc.name, func() {
//			suite.IncreaseHeight(1)
//			tc.pre()
//			err := suite.app.ClobKeeper.Exchange(suite.ctx, tc.input)
//			if tc.expectedErrMsg != "" {
//				suite.Require().Error(err)
//				suite.Contains(err.Error(), tc.expectedErrMsg)
//			} else {
//				suite.Require().NoError(err)
//			}
//			tc.post()
//		})
//	}
//}
