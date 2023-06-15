package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/elys-network/elys/app"
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

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
