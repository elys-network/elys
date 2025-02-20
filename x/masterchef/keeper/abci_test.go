package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"

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

func (suite *MasterchefKeeperTestSuite) TestABCI_EndBlocker() {
	suite.ResetSuite(true)

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(2000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, sdkmath.NewInt(2000))
	unclaimed = unclaimed.Add(uedenToken, uedenBToken)

	// Mint coins
	suite.MintMultipleTokenToAddress(suite.genesisAccount, unclaimed)

	// Add testing commitment
	simapp.AddTestCommitment(suite.app, suite.ctx, suite.genesisAccount, committed)
	suite.app.MasterchefKeeper.EndBlocker(suite.ctx)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	srv := tokenomicskeeper.NewMsgServerImpl(suite.app.TokenomicsKeeper)

	expected := &tokenomicstypes.MsgCreateTimeBasedInflation{
		Description:      "description",
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
		Description:      "description",
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

	// Set tokenomics params
	listTimeBasdInflations := suite.app.TokenomicsKeeper.GetAllTimeBasedInflation(suite.ctx)

	// After the first year
	ctx := suite.ctx.WithBlockHeight(1)
	suite.app.MasterchefKeeper.ProcessUpdateIncentiveParams(ctx)

	// Check if the params are correctly set
	params := suite.app.MasterchefKeeper.GetParams(ctx)
	suite.Require().NotNil(params.LpIncentives)
	suite.Require().Equal(params.LpIncentives.EdenAmountPerYear, sdkmath.NewInt(int64(listTimeBasdInflations[0].Inflation.LmRewards)))

	// After the first year
	ctx = ctx.WithBlockHeight(6307210)

	// After reading tokenomics again
	suite.app.MasterchefKeeper.ProcessUpdateIncentiveParams(ctx)

	// Check if the params are correctly set
	params = suite.app.MasterchefKeeper.GetParams(ctx)
	suite.Require().NotNil(params.LpIncentives)
	suite.Require().Equal(params.LpIncentives.EdenAmountPerYear, sdkmath.NewInt(int64(listTimeBasdInflations[0].Inflation.LmRewards)))
}

func (suite *MasterchefKeeperTestSuite) TestCollectGasFees() {
	suite.ResetSuite(true)
	// Collect gas fees
	collectedAmt, err := suite.app.MasterchefKeeper.CollectGasFees(suite.ctx, ptypes.BaseCurrency)
	suite.Require().NoError(err)

	// rewards should be zero
	suite.Require().True(collectedAmt.IsZero())

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdkmath.NewInt(1000000))
	transferAmt := sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100))

	// Set revenue address
	params := suite.app.MasterchefKeeper.GetParams(suite.ctx)
	params.ProtocolRevenueAddress = addr[0].String()
	suite.app.MasterchefKeeper.SetParams(suite.ctx, params)

	// Deposit 100elys to FeeCollectorName wallet
	err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, addr[0], authtypes.FeeCollectorName, sdk.NewCoins(transferAmt))
	suite.Require().NoError(err)

	// Create a pool
	// Mint 100000USDC
	suite.MintTokenToAddress(addr[0], sdkmath.NewInt(100000), ptypes.BaseCurrency)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000)),
		},
	}

	argSwapFee := sdkmath.LegacyMustNewDecFromStr("0.01")

	poolParams := ammtypes.PoolParams{
		SwapFee: argSwapFee,
	}

	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(addr[0], poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)

	// check length of pools
	suite.Require().Equal(len(pools), 1)

	// check block height
	suite.Require().Equal(int64(0), suite.ctx.BlockHeight())

	// Collect gas fees again
	collectedAmt, err = suite.app.MasterchefKeeper.CollectGasFees(suite.ctx, ptypes.BaseCurrency)
	suite.Require().NoError(err)

	// check block height
	suite.Require().Equal(int64(0), suite.ctx.BlockHeight())

	// It should be 5.4 usdc
	suite.Require().Equal(collectedAmt.String(), "5.400000000000000000uusdc")
}

