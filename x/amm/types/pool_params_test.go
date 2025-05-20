package types_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/v4/x/amm/types"
)

func (suite *TestSuite) TestPoolParamsValidate() {
	suite.SetupTest()
	poolParams := types.PoolParams{
		SwapFee:   sdkmath.LegacyDec{},
		UseOracle: false,
		FeeDenom:  "",
	}
	for _, tc := range []struct {
		name     string
		errorMsg string
		function func()
	}{
		{
			name:     "swap fee is nil",
			errorMsg: "swap_fee is nil",
			function: func() {
			},
		},
		{
			name:     "swap fee is -ve",
			errorMsg: types.ErrNegativeSwapFee.Error(),
			function: func() {
				poolParams.SwapFee = sdkmath.LegacyOneDec().MulInt64(-1)
			},
		},
		{
			name:     "swap fee is 1",
			errorMsg: types.ErrTooMuchSwapFee.Error(),
			function: func() {
				poolParams.SwapFee = sdkmath.LegacyOneDec()
			},
		},
		{
			name:     "invalid fee denom",
			errorMsg: "invalid denom: %%%",
			function: func() {
				poolParams.SwapFee = sdkmath.LegacyMustNewDecFromStr("0.001")
				poolParams.FeeDenom = "%%%"
			},
		},
		{
			name:     "success",
			errorMsg: "",
			function: func() {
				poolParams.FeeDenom = "uatom"
			},
		},
	} {
		suite.Run(tc.name, func() {
			tc.function()
			err := poolParams.Validate()
			if err != nil {
				suite.Require().Equal(tc.errorMsg, err.Error())
			} else {
				suite.Require().Equal(tc.errorMsg, "")
			}
		})
	}
}
