package types_test

import (
	"fmt"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateVestingInfo_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateVestingInfo
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateVestingInfo{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgUpdateVestingInfo{
				Authority:      sample.AccAddress(),
				VestingDenom:   ptypes.ATOM,
				BaseDenom:      ptypes.BaseCurrency,
				NumBlocks:      10,
				VestNowFactor:  10,
				NumMaxVestings: 10,
			},
		},
		{
			name: "invalid vesting denom",
			msg: types.MsgUpdateVestingInfo{
				Authority:      sample.AccAddress(),
				VestingDenom:   "",
				BaseDenom:      ptypes.BaseCurrency,
				NumBlocks:      10,
				VestNowFactor:  10,
				NumMaxVestings: 10,
			},
			err: fmt.Errorf("invalid denom"),
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
