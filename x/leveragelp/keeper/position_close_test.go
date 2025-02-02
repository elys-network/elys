package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (suite *KeeperTestSuite) OpenPosition(addr sdk.AccAddress) (*types.Position, math.LegacyDec, types.Pool) {
	k := suite.app.LeveragelpKeeper
	suite.SetupCoinPrices(suite.ctx)
	poolAddr := ammtypes.NewPoolAddress(uint64(1))
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount := int64(10_000_000)
	pool := types.Pool{
		AmmPoolId:          1,
		Health:             math.LegacyZeroDec(),
		LeveragedLpAmount:  math.ZeroInt(),
		LeverageMax:        math.LegacyOneDec().MulInt64(10),
		MaxLeveragelpRatio: math.LegacyMustNewDecFromStr("0.6"),
	}
	poolInit := sdk.Coins{sdk.NewInt64Coin("uusdc", amount), sdk.NewInt64Coin("uusdt", amount)}

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolInit)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolInit)
	suite.Require().NoError(err)

	ammParams := suite.app.AmmKeeper.GetParams(suite.ctx)
	ammParams.WeightBreakingFeeMultiplier = math.LegacyZeroDec()
	ammParams.WeightBreakingFeeExponent = math.LegacyNewDecWithPrec(25, 1) // 2.5
	ammParams.WeightRecoveryFeePortion = math.LegacyNewDecWithPrec(10, 2)  // 10%
	ammParams.ThresholdWeightDifference = math.LegacyZeroDec()
	suite.app.AmmKeeper.SetParams(suite.ctx, ammParams)

	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId:            1,
		Address:           poolAddr.String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:   math.LegacyZeroDec(),
			UseOracle: true,
			FeeDenom:  "uusdc",
		},
		TotalShares: sdk.NewCoin("amm/pool/1", math.NewInt(2).Mul(ammtypes.OneShare)),
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:                  poolInit[0],
				Weight:                 math.NewInt(10),
				ExternalLiquidityRatio: math.LegacyOneDec(),
			},
			{
				Token:                  poolInit[1],
				Weight:                 math.NewInt(10),
				ExternalLiquidityRatio: math.LegacyOneDec(),
			},
		},
		TotalWeight: math.NewInt(20),
	})
	suite.Require().NoError(err)
	k.SetPool(suite.ctx, pool)
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdc",
		Liquidity: math.NewInt(amount),
	})
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdt",
		Liquidity: math.NewInt(amount),
	})
	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, 1)
	suite.Require().True(found)
	err = suite.app.PerpetualKeeper.OnLeverageLpEnablePool(suite.ctx, ammPool)

	usdcToken := sdk.NewInt64Coin("uusdc", amount*20)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, sdk.Coins{usdcToken})
	suite.Require().NoError(err)

	stableMsgServer := stablestakekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
	_, err = stableMsgServer.Bond(suite.ctx, &stablestaketypes.MsgBond{
		Creator: addr.String(),
		Amount:  math.NewInt(amount * 10),
		PoolId:  1,
	})
	suite.Require().NoError(err)

	leverage := math.LegacyNewDec(5)
	// open a position
	position, err := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: math.NewInt(amount).QuoRaw(1000),
		AmmPoolId:        1,
		Leverage:         leverage,
	}, 1)
	suite.Require().NoError(err)
	return position, leverage, pool
}

func (suite *KeeperTestSuite) TestForceCloseLong() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, leverage, pool := suite.OpenPosition(addr)
	timeDifference := suite.ctx.BlockTime().Add(time.Hour).Unix() - suite.ctx.BlockTime().Unix()
	interestRate := suite.app.StablestakeKeeper.GetParams(suite.ctx).InterestRate
	borrowed := leverage.Sub(math.LegacyOneDec()).MulInt(position.Collateral.Amount)
	repayAmount := borrowed.Add(borrowed.
		Mul(interestRate).
		Mul(math.LegacyNewDec(timeDifference)).
		Quo(math.LegacyNewDec(86400 * 365))).TruncateInt()

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))
	_, _, _, repayAmountOut, _, _, _, _, err := k.CheckHealthStopLossThenRepayAndClose(suite.ctx, position, &pool, math.LegacyOneDec(), false)
	suite.Require().NoError(err)
	suite.Require().Equal(repayAmount.String(), repayAmountOut.String())
}

