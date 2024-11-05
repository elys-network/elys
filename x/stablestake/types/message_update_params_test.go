package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateParams
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateParams{
				Authority: "invalid_address",
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdk.NewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdk.NewDecWithPrec(3, 2), // 0.03
					InterestRateMax:      sdk.NewDecWithPrec(5, 2), // 0.05
					InterestRateMin:      sdk.NewDecWithPrec(1, 2), // 0.01
					InterestRateIncrease: sdk.NewDecWithPrec(2, 2), // 0.02
					InterestRateDecrease: sdk.NewDecWithPrec(1, 2), // 0.01
					HealthGainFactor:     sdk.NewDecWithPrec(1, 1), // 0.1
					TotalValue:           sdk.NewInt(1000),
					MaxLeverageRatio:     sdk.MustNewDecFromStr("0.7"),
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdk.NewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdk.NewDecWithPrec(3, 2), // 0.03
					InterestRateMax:      sdk.NewDecWithPrec(5, 2), // 0.05
					InterestRateMin:      sdk.NewDecWithPrec(1, 2), // 0.01
					InterestRateIncrease: sdk.NewDecWithPrec(2, 2), // 0.02
					InterestRateDecrease: sdk.NewDecWithPrec(1, 2), // 0.01
					HealthGainFactor:     sdk.NewDecWithPrec(1, 1), // 0.1
					TotalValue:           sdk.NewInt(1000),
					MaxLeverageRatio:     sdk.MustNewDecFromStr("0.7"),
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

func TestNewMsgUpdateParams(t *testing.T) {

	accAdress := sample.AccAddress()

	params := DefaultParams()
	got := NewMsgUpdateParams(
		accAdress,
		&params,
	)

	want := &MsgUpdateParams{
		Authority: accAdress,
		Params:    &params,
	}

	assert.Equal(t, want, got)
}

func TestMsgUpdateParams_Route(t *testing.T) {
	msg := MsgUpdateParams{}
	assert.Equal(t, "stablestake", msg.Route())
}

func TestMsgUpdateParams_Type(t *testing.T) {
	msg := MsgUpdateParams{}
	assert.Equal(t, "update_params", msg.Type())
}

func TestMsgUpdateParams_GetSigners(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgUpdateParams{Authority: accAdress}

	creator, err := sdk.AccAddressFromBech32(accAdress)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, []sdk.AccAddress{creator}, msg.GetSigners())
}

func TestMsgUpdateParams_GetSignBytes(t *testing.T) {
	accAdress := sample.AccAddress()
	msg := MsgUpdateParams{Authority: accAdress}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	b := sdk.MustSortJSON(bz)

	assert.Equal(t, b, msg.GetSignBytes())
}
