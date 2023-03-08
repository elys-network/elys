package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func TestGenesisInflationMsgServerUpdate(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateGenesisInflation
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateGenesisInflation{Authority: authority},
		},
		{
			desc:    "InvalidSigner",
			request: &types.MsgUpdateGenesisInflation{Authority: "B"},
			err:     govtypes.ErrInvalidSigner,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgUpdateGenesisInflation{Authority: authority}
			_, err := srv.UpdateGenesisInflation(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateGenesisInflation(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetGenesisInflation(ctx)
				require.True(t, found)
				require.Equal(t, expected.Authority, rst.Authority)
			}
		})
	}
}
