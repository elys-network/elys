package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/testutil/nullify"
	"github.com/elys-network/elys/v6/x/tokenomics/types"
)

func TestGenesisInflationQuery(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)

	item := createTestGenesisInflation(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetGenesisInflationRequest
		response *types.QueryGetGenesisInflationResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetGenesisInflationRequest{},
			response: &types.QueryGetGenesisInflationResponse{GenesisInflation: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GenesisInflation(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
