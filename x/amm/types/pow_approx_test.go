package types_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
	"math"
	"math/rand"
	"testing"
)

func TestPowApproximation(t *testing.T) {
	precision := sdk.MustNewDecFromStr("0.0001")
	baseInt := 100
	count := 0
	iterations := 100000
	for i := 0; i < iterations; i++ {
		numberF := float64(rand.Int63n(int64(baseInt))) + rand.Float64()
		expF := rand.Float64()
		number := sdk.MustNewDecFromStr(fmt.Sprintf("%f", numberF))
		exp := sdk.MustNewDecFromStr(fmt.Sprintf("%f", expF))
		value, _ := types.PowerApproximation(number, exp)
		libraryValue := sdk.MustNewDecFromStr(fmt.Sprintf("%f", math.Pow(numberF, expF)))
		err := libraryValue.Sub(value).Abs()
		assert := err.LTE(precision)
		if !assert {
			count++
		}
	}
	highInaccuracy := count < (iterations / 10)
	require.True(t, highInaccuracy)

}
