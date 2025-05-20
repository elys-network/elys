package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/amm/types"
)

func (suite *AmmKeeperTestSuite) TestOraclePoolSlippageTrack() {
	tracks := []types.OraclePoolSlippageTrack{
		{
			PoolId:    1,
			Timestamp: 100000,
			Tracked:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000)},
		},
		{
			PoolId:    1,
			Timestamp: 100001,
			Tracked:   sdk.Coins{sdk.NewInt64Coin("uelys", 2000)},
		},
		{
			PoolId:    2,
			Timestamp: 100001,
			Tracked:   sdk.Coins{sdk.NewInt64Coin("uelys", 1000)},
		},
	}
	for _, track := range tracks {
		suite.app.AmmKeeper.SetSlippageTrack(suite.ctx, track)
	}
	for _, track := range tracks {
		t := suite.app.AmmKeeper.GetSlippageTrack(suite.ctx, track.PoolId, track.Timestamp)
		suite.Require().Equal(t, track)
	}
	tracksStored := suite.app.AmmKeeper.AllSlippageTracks(suite.ctx)
	suite.Require().Equal(tracks, tracksStored)

	track := suite.app.AmmKeeper.GetLastSlippageTrack(suite.ctx, 1)
	suite.Require().Equal(track, tracks[1])
	track = suite.app.AmmKeeper.GetFirstSlippageTrack(suite.ctx, 1)
	suite.Require().Equal(track, tracks[0])

	diff := suite.app.AmmKeeper.GetTrackedSlippageDiff(suite.ctx, 1)
	suite.Require().Equal(diff.String(), tracks[1].Tracked.Sub(tracks[0].Tracked...).String())

	suite.app.AmmKeeper.DeleteSlippageTrack(suite.ctx, tracks[0])
	tracksStored = suite.app.AmmKeeper.AllSlippageTracks(suite.ctx)
	suite.Require().Len(tracksStored, 2)

	suite.app.AmmKeeper.TrackSlippage(suite.ctx, 3, sdk.NewInt64Coin("uelys", 1))
	track = suite.app.AmmKeeper.GetSlippageTrack(suite.ctx, 3, uint64(suite.ctx.BlockTime().Unix()))
	suite.Require().Equal(track.Tracked.String(), "1uelys")
}