func (suite *MasterchefKeeperTestSuite) TestCollectDEXRevenue() {

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 2, sdkmath.NewInt(1000000))

	// Create 2 pools

	// #######################
	// ####### POOL 1 ########
	// Mint 100000USDC
	suite.MintTokenToAddress(addr[0], sdkmath.NewInt(100000), ptypes.BaseCurrency)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000)),
		},
	}

	argSwapFee := sdkmath.LegacyMustNewDecFromStr("0.01")

	poolParams := ammtypes.PoolParams{
		SwapFee: argSwapFee,
	}

	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(sdk.AccAddress(addr[0]), poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	// ####### POOL 2 ########
	// ATOM+USDC pool
	// Mint uusdc
	suite.MintTokenToAddress(addr[1], sdkmath.NewInt(200000), ptypes.BaseCurrency)
	// Mint uatom
	suite.MintTokenToAddress(addr[1], sdkmath.NewInt(200000), ptypes.ATOM)

	poolAssets2 := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(150000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000)),
		},
	}
	// Create a ATOM+USDC pool
	ammPool = suite.CreateNewAmmPool(sdk.AccAddress(addr[1]), poolAssets2, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(2))

	pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)

	// check length of pools
	suite.Require().Equal(len(pools), 2)

	// check block height
	suite.Require().Equal(int64(0), suite.ctx.BlockHeight())

	// Fill in pool #1 revenue wallet
	revenueAddress1 := ammtypes.NewPoolRevenueAddress(0)
	suite.MintTokenToAddress(revenueAddress1, sdkmath.NewInt(1000), ptypes.BaseCurrency)

	// Fill in pool #2 revenue wallet
	revenueAddress2 := ammtypes.NewPoolRevenueAddress(1)
	suite.MintTokenToAddress(revenueAddress2, sdkmath.NewInt(3000), ptypes.BaseCurrency)

	// Set revenue address
	params := suite.app.MasterchefKeeper.GetParams(suite.ctx)
	params.ProtocolRevenueAddress = addr[0].String()
	suite.app.MasterchefKeeper.SetParams(suite.ctx, params)

	// Collect revenue
	collectedAmt, rewardForLpsAmt, _, err := suite.app.MasterchefKeeper.CollectDEXRevenue(suite.ctx)
	suite.Require().NoError(err)

	// check block height
	suite.Require().Equal(int64(0), suite.ctx.BlockHeight())

	// It should be 3000=1000+2000 usdc
	suite.Require().Equal(collectedAmt, sdk.Coins{sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(3000))})
	// It should be 1950=3000*0.65 usdc
	suite.Require().Equal(rewardForLpsAmt, sdk.DecCoins{sdk.NewDecCoin(ptypes.BaseCurrency, sdkmath.NewInt(1800))})
}

