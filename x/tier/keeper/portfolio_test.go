package keeper_test

import (
	"cosmossdk.io/math"
	"strconv"
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"

	assetprofilerkeeper "github.com/elys-network/elys/x/assetprofile/keeper"
	profiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/tier/keeper"
	"github.com/elys-network/elys/x/tier/types"

	"github.com/elys-network/elys/app"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

//const userAddress string = "cosmos10t3g865e53yhhzvwwr5ldg50yq7vdwwfemrdeg"

func createNPortfolio(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Portfolio {
	items := make([]types.Portfolio, n)
	addresses := app.CreateRandomAccounts(n)
	for i := range items {
		items[i].Creator = addresses[i].String()
		items[i].Portfolio = math.LegacyNewDec(1000)
		items[i].Date = keeper.GetDateFromContext(ctx)

		keeper.SetPortfolio(ctx, items[i])
	}
	return items
}

func TestPortfolioGet(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	for _, item := range items {
		_, found := keeper.GetPortfolio(ctx,
			sdk.MustAccAddressFromBech32(item.Creator),
			keeper.GetDateFromContext(ctx),
		)
		require.True(t, found)
	}
}

func TestPortfolioRemoveLast(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	count := keeper.RemovePortfolioLast(ctx,
		keeper.GetDateFromContext(ctx),
		100,
	)
	_, found := keeper.GetPortfolio(ctx,
		sdk.MustAccAddressFromBech32(items[9].Creator),
		keeper.GetDateFromContext(ctx),
	)
	require.Equal(t, count, uint64(10))
	require.False(t, found)

	// Try to remove again
	count = keeper.RemovePortfolioLast(ctx,
		keeper.GetDateFromContext(ctx),
		100,
	)
	require.Equal(t, count, uint64(0))
}

func TestPortfolioGetAll(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPortfolio(ctx)))
}

func TestGetPortfolioNative(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	_, _, oracle, tier, assetProfiler := app.MasterchefKeeper, app.AmmKeeper, app.OracleKeeper, app.TierKeeper, app.AssetprofileKeeper

	// Setup coin prices
	SetupCoinPrices(ctx, oracle, assetProfiler)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(0))

	// Create a pool
	// Mint 100000USDC + 10 ELYS (pool creation fee)
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000))
	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], coins)
	require.NoError(t, err)

	tier.RetrieveAllPortfolio(ctx, addr[0])

	portfolio, found := tier.GetPortfolio(ctx, addr[0], tier.GetDateFromContext(ctx))
	require.True(t, found)
	require.Equal(t, portfolio, math.LegacyNewDec(101000))
}

func TestGetPortfolioAmm(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	_, amm, oracle, tier, assetProfiler := app.MasterchefKeeper, app.AmmKeeper, app.OracleKeeper, app.TierKeeper, app.AssetprofileKeeper

	// Setup coin prices
	SetupCoinPrices(ctx, oracle, assetProfiler)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000))

	// Create a pool
	// Mint 100000USDC + 10 ELYS (pool creation fee)
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000))
	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], coins)
	require.NoError(t, err)

	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: math.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, math.NewInt(100000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: math.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
	})

	poolParams := &ammtypes.PoolParams{
		SwapFee:                     math.LegacyZeroDec(),
		ExitFee:                     math.LegacyZeroDec(),
		UseOracle:                   false,
		WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
		WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
		ExternalLiquidityRatio:      math.LegacyOneDec(),
		WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   math.LegacyZeroDec(),
		FeeDenom:                    "",
	}

	// Create a Elys+USDC pool
	msgServer := ammkeeper.NewMsgServerImpl(amm)
	resp, err := msgServer.CreatePool(
		sdk.WrapSDKContext(ctx),
		&ammtypes.MsgCreatePool{
			Sender:     addr[0].String(),
			PoolParams: poolParams,
			PoolAssets: poolAssets,
		})

	// TODO: Check lp token price
	//pool := amm.GetAllPool(ctx)[0]
	//info := amm.PoolExtraInfo(ctx, pool)
	//require.Equal(t, pool.TotalShares, pool)

	//require.Equal(t, info.Tvl, math.LegacyNewDec(2))

	//require.Equal(t, info.LpTokenPrice, math.LegacyNewDec(2))

	require.NoError(t, err)
	require.Equal(t, resp.PoolID, uint64(1))

	tier.RetrieveAllPortfolio(ctx, addr[0])

	portfolio, found := tier.GetPortfolio(ctx, addr[0], tier.GetDateFromContext(ctx))
	require.True(t, found)
	require.Equal(t, portfolio, math.LegacyNewDec(100100))
}

