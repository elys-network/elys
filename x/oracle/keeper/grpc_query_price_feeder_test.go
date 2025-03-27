package keeper_test

import (
	"strconv"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/oracle/types"
)

func (suite *KeeperTestSuite) TestPriceFeederQuerySingle() {
	keeper, ctx := suite.app.OracleKeeper, suite.ctx

	msgs := createNPriceFeeder(&keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetPriceFeederRequest
		response *types.QueryGetPriceFeederResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPriceFeederRequest{
				Feeder: msgs[0].Feeder,
			},
			response: &types.QueryGetPriceFeederResponse{PriceFeeder: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPriceFeederRequest{
				Feeder: msgs[1].Feeder,
			},
			response: &types.QueryGetPriceFeederResponse{PriceFeeder: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPriceFeederRequest{
				Feeder: authtypes.NewModuleAddress(strconv.Itoa(100000)).String(),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		suite.Run(tc.desc, func() {
			response, err := keeper.PriceFeeder(ctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPriceFeederQueryPaginated() {
	keeper, ctx := suite.app.OracleKeeper, suite.ctx

	msgs := createNPriceFeeder(&keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPriceFeederRequest {
		return &types.QueryAllPriceFeederRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	suite.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PriceFeederAll(ctx, request(nil, uint64(i), uint64(step), false))
			suite.Require().NoError(err)
			suite.Require().LessOrEqual(len(resp.PriceFeeder), step)
			suite.Require().Subset(
				nullify.Fill(msgs),
				nullify.Fill(resp.PriceFeeder),
			)
		}
	})
	suite.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PriceFeederAll(ctx, request(next, 0, uint64(step), false))
			suite.Require().NoError(err)
			suite.Require().LessOrEqual(len(resp.PriceFeeder), step)
			suite.Require().Subset(
				nullify.Fill(msgs),
				nullify.Fill(resp.PriceFeeder),
			)
			next = resp.Pagination.NextKey
		}
	})
	suite.Run("Total", func() {
		resp, err := keeper.PriceFeederAll(ctx, request(nil, 0, 0, true))
		suite.Require().NoError(err)
		suite.Require().Equal(len(msgs), int(resp.Pagination.Total))
		suite.Require().ElementsMatch(
			nullify.Fill(msgs),
			nullify.Fill(resp.PriceFeeder),
		)
	})
	suite.Run("InvalidRequest", func() {
		_, err := keeper.PriceFeederAll(ctx, nil)
		suite.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
