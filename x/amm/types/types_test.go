package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/elys-network/elys/app"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

const (
	initChain = true
)

type TestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (suite *TestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain, tmproto.Header{})
	suite.app = app
}

func (suite *TestSuite) SetupStableCoinPrices() {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusdc",
		Display: "USDC",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusdt",
		Display: "USDT",
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
	suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdk.NewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(suite.ctx.BlockTime().Unix()),
	})
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