func TestPortfolioGetDiscount(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := make([]types.Portfolio, 10)
	addresses := app.CreateRandomAccounts(10)
	for j := 0; j < 8; j++ {
		ctx = ctx.WithBlockTime(ctx.BlockTime().AddDate(0, 0, 1))
		for i := range items {
			items[i].Creator = addresses[i].String()
			items[i].Portfolio = math.LegacyNewDec(400000)
			items[i].Date = keeper.GetDateFromContext(ctx)

			keeper.SetPortfolio(ctx, items[i])
		}
	}

	items[9].Portfolio = math.LegacyNewDec(500)
	items[9].Date = keeper.GetDateFromContext(ctx)

	keeper.SetPortfolio(ctx, items[9])

	_, _, discount := keeper.GetMembershipTier(ctx, sdk.MustAccAddressFromBech32(items[0].Creator))
	require.Equal(t, discount, uint64(20))

	_, _, discount = keeper.GetMembershipTier(ctx, sdk.MustAccAddressFromBech32(items[9].Creator))
	require.Equal(t, discount, uint64(0))
}

func TestGetPortfolioPerpetual(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	perpetual, amm, oracle, tier, assetProfiler := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper, app.TierKeeper, app.AssetprofileKeeper

	// Setup coin prices
	SetupCoinPrices(ctx, oracle, assetProfiler)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000))

	// Create a pool
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 1000000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 10000000))
	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], coins)
	require.NoError(t, err)

	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: math.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, math.NewInt(10000000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: math.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000)),
	})

	poolParams := &ammtypes.PoolParams{
		SwapFee:                     math.LegacyZeroDec(),
		ExitFee:                     math.LegacyZeroDec(),
		UseOracle:                   false,
		WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
		WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
		ExternalLiquidityRatio:      math.LegacyOneDec(),
		WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   math.LegacyZeroDec(),
		FeeDenom:                    "",
	}

	// Create a Elys+USDC pool
	msgServer := ammkeeper.NewMsgServerImpl(amm)
	resp, err := msgServer.CreatePool(
		sdk.WrapSDKContext(ctx),
		&ammtypes.MsgCreatePool{
			Sender:     addr[0].String(),
			PoolParams: poolParams,
			PoolAssets: poolAssets,
		})

	require.NoError(t, err)
	require.Equal(t, resp.PoolID, uint64(1))

	err = perpetual.SetMTP(ctx, &perpetualtypes.MTP{
		Address:                        addr[0].String(),
		CollateralAsset:                ptypes.BaseCurrency,
		CustodyAsset:                   ptypes.Elys,
		Collateral:                     math.NewInt(0),
		Liabilities:                    math.NewInt(0),
		BorrowInterestPaidCollateral:   math.NewInt(0),
		BorrowInterestPaidCustody:      math.NewInt(0),
		BorrowInterestUnpaidCollateral: math.NewInt(0),
		Custody:                        math.NewInt(10000),
		MtpHealth:                      math.LegacyNewDec(0),
		Position:                       perpetualtypes.Position_LONG,
		Id:                             0,
		ConsolidateLeverage:            math.LegacyZeroDec(),
		SumCollateral:                  math.ZeroInt(),
	})
	require.NoError(t, err)

	tier.RetrieveAllPortfolio(ctx, addr[0])

	portfolio, found := tier.GetPortfolio(ctx, addr[0], tier.GetDateFromContext(ctx))
	require.True(t, found)
	require.Equal(t, portfolio, math.LegacyNewDec(10099100))
}

// TODO
// 3: staked + rewards

func SetupCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper, assetProfiler assetprofilerkeeper.Keeper) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Display: "USDC",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   "uusdt",
		Display: "USDT",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.Elys,
		Display: "ELYS",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.ATOM,
		Display: "ATOM",
		Decimal: 6,
	})
	assetProfiler.SetEntry(ctx, profiletypes.Entry{BaseDenom: ptypes.Elys, Denom: ptypes.Elys})
	assetProfiler.SetEntry(ctx, profiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency})

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     math.LegacyNewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     math.LegacyNewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     math.LegacyNewDec(100),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     math.LegacyNewDec(100),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
}
