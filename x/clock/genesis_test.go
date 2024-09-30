package clock_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/app"
	clock "github.com/elys-network/elys/x/clock"
	"github.com/elys-network/elys/x/clock/types"
)

type GenesisTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app *app.ElysApp
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.app = app.InitElysTestApp(true)
	suite.ctx = suite.app.BaseApp.NewContext(false)
}

func (suite *GenesisTestSuite) TestClockInitGenesis() {
	_, _, addr := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()

	defaultParams := types.DefaultParams()

	testCases := []struct {
		name     string
		genesis  types.GenesisState
		expPanic bool
	}{
		{
			"default genesis",
			*clock.DefaultGenesisState(),
			false,
		},
		{
			"custom genesis - none",
			types.GenesisState{
				Params: types.Params{
					ContractAddresses: []string(nil),
					ContractGasLimit:  defaultParams.ContractGasLimit,
				},
			},
			false,
		},
		{
			"custom genesis - incorrect addr",
			types.GenesisState{
				Params: types.Params{
					ContractAddresses: []string{"incorrectaddr"},
					ContractGasLimit:  defaultParams.ContractGasLimit,
				},
			},
			true,
		},
		{
			"custom genesis - only one addr allowed",
			types.GenesisState{
				Params: types.Params{
					ContractAddresses: []string{addr.String(), addr2.String()},
					ContractGasLimit:  defaultParams.ContractGasLimit,
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			if tc.expPanic {
				suite.Require().Panics(func() {
					clock.InitGenesis(suite.ctx, suite.app.ClockKeeper, tc.genesis)
				})
			} else {
				suite.Require().NotPanics(func() {
					clock.InitGenesis(suite.ctx, suite.app.ClockKeeper, tc.genesis)
				})

				params := suite.app.ClockKeeper.GetParams(suite.ctx)
				suite.Require().Equal(tc.genesis.Params, params)
			}
		})
	}
}
