package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
)

func (suite *MasterchefKeeperTestSuite) TestExternalIncentive() {
	externalIncentives := []types.ExternalIncentive{
		{
			Id:             0,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: math.OneInt(),
			Apr:            math.LegacyZeroDec(),
		},
		{
			Id:             1,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: math.OneInt(),
			Apr:            math.LegacyZeroDec(),
		},
		{
			Id:             2,
			RewardDenom:    "reward1",
			PoolId:         2,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: math.OneInt(),
			Apr:            math.LegacyZeroDec(),
		},
	}
	for _, externalIncentive := range externalIncentives {
		suite.app.MasterchefKeeper.SetExternalIncentive(suite.ctx, externalIncentive)
	}
	for _, externalIncentive := range externalIncentives {
		info, found := suite.app.MasterchefKeeper.GetExternalIncentive(suite.ctx, externalIncentive.Id)
		suite.Require().True(found)
		suite.Require().Equal(info, externalIncentive)
	}
	externalIncentivesStored := suite.app.MasterchefKeeper.GetAllExternalIncentives(suite.ctx)
	suite.Require().Len(externalIncentivesStored, 3)

	suite.app.MasterchefKeeper.RemoveExternalIncentive(suite.ctx, externalIncentives[0].Id)
	externalIncentivesStored = suite.app.MasterchefKeeper.GetAllExternalIncentives(suite.ctx)
	suite.Require().Len(externalIncentivesStored, 2)
}

