package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func TestGenesisInflationQuery(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
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
			response, err := keeper.GenesisInflation(wctx, tc.request)
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
