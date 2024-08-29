package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/parameter/keeper"
	"github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_UpdateRewardsDataLifetime(t *testing.T) {
	k, ctx := keepertest.ParameterKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)

	tests := []struct {
		name    string
		msg     *types.MsgUpdateRewardsDataLifetime
		want    *types.MsgUpdateRewardsDataLifetimeResponse
		wantErr bool
	}{
		{
			name: "Valid positive rewards data lifetime",
			msg: &types.MsgUpdateRewardsDataLifetime{
				Creator:             "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				RewardsDataLifetime: "100",
			},
			want:    &types.MsgUpdateRewardsDataLifetimeResponse{},
			wantErr: false,
		},
		{
			name: "Invalid zero rewards data lifetime",
			msg: &types.MsgUpdateRewardsDataLifetime{
				Creator:             "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				RewardsDataLifetime: "0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid negative rewards data lifetime",
			msg: &types.MsgUpdateRewardsDataLifetime{
				Creator:             "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				RewardsDataLifetime: "-10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid non-numeric rewards data lifetime",
			msg: &types.MsgUpdateRewardsDataLifetime{
				Creator:             "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				RewardsDataLifetime: "abc",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.UpdateRewardsDataLifetime(wctx, tt.msg)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