func (suite *KeeperTestSuite) TestForceCloseLongWithNoFullRepayment() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, leverage, pool := suite.OpenPosition(addr)
	timeDifference := suite.ctx.BlockTime().Add(time.Hour*24*365*5).Unix() - suite.ctx.BlockTime().Unix()
	interestRate := suite.app.StablestakeKeeper.GetParams(suite.ctx).InterestRate
	borrowed := leverage.Sub(math.LegacyOneDec()).MulInt(position.Collateral.Amount)
	repayAmount := borrowed.Add(borrowed.
		Mul(interestRate).
		Mul(math.LegacyNewDec(timeDifference)).
		Quo(math.LegacyNewDec(86400 * 365))).RoundInt()

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 365 * 5))
	_, _, _, repayAmountOut, _, _, _, _, err := k.CheckHealthStopLossThenRepayAndClose(suite.ctx, position, &pool, math.LegacyOneDec(), false)
	suite.Require().NoError(err)
	suite.Require().Greater(repayAmount.String(), repayAmountOut.String())
}

func (suite *KeeperTestSuite) TestForceCloseLongPartial() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, leverage, pool := suite.OpenPosition(addr)
	originalPosition := *position
	timeDifference := suite.ctx.BlockTime().Add(time.Hour).Unix() - suite.ctx.BlockTime().Unix()
	interestRate := suite.app.StablestakeKeeper.GetParams(suite.ctx).InterestRate
	borrowed := leverage.Sub(math.LegacyOneDec()).MulInt(position.Collateral.Amount)
	repayAmount := borrowed.Add(borrowed.
		Mul(interestRate).
		Mul(math.LegacyNewDec(timeDifference)).
		Quo(math.LegacyNewDec(86400 * 365))).TruncateInt()
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))
	suite.SetCurrentHeight(suite.ctx.BlockHeight() + 1)
	// close 50%
	_, _, _, repayAmountOut, _, _, _, _, err := k.CheckHealthStopLossThenRepayAndClose(suite.ctx, position, &pool, math.LegacyOneDec().QuoInt64(2), false)
	suite.Require().NoError(err)
	suite.Require().Equal(repayAmount.Quo(math.NewInt(2)).String(), repayAmountOut.String())

	// Collateral should be reduced by 50%
	after, _ := k.GetPosition(suite.ctx, addr, 1)
	suite.Require().Equal(originalPosition.Collateral.Amount.Quo(math.NewInt(2)).String(), after.Collateral.Amount.String())
}

func (suite *KeeperTestSuite) TestHealthDecreaseForInterest() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _, _ := suite.OpenPosition(addr)
	_, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.221000000000000000") // slippage enabled on amm
	suite.Require().Equal("1.248428922744246025", health.String()) // slippage disabled on amm

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 365))
	suite.SetCurrentHeight(suite.ctx.BlockHeight() + 1)
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	suite.app.StablestakeKeeper.UpdateInterestAndGetDebt(suite.ctx, position.GetPositionAddress(), 1, 1)
	health, err = k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "0.610500000000000000") // slippage enabled on amm
	suite.Require().Equal("1.095113090126531601", health.String()) // slippage disabled on amm
}

// test positionHealth should be maxDec when liablities is zero
func (suite *KeeperTestSuite) TestPositionHealth() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _, _ := suite.OpenPosition(addr)
	_, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	suite.Require().Equal("1.248428922744246025", health.String())

	//setting position debt/liablities to zero
	debt := suite.app.StablestakeKeeper.GetDebt(suite.ctx, position.GetPositionAddress(), position.BorrowPoolId)
	debt.Borrowed = math.ZeroInt()
	debt.InterestStacked = math.ZeroInt()
	debt.InterestPaid = math.ZeroInt()
	suite.app.StablestakeKeeper.SetDebt(suite.ctx, debt)

	//get position health
	positionHealth, _ := suite.app.LeveragelpKeeper.GetPositionHealth(suite.ctx, *position)
	suite.Require().Equal(math.LegacyMaxSortableDec, positionHealth)
}
