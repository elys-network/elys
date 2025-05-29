package types_test

import (
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/v6/x/commitment/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCommitClaimedRewards_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgCommitClaimedRewards
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgCommitClaimedRewards{
				Creator: "invalid_address",
				Amount:  sdkmath.ZeroInt(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgCommitClaimedRewards{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.OneInt(),
				Denom:   ptypes.ATOM,
			},
		},
		{
			name: "invalid denom",
			msg: types.MsgCommitClaimedRewards{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.OneInt(),
				Denom:   "@@@@@@",
			},
			err: errors.New("invalid denom"),
		},
		{
			name: "invalid amount - negative",
			msg:  types.MsgCommitClaimedRewards{sample.AccAddress(), sdkmath.NewInt(-200), ptypes.Eden},
			err:  types.ErrInvalidAmount,
		},
		{
			name: "invalid amount - nil",
			msg:  types.MsgCommitClaimedRewards{sample.AccAddress(), sdkmath.Int{}, ptypes.Eden},
			err:  types.ErrInvalidAmount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
