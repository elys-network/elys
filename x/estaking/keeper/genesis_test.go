package keeper_test

import (
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/estaking/types"
)

func (suite *EstakingKeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	suite.app.EstakingKeeper.InitGenesis(suite.ctx, genesisState)
	got := suite.app.EstakingKeeper.ExportGenesis(suite.ctx)
	suite.Require().NotNil(got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
