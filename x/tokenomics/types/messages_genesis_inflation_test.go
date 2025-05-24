package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateGenesisInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateGenesisInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateGenesisInflation{
				Authority: "invalid_address",
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
				SeedVesting:           100,
				StrategicSalesVesting: 200,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateGenesisInflation{
				Authority: sample.AccAddress(),
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
				SeedVesting:           100,
				StrategicSalesVesting: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