// Test USDC reward as external and via dex collection
func (suite *MasterchefKeeperTestSuite) TestUSDCExternalIncentive() {

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	srv := tokenomicskeeper.NewMsgServerImpl(suite.app.TokenomicsKeeper)

	expected := &tokenomicstypes.MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: uint64(1),
		EndBlockHeight:   uint64(6307200),
		Inflation: &tokenomicstypes.InflationEntry{
			LmRewards:         9999999,
			IcsStakingRewards: 9999999,
			CommunityFund:     9999999,
			StrategicReserve:  9999999,
			TeamTokensVested:  9999999,
		},
		Description: "Description",
	}

	_, err := srv.CreateTimeBasedInflation(suite.ctx, expected)
	suite.Require().NoError(err)

	expected = &tokenomicstypes.MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: uint64(6307201),
		EndBlockHeight:   uint64(12614401),
		Inflation: &tokenomicstypes.InflationEntry{
			LmRewards:         9999999,
			IcsStakingRewards: 9999999,
			CommunityFund:     9999999,
			StrategicReserve:  9999999,
			TeamTokensVested:  9999999,
		},
		Description: "Description",
	}
	_, err = srv.CreateTimeBasedInflation(suite.ctx, expected)
	suite.Require().NoError(err)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 2, math.NewInt(10000000000))

	// Create a pool
	// Mint 100000USDC
	suite.MintTokenToAddress(addr[0], math.NewInt(10000000000), ptypes.BaseCurrency)
	suite.MintTokenToAddress(addr[1], math.NewInt(10000000000), ptypes.BaseCurrency)

	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, ammtypes.ModuleName, usdcToken.MulInt(math.NewInt(2)))
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, ammtypes.ModuleName, addr[0], usdcToken)
	suite.Require().NoError(err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, math.NewInt(10000000)),
		},
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000000)),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := ammtypes.PoolParams{
		SwapFee: argSwapFee,
	}

	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(sdk.AccAddress(addr[0]), poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)

	// check length of pools
	suite.Require().Equal(len(pools), 1)

	_, _, err = suite.app.AmmKeeper.ExitPool(suite.ctx, addr[0], pools[0].PoolId, math.NewIntWithDecimal(1, 21), sdk.NewCoins(), "", false)
	suite.Require().NoError(err)

	// new user join pool with same shares
	share := ammtypes.InitPoolSharesSupply.Mul(math.NewIntWithDecimal(1, 5))
	suite.T().Log(suite.app.MasterchefKeeper.GetPoolTotalCommit(suite.ctx, pools[0].PoolId))
	suite.Require().Equal(suite.app.MasterchefKeeper.GetPoolTotalCommit(suite.ctx, pools[0].PoolId).String(), "10000000000000000000000000")
	suite.Require().Equal(suite.app.MasterchefKeeper.GetPoolBalance(suite.ctx, pools[0].PoolId, addr[0]).String(), "10000000000000000000000000")
	_, _, err = suite.app.AmmKeeper.JoinPoolNoSwap(suite.ctx, addr[1], pools[0].PoolId, share, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, math.NewInt(10000000)), sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000000))))
	suite.Require().NoError(err)
	suite.Require().Equal(suite.app.MasterchefKeeper.GetPoolTotalCommit(suite.ctx, pools[0].PoolId).String(), "20000000000000000000000000")
	suite.Require().Equal(suite.app.MasterchefKeeper.GetPoolBalance(suite.ctx, pools[0].PoolId, addr[1]), share)

	suite.MintTokenToAddress(addr[0], math.NewIntWithDecimal(100000000, 6), ptypes.ATOM)

	_, err = suite.msgServer.AddExternalRewardDenom(suite.ctx, &types.MsgAddExternalRewardDenom{
		Authority:   suite.app.GovKeeper.GetAuthority(),
		RewardDenom: ptypes.BaseCurrency,
		MinAmount:   math.OneInt(),
		Supported:   true,
	})
	suite.Require().NoError(err)
	_, err = suite.msgServer.AddExternalIncentive(suite.ctx, &types.MsgAddExternalIncentive{
		Sender:         addr[0].String(),
		RewardDenom:    ptypes.BaseCurrency,
		PoolId:         pools[0].PoolId,
		AmountPerBlock: math.NewIntWithDecimal(100, 6),
		FromBlock:      0,
		ToBlock:        1000,
	})
	suite.Require().NoError(err)

	info, _ := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, pools[0].PoolId)
	info.EnableEdenRewards = true
	suite.app.MasterchefKeeper.SetPoolInfo(suite.ctx, info)

	// Fill in pool revenue wallet
	revenueAddress1 := ammtypes.NewPoolRevenueAddress(1)
	suite.MintTokenToAddress(revenueAddress1, math.NewInt(100000), ptypes.BaseCurrency)

	// check rewards after 100 block
	for i := 1; i <= 100; i++ {
		suite.app.MasterchefKeeper.EndBlocker(suite.ctx)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	}

	suite.Require().Equal(suite.ctx.BlockHeight(), int64(100))
	poolRewardInfo, _ := suite.app.MasterchefKeeper.GetPoolRewardInfo(suite.ctx, pools[0].PoolId, ptypes.BaseCurrency)
	suite.Require().Equal(poolRewardInfo.LastUpdatedBlock, uint64(99))

	res, err := suite.app.MasterchefKeeper.UserPendingReward(suite.ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(res.TotalRewards[0].Amount.String(), "49")
	suite.Require().Equal(res.TotalRewards[1].Amount.String(), "4950030000")
	res, err = suite.app.MasterchefKeeper.UserPendingReward(suite.ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(res.TotalRewards[0].Amount.String(), "49")
	suite.Require().Equal(res.TotalRewards[1].Amount.String(), "4950030000")

	prevUSDCBal := suite.app.BankKeeper.GetBalance(suite.ctx, addr[1], ptypes.BaseCurrency)

	// check rewards claimed
	_, err = suite.msgServer.ClaimRewards(suite.ctx, &types.MsgClaimRewards{
		Sender:  addr[0].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	suite.Require().NoError(err)
	_, err = suite.msgServer.ClaimRewards(suite.ctx, &types.MsgClaimRewards{
		Sender:  addr[1].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	suite.Require().NoError(err)

	curUSDCBal := suite.app.BankKeeper.GetBalance(suite.ctx, addr[1], ptypes.BaseCurrency)
	amount, _ := math.NewIntFromString("4950030000")
	suite.Require().Equal(curUSDCBal.Amount.String(), prevUSDCBal.Amount.Add(amount).String())

	// no pending rewards
	res, err = suite.app.MasterchefKeeper.UserPendingReward(suite.ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(res.TotalRewards, 0)
	res, err = suite.app.MasterchefKeeper.UserPendingReward(suite.ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(res.TotalRewards, 0)

	// Eden should not be claimable
	curEdenBal := suite.app.BankKeeper.GetBalance(suite.ctx, addr[1], ptypes.Eden)
	suite.Require().Equal(curEdenBal.Amount.String(), "0")
}
