package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	initChain = true
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	app         *simapp.ElysApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.InitElysTestApp(initChain)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(initChain, tmproto.Header{})
	suite.app = app
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// TestKeeper_Logger tests the Logger function
func TestKeeper_Logger(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	// Create a test context and keeper
	keeper := app.StablestakeKeeper

	logger := app.Logger()

	keeper.Logger(ctx).Info("test")
	logger.Info("test")
}

// TestKeeper_SetHooks_Panic tests the SetHooks function with a nil argument
func TestKeeper_SetHooks_Panic(t *testing.T) {
	app := simapp.InitElysTestApp(true)

	keeper := app.StablestakeKeeper

	assert.Panics(t, func() {
		keeper.SetHooks(nil)
	})
}
