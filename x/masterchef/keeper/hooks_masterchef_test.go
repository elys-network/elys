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

func (suite *MasterchefKeeperTestSuite) TestHookMasterchef() {
	suite.ResetSuite(true)

	simapp.SetPerpetualParams(suite.app, suite.ctx)
	simapp.SetStableStake(suite.app, suite.ctx)
	simapp.SetParameters(suite.app, suite.ctx)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	srv := tokenomicskeeper.NewMsgServerImpl(suite.app.TokenomicsKeeper)

	expected := &tokenomicstypes.MsgCreateTimeBasedInflation{
		Description:      "Description",
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
	}

	_, err := srv.CreateTimeBasedInflation(suite.ctx, expected)
	suite.Require().NoError(err)

	expected = &tokenomicstypes.MsgCreateTimeBasedInflation{
		Description:      "Description",
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
	}
	_, err = srv.CreateTimeBasedInflation(suite.ctx, expected)
	suite.Require().NoError(err)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 2, math.NewInt(10000000000))

	// Create a pool
	// Mint 100000USDC
	suite.MintTokenToAddress(addr[0], math.NewInt(10000000000), ptypes.BaseCurrency)
	suite.MintTokenToAddress(addr[1], math.NewInt(10000000000), ptypes.BaseCurrency)

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

	// Mint uatom
	suite.MintTokenToAddress(addr[0], math.NewIntWithDecimal(100000000, 6), "uatom")

	// external reward distribution
	_, err = suite.msgServer.AddExternalRewardDenom(suite.ctx, &types.MsgAddExternalRewardDenom{
		Authority:   suite.app.GovKeeper.GetAuthority(),
		RewardDenom: "uatom",
		MinAmount:   math.OneInt(),
		Supported:   true,
	})
	suite.Require().NoError(err)
	_, err = suite.msgServer.AddExternalIncentive(suite.ctx, &types.MsgAddExternalIncentive{
		Sender:         addr[0].String(),
		RewardDenom:    "uatom",
		PoolId:         pools[0].PoolId,
		AmountPerBlock: math.NewIntWithDecimal(100, 6),
		FromBlock:      0,
		ToBlock:        1000,
	})
	suite.Require().NoError(err)

	// check rewards after 100 block
	ctx := suite.ctx
	for i := 1; i <= 100; i++ {
		suite.app.MasterchefKeeper.EndBlocker(ctx)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

	suite.Require().Equal(ctx.BlockHeight(), int64(100))
	poolRewardInfo, _ := suite.app.MasterchefKeeper.GetPoolRewardInfo(ctx, pools[0].PoolId, "uatom")
	suite.Require().Equal(poolRewardInfo.LastUpdatedBlock, uint64(99))

	res, err := suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(res.TotalRewards[0].Amount.String(), "4950000000")
	res, err = suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(res.TotalRewards[0].Amount.String(), "4950000000")

	// check rewards claimed
	_, err = suite.msgServer.ClaimRewards(ctx, &types.MsgClaimRewards{
		Sender:  addr[0].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	suite.Require().NoError(err)
	_, err = suite.msgServer.ClaimRewards(ctx, &types.MsgClaimRewards{
		Sender:  addr[1].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	suite.Require().NoError(err)

	atomAmount := suite.app.BankKeeper.GetBalance(ctx, addr[1], "uatom")
	suite.Require().Equal(atomAmount.Amount.String(), "4950000000")

	// no pending rewards
	res, err = suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(res.TotalRewards, 0)
	res, err = suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(res.TotalRewards, 0)

	// first user exit pool
	_, _, err = suite.app.AmmKeeper.ExitPool(ctx, addr[1], pools[0].PoolId, share.Quo(math.NewInt(2)), sdk.NewCoins(), "", false)
	suite.Require().NoError(err)

	// check rewards after 100 block
	for i := 1; i <= 100; i++ {
		suite.app.MasterchefKeeper.EndBlocker(ctx)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

	suite.Require().Equal(ctx.BlockHeight(), int64(200))
	poolRewardInfo, _ = suite.app.MasterchefKeeper.GetPoolRewardInfo(ctx, pools[0].PoolId, "uatom")
	suite.Require().Equal(poolRewardInfo.LastUpdatedBlock, uint64(199))

	res, err = suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(res.TotalRewards[0].String(), "6666666666uatom")
	res, err = suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(res.TotalRewards[0].String(), "3333333333uatom")

	// check rewards claimed
	_, err = suite.msgServer.ClaimRewards(ctx, &types.MsgClaimRewards{
		Sender:  addr[0].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	suite.Require().NoError(err)
	_, err = suite.msgServer.ClaimRewards(ctx, &types.MsgClaimRewards{
		Sender:  addr[1].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	suite.Require().NoError(err)

	atomAmount = suite.app.BankKeeper.GetBalance(ctx, addr[1], "uatom")
	suite.Require().Equal(atomAmount.String(), "8283333333uatom")

	// no pending rewards
	res, err = suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(res.TotalRewards, 0)
	res, err = suite.app.MasterchefKeeper.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(res.TotalRewards, 0)

	pool, found := suite.app.MasterchefKeeper.GetPoolInfo(ctx, pools[0].PoolId)
	suite.Require().Equal(true, found)
	suite.Require().Equal(pool.ExternalIncentiveApr.String(), "4204.799481351999973502")
}
