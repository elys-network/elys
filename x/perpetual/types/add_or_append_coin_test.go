package types

import (
	"cosmossdk.io/math"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAddOrAppendCoin(t *testing.T) {
	tests := []struct {
		name     string
		coins    sdk.Coins
		newCoin  sdk.Coin
		expected sdk.Coins
	}{
		{
			name:     "Append new coin",
			coins:    sdk.NewCoins(sdk.NewCoin("atom", math.NewInt(100))),
			newCoin:  sdk.NewCoin("eth", math.NewInt(50)),
			expected: sdk.NewCoins(sdk.NewCoin("atom", math.NewInt(100)), sdk.NewCoin("eth", math.NewInt(50))),
		},
		{
			name:     "Aggregate coin amount",
			coins:    sdk.NewCoins(sdk.NewCoin("atom", math.NewInt(100))),
			newCoin:  sdk.NewCoin("atom", math.NewInt(50)),
			expected: sdk.NewCoins(sdk.NewCoin("atom", math.NewInt(150))),
		},
		{
			name:     "Aggregate coin amount in larger slice",
			coins:    sdk.NewCoins(sdk.NewCoin("atom", math.NewInt(100)), sdk.NewCoin("btc", math.NewInt(200)), sdk.NewCoin("eth", math.NewInt(300))),
			newCoin:  sdk.NewCoin("btc", math.NewInt(50)),
			expected: sdk.NewCoins(sdk.NewCoin("atom", math.NewInt(100)), sdk.NewCoin("btc", math.NewInt(250)), sdk.NewCoin("eth", math.NewInt(300))),
		},
		{
			name:     "Append new coin to empty slice",
			coins:    sdk.NewCoins(),
			newCoin:  sdk.NewCoin("eth", math.NewInt(50)),
			expected: sdk.NewCoins(sdk.NewCoin("eth", math.NewInt(50))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddOrAppendCoin(tt.coins, tt.newCoin)
			if !reflect.DeepEqual(sdk.NewCoins(got...), tt.expected) {
				t.Errorf("AddOrAppendCoin() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
