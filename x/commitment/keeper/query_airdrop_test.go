// query_airdrop_test.go
package keeper_test

import (
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAirDrop(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(false)
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))
	app.CommitmentKeeper.SetAtomStaker(ctx, types.AtomStaker{
		Address: addr[0].String(),
		Amount:  sdkmath.NewInt(1000),
	})

	for _, tc := range []struct {
		desc     string
		request  *types.QueryAirDropRequest
		response *types.QueryAirDropResponse
		err      error
	}{
		{
			desc: "ValidRequest",
			request: &types.QueryAirDropRequest{
				Address: addr[0].String(),
			},
			response: &types.QueryAirDropResponse{
				AtomStaking: sdkmath.NewInt(1000),
				Cadet:       sdkmath.NewInt(0),
				NftHolder:   sdkmath.NewInt(0),
				Governor:    sdkmath.NewInt(0),
			},
		},
		{
			desc:    "InvalidRequest",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
		{
			desc: "InvalidAddress",
			request: &types.QueryAirDropRequest{
				Address: "invalidaddress",
			},
			err: errors.New("decoding bech32 failed"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := app.CommitmentKeeper.AirDrop(ctx, tc.request)
			if tc.err != nil {
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
