package keeper_test

import (
	"sort"
	"strings"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
	"github.com/stretchr/testify/suite"
)

type assetPriceInfo struct {
	denom   string
	display string
	price   osmomath.BigDec
}

const (
	initChain = true
)

var (
	priceMap = map[string]assetPriceInfo{
		"uusdc": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   osmomath.OneBigDec(),
		},
		"uusdt": {
			denom:   "uusdt",
			display: "USDT",
			price:   osmomath.OneBigDec(),
		},
		"USDC": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   osmomath.OneBigDec(),
		},
		"uelys": {
			denom:   ptypes.Elys,
			display: "ELYS",
			price:   osmomath.MustNewBigDecFromStr("3.0"),
		},
		"uatom": {
			denom:   ptypes.ATOM,
			display: "ATOM",
			price:   osmomath.MustNewBigDecFromStr("1.0"),
		},
	}
)

type AmmKeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (suite *AmmKeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain, suite.Suite.T())

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain)
	suite.app = app
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(AmmKeeperTestSuite))
}

func (suite *AmmKeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *AmmKeeperTestSuite) GetAccountIssueAmount() math.Int {
	return math.NewInt(10_000_000_000_000)
}

func (suite *AmmKeeperTestSuite) AddAccounts(n int, given []sdk.AccAddress) []sdk.AccAddress {
	issueAmount := suite.GetAccountIssueAmount()
	var addresses []sdk.AccAddress
	if n > len(given) {
		addresses = simapp.AddTestAddrs(suite.app, suite.ctx, n-len(given), issueAmount)
		addresses = append(addresses, given...)
	} else {
		addresses = given
	}
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin(ptypes.ATOM, issueAmount),
			sdk.NewCoin(ptypes.Elys, issueAmount),
			sdk.NewCoin(ptypes.BaseCurrency, issueAmount),
		)
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
		if err != nil {
			panic(err)
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, address, coins)
		if err != nil {
			panic(err)
		}
	}
	return addresses
}

func (suite *AmmKeeperTestSuite) SetAmmParams() {
	suite.app.AmmKeeper.SetParams(suite.ctx, types.Params{
		PoolCreationFee:             math.NewInt(10_000_000),
		SlippageTrackDuration:       604800,
		BaseAssets:                  []string{ptypes.BaseCurrency},
		AllowedPoolCreators:         []string{authtypes.NewModuleAddress(govtypes.ModuleName).String()},
		WeightBreakingFeeExponent:   math.LegacyMustNewDecFromStr("2.5"),
		WeightBreakingFeeMultiplier: math.LegacyMustNewDecFromStr("0.0005"),
		WeightBreakingFeePortion:    math.LegacyMustNewDecFromStr("0.5"),
		WeightRecoveryFeePortion:    math.LegacyMustNewDecFromStr("0.1"),
		ThresholdWeightDifference:   math.LegacyMustNewDecFromStr("0.3"),
	})
}

func (suite *AmmKeeperTestSuite) SetupAssetProfile() {
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, atypes.Entry{
		BaseDenom:                "uusdc",
		Decimals:                 6,
		Denom:                    "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
		Path:                     "transfer/channel-12",
		IbcChannelId:             "channel-12",
		IbcCounterpartyChannelId: "channel-19",
		DisplayName:              "USDC",
		DisplaySymbol:            "uUSDC",
		Network:                  "",
		Address:                  "",
		ExternalSymbol:           "uUSDC",
		TransferLimit:            "",
		Permissions:              []string{},
		UnitDenom:                "uusdc",
		IbcCounterpartyDenom:     "",
		IbcCounterpartyChainId:   "",
		Authority:                "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
		CommitEnabled:            true,
		WithdrawEnabled:          true,
	})
}

func (suite *AmmKeeperTestSuite) SetupStableCoinPrices() {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Display: "USDC",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusdt",
		Display: "USDT",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusda",
		Display: "USDA",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     math.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     math.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDA",
		Price:     math.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
}

func (suite *AmmKeeperTestSuite) SetupCoinPrices() {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range priceMap {
		suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
			Denom:   v.denom,
			Display: v.display,
			Decimal: 6,
		})
		suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
			Asset:     v.display,
			Price:     v.price.Dec(),
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(suite.ctx.BlockTime().Unix()),
		})
	}
}

func (suite *AmmKeeperTestSuite) CreateNewAmmPool(creator sdk.AccAddress, useOracle bool, swapFee, exitFee osmomath.BigDec, asset2 string, baseTokenAmount, assetAmount math.Int) types.Pool {
	poolAssets := []types.PoolAsset{
		{
			Token:                  sdk.NewCoin(ptypes.BaseCurrency, baseTokenAmount),
			Weight:                 math.NewInt(10),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
		{
			Token:                  sdk.NewCoin(asset2, assetAmount),
			Weight:                 math.NewInt(10),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
	}
	sort.Slice(poolAssets, func(i, j int) bool {
		return strings.Compare(poolAssets[i].Token.Denom, poolAssets[j].Token.Denom) <= 0
	})
	poolParams := types.PoolParams{
		UseOracle: useOracle,
		SwapFee:   swapFee.Dec(),
		FeeDenom:  ptypes.BaseCurrency,
	}

	createPoolMsg := &types.MsgCreatePool{
		Sender:     creator.String(),
		PoolParams: poolParams,
		PoolAssets: poolAssets,
	}

	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, createPoolMsg)
	suite.Require().NoError(err)
	ammPool, _ := suite.app.AmmKeeper.GetPool(suite.ctx, poolId)

	return ammPool
}

func SetupMockPools(k *keeper.Keeper, ctx sdk.Context) {
	// Create and set mock pools
	pools := []types.Pool{
		{
			PoolId:  1,
			Address: types.NewPoolAddress(uint64(1)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("denom1", math.NewInt(1000)), Weight: math.OneInt()},
				{Token: sdk.NewCoin("denom2", math.NewInt(1000)), Weight: math.OneInt()},
			},
			TotalWeight: math.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(1), types.OneShare),
		},
		{
			PoolId:  2,
			Address: types.NewPoolAddress(uint64(2)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("uusdc", math.NewInt(1000)), Weight: math.OneInt()},
				{Token: sdk.NewCoin("denom1", math.NewInt(1000)), Weight: math.OneInt()},
			},
			TotalWeight: math.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(2), types.OneShare),
		},
		{
			PoolId:  3,
			Address: types.NewPoolAddress(uint64(3)).String(),
			PoolAssets: []types.PoolAsset{
				{Token: sdk.NewCoin("uusdc", math.NewInt(1000)), Weight: math.OneInt()},
				{Token: sdk.NewCoin("denom3", math.NewInt(1000)), Weight: math.OneInt()},
			},
			TotalWeight: math.NewInt(2),
			PoolParams: types.PoolParams{
				UseOracle: false,
			},
			TotalShares: sdk.NewCoin(types.GetPoolShareDenom(3), types.OneShare),
		},
	}

	for _, pool := range pools {
		k.SetPool(ctx, pool)
	}
}

func (suite *AmmKeeperTestSuite) VerifyPoolAssetWithBalance(poolId uint64) bool {
	pool, found := suite.app.AmmKeeper.GetPool(suite.ctx, poolId)
	if !found {
		return false
	}
	for _, asset := range pool.PoolAssets {
		bal := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(pool.Address), asset.Token.Denom)
		if !asset.Token.Amount.Equal(bal.Amount) {
			println("pool Asset DS: ", asset.Token.String())
			println("pool balance: ", bal.String())
			return false
		}
	}

	return true
}
