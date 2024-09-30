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

func (suite KeeperTestSuite) OpenPosition(addr sdk.AccAddress) (*types.Position, math.LegacyDec, types.Pool) {
	k := suite.app.LeveragelpKeeper
	suite.SetupCoinPrices(suite.ctx)
	poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount := int64(10_000_000)
	pool := types.Pool{
		AmmPoolId:         1,
		Enabled:           true,
		Closed:            false,
		Health:            math.LegacyZeroDec(),
		LeveragedLpAmount: math.ZeroInt(),
		LeverageMax:       math.LegacyOneDec().MulInt64(10),
	}
	poolInit := sdk.Coins{sdk.NewInt64Coin("uusdc", amount), sdk.NewInt64Coin("uusdt", amount)}

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolInit)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolInit)
	suite.Require().NoError(err)

	err = suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId:            1,
		Address:           poolAddr.String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:                     math.LegacyZeroDec(),
			ExitFee:                     math.LegacyZeroDec(),
			UseOracle:                   true,
			WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
			WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
			ExternalLiquidityRatio:      math.LegacyNewDec(1),
			WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
			ThresholdWeightDifference:   math.LegacyZeroDec(),
			FeeDenom:                    "uusdc",
		},
		TotalShares: sdk.NewCoin("amm/pool/1", math.NewInt(2).Mul(ammtypes.OneShare)),
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  poolInit[0],
				Weight: math.NewInt(10),
			},
			{
				Token:  poolInit[1],
				Weight: math.NewInt(10),
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

	usdcToken := sdk.NewInt64Coin("uusdc", amount*20)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, sdk.Coins{usdcToken})
	suite.Require().NoError(err)

	stableMsgServer := stablestakekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
	_, err = stableMsgServer.Bond(sdk.WrapSDKContext(suite.ctx), &stablestaketypes.MsgBond{
		Creator: addr.String(),
		Amount:  math.NewInt(amount * 10),
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
	})
	suite.Require().NoError(err)
	return position, leverage, pool
}

func (suite KeeperTestSuite) TestCloseLong() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	var (
		msg = &types.MsgClose{
			Creator:  addr.String(),
			Id:       1,
			LpAmount: math.ZeroInt(),
		}
		repayAmount = math.NewInt(0)
	)

	_, repayAmountOut, err := k.CloseLong(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Equal(repayAmount.String(), repayAmountOut.String())
}

func (suite KeeperTestSuite) TestForceCloseLong() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, leverage, pool := suite.OpenPosition(addr)
	timeDifference := suite.ctx.BlockTime().Add(time.Hour).Unix() - suite.ctx.BlockTime().Unix()
	interestRate := suite.app.StablestakeKeeper.GetParams(suite.ctx).InterestRate
	borrowed := leverage.Sub(math.LegacyOneDec()).MulInt(position.Collateral.Amount)
	repayAmount := borrowed.Add(borrowed.
		Mul(interestRate).
		Mul(math.LegacyNewDec(timeDifference)).
		Quo(math.LegacyNewDec(86400 * 365))).RoundInt()

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))
	repayAmountOut, err := k.ForceCloseLong(suite.ctx, *position, pool, position.LeveragedLpAmount)
	suite.Require().NoError(err)
	suite.Require().Equal(repayAmount.String(), repayAmountOut.String())
}

func (suite KeeperTestSuite) TestForceCloseLongWithNoFullRepayment() {
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
	repayAmountOut, err := k.ForceCloseLong(suite.ctx, *position, pool, position.LeveragedLpAmount)
	suite.Require().NoError(err)
	suite.Require().Greater(repayAmount.String(), repayAmountOut.String())
}

func (suite KeeperTestSuite) TestForceCloseLongPartial() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, leverage, pool := suite.OpenPosition(addr)
	timeDifference := suite.ctx.BlockTime().Add(time.Hour).Unix() - suite.ctx.BlockTime().Unix()
	interestRate := suite.app.StablestakeKeeper.GetParams(suite.ctx).InterestRate
	borrowed := leverage.Sub(math.LegacyOneDec()).MulInt(position.Collateral.Amount)
	repayAmount := borrowed.Add(borrowed.
		Mul(interestRate).
		Mul(math.LegacyNewDec(timeDifference)).
		Quo(math.LegacyNewDec(86400 * 365))).RoundInt()
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))
	// close 50%
	repayAmountOut, err := k.ForceCloseLong(suite.ctx, *position, pool, position.LeveragedLpAmount.Quo(math.NewInt(2)))
	suite.Require().NoError(err)
	suite.Require().Equal(repayAmount.Quo(math.NewInt(2)).String(), repayAmountOut.String())
}

func (suite KeeperTestSuite) TestHealthDecreaseForInterest() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _, _ := suite.OpenPosition(addr)
	_, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "1.221000000000000000") // slippage enabled on amm
	suite.Require().Equal("1.250000000000000000", health.String()) // slippage disabled on amm

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 365))
	suite.app.StablestakeKeeper.BeginBlocker(suite.ctx)
	suite.app.StablestakeKeeper.UpdateInterestAndGetDebt(suite.ctx, position.GetPositionAddress())
	health, err = k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	// suite.Require().Equal(health.String(), "0.610500000000000000") // slippage enabled on amm
	suite.Require().Equal("1.096491228070175439", health.String()) // slippage disabled on amm
}

// test positionHealth should be maxDec when liablities is zero
func (suite KeeperTestSuite) TestPositionHealth() {
	k := suite.app.LeveragelpKeeper
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	position, _, _ := suite.OpenPosition(addr)
	_, found := suite.app.AmmKeeper.GetPool(suite.ctx, position.AmmPoolId)
	suite.Require().True(found)
	health, err := k.GetPositionHealth(suite.ctx, *position)
	suite.Require().NoError(err)
	suite.Require().Equal("1.250000000000000000", health.String())

	//setting position debt/liablities to zero
	debt := suite.app.StablestakeKeeper.GetDebt(suite.ctx, position.GetPositionAddress())
	debt.Borrowed = math.ZeroInt()
	debt.InterestStacked = math.ZeroInt()
	debt.InterestPaid = math.ZeroInt()
	suite.app.StablestakeKeeper.SetDebt(suite.ctx, debt)

	//get position health
	positionHealth, _ := suite.app.LeveragelpKeeper.GetPositionHealth(suite.ctx, *position)
	suite.Require().Equal(math.LegacyMaxSortableDec, positionHealth)
}
