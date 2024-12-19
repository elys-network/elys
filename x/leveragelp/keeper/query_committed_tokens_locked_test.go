// query_committed_tokens_locked_test.go
package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	simapp "github.com/elys-network/elys/app"
	cmttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (suite *KeeperTestSuite) TestCommittedTokensLocked() {
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdkmath.NewInt(1000000))
	position := types.Position{
		Address: addr[0].String(),
	}
	suite.app.LeveragelpKeeper.SetPosition(suite.ctx, &position)
	commitments := cmttypes.Commitments{
		Creator:         addr[0].String(),
		CommittedTokens: []*cmttypes.CommittedTokens{{Denom: "uelys", Amount: sdkmath.NewInt(100)}},
	}
	suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitments)

	_, err := suite.app.LeveragelpKeeper.CommittedTokensLocked(suite.ctx, nil)
	suite.Require().Error(err)

	request := &types.QueryCommittedTokensLockedRequest{
		Address: "nonbech32",
	}
	_, err = suite.app.LeveragelpKeeper.CommittedTokensLocked(suite.ctx, request)
	suite.Require().Error(err)

	request = &types.QueryCommittedTokensLockedRequest{
		Address: position.Address,
	}
	_, err = suite.app.LeveragelpKeeper.CommittedTokensLocked(suite.ctx, request)
	suite.Require().NoError(err)
}