func (suite *MasterchefKeeperTestSuite) TestExternalRewardsDistribution() {
	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 2, sdkmath.NewInt(1000000))

	// Create 2 pools

	// #######################
	// ####### POOL 1 ########
	// Mint 100000USDC
	suite.MintTokenToAddress(addr[0], sdkmath.NewInt(100000), ptypes.BaseCurrency)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000)),
		},
	}

	argSwapFee := sdkmath.LegacyMustNewDecFromStr("0.01")

	poolParams := ammtypes.PoolParams{
		SwapFee: argSwapFee,
	}

	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(sdk.AccAddress(addr[0]), poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	// ####### POOL 2 ########
	// ATOM+USDC pool
	// Mint uusdc
	suite.MintTokenToAddress(addr[1], sdkmath.NewInt(200000), ptypes.BaseCurrency)
	// Mint uatom
	suite.MintTokenToAddress(addr[1], sdkmath.NewInt(200000), ptypes.ATOM)

	poolAssets2 := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(150000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000)),
		},
	}

	// Create a ATOM+USDC pool
	ammPool = suite.CreateNewAmmPool(sdk.AccAddress(addr[1]), poolAssets2, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(2))

	pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)

	// check length of pools
	suite.Require().Equal(len(pools), 2)

	externalIncentive := types.ExternalIncentive{
		Id:             0,
		RewardDenom:    "reward1",
		PoolId:         1,
		FromBlock:      suite.ctx.BlockHeight() - 1,
		ToBlock:        suite.ctx.BlockHeight() + 101,
		AmountPerBlock: sdkmath.OneInt(),
		Apr:            sdkmath.LegacyZeroDec(),
	}

	suite.app.MasterchefKeeper.SetExternalIncentive(suite.ctx, externalIncentive)

	_, found := suite.app.MasterchefKeeper.GetPoolRewardInfo(suite.ctx, externalIncentive.PoolId, externalIncentive.RewardDenom)
	suite.Require().False(found)

	suite.app.MasterchefKeeper.ProcessExternalRewardsDistribution(suite.ctx)

	pool, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, externalIncentive.PoolId)
	suite.Require().True(found)
	suite.Require().Equal(pool.ExternalRewardDenoms, []string{"reward1"})

	rewardInfo, found := suite.app.MasterchefKeeper.GetPoolRewardInfo(suite.ctx, externalIncentive.PoolId, externalIncentive.RewardDenom)
	suite.Require().True(found)
	suite.Require().Equal(rewardInfo.RewardDenom, externalIncentive.RewardDenom)
	suite.Require().Equal(rewardInfo.PoolAccRewardPerShare, sdkmath.LegacyMustNewDecFromStr("0.000099900099900099"))

	// Test multiple external incentives
	externalIncentive2 := types.ExternalIncentive{
		Id:             0,
		RewardDenom:    "reward2",
		PoolId:         1,
		FromBlock:      suite.ctx.BlockHeight() - 1,
		ToBlock:        suite.ctx.BlockHeight() + 101,
		AmountPerBlock: sdkmath.OneInt(),
		Apr:            sdkmath.LegacyZeroDec(),
	}
	suite.app.MasterchefKeeper.SetExternalIncentive(suite.ctx, externalIncentive2)

	suite.app.MasterchefKeeper.ProcessExternalRewardsDistribution(suite.ctx)

	pool, found = suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, externalIncentive.PoolId)
	suite.Require().True(found)
	suite.Require().Equal(pool.ExternalRewardDenoms, []string{"reward1", "reward2"})

	rewardInfo, found = suite.app.MasterchefKeeper.GetPoolRewardInfo(suite.ctx, externalIncentive2.PoolId, externalIncentive2.RewardDenom)
	suite.Require().True(found)
	suite.Require().Equal(rewardInfo.RewardDenom, externalIncentive2.RewardDenom)
	suite.Require().Equal(rewardInfo.PoolAccRewardPerShare, sdkmath.LegacyMustNewDecFromStr("0.000099900099900099"))

	// Get Tvl for non-existent pool
	res := suite.app.MasterchefKeeper.GetPoolTVL(suite.ctx, 1000)
	suite.Require().Equal(res, sdkmath.LegacyZeroDec())

	// increase timestamp
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour))
	suite.app.MasterchefKeeper.ProcessExternalRewardsDistribution(suite.ctx)

	// set current block to last block
	suite.ctx = suite.ctx.WithBlockHeight(externalIncentive.ToBlock)

	suite.app.MasterchefKeeper.ProcessExternalRewardsDistribution(suite.ctx)

	// should delete incentives
	listInc := suite.app.MasterchefKeeper.GetAllExternalIncentives(suite.ctx)
	suite.Require().Equal(len(listInc), 0)
}

func (suite *MasterchefKeeperTestSuite) TestInitialParams() {
	suite.ResetSuite(true)

	res := suite.app.MasterchefKeeper.InitPoolParams(suite.ctx, 1)
	suite.Require().Equal(res, true)

	poolInfo, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, 1)
	suite.Require().Equal(found, true)
	suite.Require().Equal(poolInfo.PoolId, uint64(1))
}

func (suite *MasterchefKeeperTestSuite) TestProcessTakerFees() {
	suite.ResetSuite(true)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdkmath.NewInt(1000000))

	// mint some tokens in taker address
	takerAddress := suite.app.ParameterKeeper.GetParams(suite.ctx).TakerFeeCollectionAddress
	suite.MintTokenToAddress(sdk.MustAccAddressFromBech32(takerAddress), sdkmath.NewInt(100), ptypes.BaseCurrency)
	suite.MintTokenToAddress(addr[0], sdkmath.NewInt(100000), ptypes.BaseCurrency)

	// Pool with 1000 ELYS and 1000 USDC
	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000)),
		},
	}

	argSwapFee := sdkmath.LegacyMustNewDecFromStr("0.01")

	poolParams := ammtypes.PoolParams{
		SwapFee: argSwapFee,
	}

	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(addr[0], poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	pools := suite.app.AmmKeeper.GetAllPool(suite.ctx)
	suite.Require().Equal(len(pools), 1)

	elysSupplyBefore := suite.app.BankKeeper.GetSupply(suite.ctx, ptypes.Elys)
	suite.Require().Equal(elysSupplyBefore.Amount.String(), "100000002000000")

	// Process taker fees
	suite.app.MasterchefKeeper.ProcessTakerFee(suite.ctx)

	// Check elys supply is reduced
	elysSupplyAfter := suite.app.BankKeeper.GetSupply(suite.ctx, ptypes.Elys)
	suite.Require().Equal(elysSupplyAfter.Amount.String(), "100000002000000")
}
