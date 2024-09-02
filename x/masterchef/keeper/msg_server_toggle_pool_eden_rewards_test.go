package keeper_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/assert"
)

func TestTogglePoolEdenRewards(t *testing.T) {

	type TestTogglePoolEdenReward struct {
		msg  types.MsgTogglePoolEdenRewards
		want error
	}

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	table := []TestTogglePoolEdenReward{
		{
			msg: types.MsgTogglePoolEdenRewards{
				Authority: "",
				PoolId:    2,
			},
			want: govtypes.ErrInvalidSigner,
		},
		{
			msg: types.MsgTogglePoolEdenRewards{
				Authority: authority,
				PoolId:    3,
			},
			want: types.ErrPoolNotFound,
		},
		{
			msg: types.MsgTogglePoolEdenRewards{
				Authority: authority,
				PoolId:    2,
				Enable:    true,
			},
			want: nil,
		},
	}
	ms, ctx := setupMsgServer(t)

	for _, test := range table {

		_, err := ms.TogglePoolEdenRewards(ctx, &test.msg)

		if err != nil {
			assert.Error(t, err, test.want)
		} else {
			assert.Equal(t, err, test.want)
		}

	}

}
