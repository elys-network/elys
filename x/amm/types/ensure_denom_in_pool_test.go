package types_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
)

func TestEnsureDenomInPool(t *testing.T) {
	poolAssetsByDenom := map[string]types.PoolAsset{
		"abc": {Token: sdk.NewCoin("abc", math.ZeroInt())},
		"def": {Token: sdk.NewCoin("def", math.ZeroInt())},
		"ghi": {Token: sdk.NewCoin("ghi", math.ZeroInt())},
	}

	tests := []struct {
		tokensIn sdk.Coins
		err      error
	}{
		{sdk.NewCoins(sdk.NewInt64Coin("abc", 100), sdk.NewInt64Coin("def", 200)), nil},
		{sdk.NewCoins(sdk.NewInt64Coin("def", 200), sdk.NewInt64Coin("jkl", 300)), sdkerrors.Wrapf(types.ErrDenomNotFoundInPool, types.InvalidInputDenomsErrFormat, "jkl")},

		{sdk.NewCoins(sdk.NewInt64Coin("xyz", 500)), sdkerrors.Wrapf(types.ErrDenomNotFoundInPool, types.InvalidInputDenomsErrFormat, "xyz")},
	}

	for _, tt := range tests {
		err := types.EnsureDenomInPool(poolAssetsByDenom, tt.tokensIn)
		if !errors.Is(err, tt.err) {
			t.Errorf("EnsureDenomInPool(%v) = %v; want %v", tt.tokensIn, err, tt.err)
		}
	}
}
