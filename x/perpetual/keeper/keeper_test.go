package keeper_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type assetPriceInfo struct {
	denom   string
	display string
	price   sdk.Dec
}

const (
	initChain = true
)

var (
	priceMap = map[string]assetPriceInfo{
		"uusdc": {
			denom:   ptypes.BaseCurrency,
			display: "USDC",
			price:   sdk.OneDec(),
		},
		"uusdt": {
			denom:   "uusdt",
			display: "USDT",
			price:   sdk.OneDec(),
		},
		"uelys": {
			denom:   ptypes.Elys,
			display: "ELYS",
			price:   sdk.MustNewDecFromStr("3.0"),
		},
		"uatom": {
			denom:   ptypes.ATOM,
			display: "ATOM",
			price:   sdk.MustNewDecFromStr("6.0"),
		},
	}
)

type PerpetualKeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (k *PerpetualKeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain)

	k.legacyAmino = app.LegacyAmino()
	k.ctx = app.BaseApp.NewContext(initChain, tmproto.Header{})
	k.app = app
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(PerpetualKeeperTestSuite))
}

func (suite *PerpetualKeeperTestSuite) ResetSuite() {
	suite.SetupTest()
}

func (suite *PerpetualKeeperTestSuite) SetCurrentHeight(h int64) {
	suite.ctx = suite.ctx.WithBlockHeight(h)
}

func (suite *PerpetualKeeperTestSuite) AddBlockTime(d time.Duration) {
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(d))
}

func (suite *PerpetualKeeperTestSuite) SetupCoinPrices(ctx sdk.Context) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range priceMap {
		suite.app.OracleKeeper.SetAssetInfo(ctx, oracletypes.AssetInfo{
			Denom:   v.denom,
			Display: v.display,
			Decimal: 6,
		})
		suite.app.OracleKeeper.SetPrice(ctx, oracletypes.Price{
			Asset:     v.display,
			Price:     v.price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(ctx.BlockTime().Unix()),
		})
	}
}

func (suite *PerpetualKeeperTestSuite) AddCoinPrices(ctx sdk.Context, denoms []string) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, v := range denoms {
		suite.app.OracleKeeper.SetAssetInfo(ctx, oracletypes.AssetInfo{
			Denom:   priceMap[v].denom,
			Display: priceMap[v].display,
			Decimal: 6,
		})
		suite.app.OracleKeeper.SetPrice(ctx, oracletypes.Price{
			Asset:     priceMap[v].display,
			Price:     priceMap[v].price,
			Source:    "elys",
			Provider:  provider.String(),
			Timestamp: uint64(ctx.BlockTime().Unix()),
		})
	}
}

func (suite *PerpetualKeeperTestSuite) RemovePrices(ctx sdk.Context, denoms []string) {
	for _, v := range denoms {
		suite.app.OracleKeeper.RemoveAssetInfo(ctx, v)
		suite.app.OracleKeeper.RemovePrice(ctx, priceMap[v].display, "elys", uint64(ctx.BlockTime().Unix()))
	}
}

func TestSetGetMTP(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	perpetual := app.PerpetualKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	for i := 0; i < 2; i++ {
		mtp := types.MTP{
			Address:                        addr[i].String(),
			CollateralAsset:                ptypes.BaseCurrency,
			CustodyAsset:                   "ATOM",
			Collateral:                     sdk.NewInt(0),
			Liabilities:                    sdk.NewInt(0),
			BorrowInterestPaidCollateral:   sdk.NewInt(0),
			BorrowInterestPaidCustody:      sdk.NewInt(0),
			BorrowInterestUnpaidCollateral: sdk.NewInt(0),
			Custody:                        sdk.NewInt(0),
			Leverage:                       sdk.NewDec(0),
			MtpHealth:                      sdk.NewDec(0),
			Position:                       types.Position_LONG,
			Id:                             0,
			ConsolidateLeverage:            sdk.ZeroDec(),
			SumCollateral:                  sdk.ZeroInt(),
		}
		err := perpetual.SetMTP(ctx, &mtp)
		require.NoError(t, err)
	}

	mtpCount := perpetual.GetMTPCount(ctx)
	require.Equal(t, mtpCount, (uint64)(2))
}

func TestGetAllWhitelistedAddress(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	perpetual := app.PerpetualKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	// Set whitelisted addresses
	perpetual.WhitelistAddress(ctx, addr[0])
	perpetual.WhitelistAddress(ctx, addr[1])

	// Get all whitelisted addresses
	whitelisted := perpetual.GetAllWhitelistedAddress(ctx)

	// length should be 2
	require.Equal(t, len(whitelisted), 2)

	// If addr[0] is whitelisted
	require.Contains(t,
		whitelisted,
		addr[0],
	)

	// If addr[1] is whitelisted
	require.Contains(t,
		whitelisted,
		addr[1],
	)
}

func SetupStableCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper) {
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

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     sdk.NewDec(100),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     sdk.NewDec(100),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
}
