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

// The test cases here are comparing computed values with go's in built function over float64.
// How to make sure those are accurate, and they only generate upto 6 decimal places.
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
	lowInaccuracy := count < ((iterations * 2) / 100) // 2% inaccuracy
	require.True(t, lowInaccuracy)

	// 2^255 + OneDec - SmallestDec, Max number lies somewhere between 2^256 and 2^255
	edgeCaseNumber := sdk.MustNewDecFromStr("57896044618658097711785492504343953926634992332820282019728792003956564819968.999999999999999999")
	edgeCaseExponent := sdk.MustNewDecFromStr("0.999999999999999999")
	_, err := types.PowerApproximation(edgeCaseNumber, edgeCaseExponent)
	require.NoError(t, err)

}
