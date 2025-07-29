package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateTimeBasedInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateTimeBasedInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateTimeBasedInflation{
				Authority:        "invalid_address",
				StartBlockHeight: 100,
				EndBlockHeight:   200,
				Description:      "Valid description",
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateTimeBasedInflation{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 100,
				EndBlockHeight:   200,
				Description:      "Valid description",
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
			},
		},
		{
			name: "end block height before start block height",
			msg: MsgCreateTimeBasedInflation{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 200,
				EndBlockHeight:   100,
				Description:      "Valid description",
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty description",
			msg: MsgCreateTimeBasedInflation{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 100,
				EndBlockHeight:   200,
				Description:      "",
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "nil inflation entry",
			msg: MsgCreateTimeBasedInflation{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 100,
				EndBlockHeight:   200,
				Description:      "Valid description",
				Inflation:        nil,
			},
			err: sdkerrors.ErrInvalidRequest,
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

func TestMsgUpdateTimeBasedInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateTimeBasedInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateTimeBasedInflation{
				Authority:        "invalid_address",
				StartBlockHeight: 100,
				EndBlockHeight:   200,
				Description:      "Valid description",
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateTimeBasedInflation{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 100,
				EndBlockHeight:   200,
				Description:      "Valid description",
				Inflation: &InflationEntry{
					LmRewards:         1000,
					IcsStakingRewards: 500,
					CommunityFund:     300,
					StrategicReserve:  200,
					TeamTokensVested:  100,
				},
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

func TestMsgDeleteTimeBasedInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteTimeBasedInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteTimeBasedInflation{
				Authority:        "invalid_address",
				StartBlockHeight: 100,
				EndBlockHeight:   200,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteTimeBasedInflation{
				Authority:        sample.AccAddress(),
				StartBlockHeight: 100,
				EndBlockHeight:   200,
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
