package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (suite *TestSuite) TestCalculateTokenARate() {
	tokenARate := types.CalculateTokenARate(sdk.NewDec(10), sdk.NewDec(3), sdk.NewDec(100), sdk.NewDec(6))
	suite.Require().Equal(tokenARate.String(), sdk.NewDec(5).String())

	tokenARate = types.CalculateTokenARate(sdk.NewDec(0), sdk.NewDec(3), sdk.NewDec(10), sdk.NewDec(6))
	suite.Require().Equal(tokenARate.String(), sdk.NewDec(0).String())

	tokenARate = types.CalculateTokenARate(sdk.NewDec(10), sdk.NewDec(3), sdk.NewDec(100), sdk.NewDec(0))
	suite.Require().Equal(tokenARate.String(), sdk.NewDec(0).String())
}
