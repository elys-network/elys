package keeper_test

import (
	"testing"

	testkeeper "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/assetprofile/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.AssetprofileKeeper(t)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryParamsRequest
		response *types.QueryParamsResponse
		err      error
	}{
		{
			desc:    "ValidRequest",
			request: &types.QueryParamsRequest{},
			response: &types.QueryParamsResponse{
				Params: keeper.GetParams(ctx),
			},
		},
		{
			desc:    "InvalidRequest",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Params(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
