package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSuccessPoolAssetsLong_CalcTotalLiabilities(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAmm := new(mocks.AmmKeeper)

	ctx := sdk.Context{}

	uusdcDenom := "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65"

	pool := ammtypes.Pool{
		PoolId: 2,
	}

	mockAmm.On("GetPool", ctx, pool.PoolId).Return(pool, true)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil, nil)

	coin := sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", sdk.NewInt(2000))
	mockChecker.On("EstimateSwapGivenOut", ctx, coin, uusdcDenom, pool).Return(sdk.NewInt(1000), nil)

	k.OpenDefineAssetsChecker = mockChecker

	ammPoolId := uint64(2)
	assets := []types.PoolAsset{
		{
			AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
			Liabilities: sdk.NewInt(2000),
		},
		{
			AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
			Liabilities: sdk.NewInt(0),
		},
	}

	got, err := k.CalcTotalLiabilities(ctx, assets, ammPoolId, uusdcDenom)

	assert.NoError(t, err)
	assert.Equal(t, got, sdk.NewInt(2000))
}

func TestSuccessPoolAssetsShort_CalcTotalLiabilities(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAmm := new(mocks.AmmKeeper)

	ctx := sdk.Context{}

	uusdcDenom := "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65"

	pool := ammtypes.Pool{
		PoolId: 2,
	}

	mockAmm.On("GetPool", ctx, pool.PoolId).Return(pool, true)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil, nil)

	coin := sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", sdk.NewInt(2000))
	mockChecker.On("EstimateSwapGivenOut", ctx, coin, uusdcDenom, pool).Return(sdk.NewInt(1000), nil)

	k.OpenDefineAssetsChecker = mockChecker

	ammPoolId := uint64(2)
	assets := []types.PoolAsset{
		{
			AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
			Liabilities: sdk.NewInt(0),
		},
		{
			AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
			Liabilities: sdk.NewInt(2000),
		},
	}

	got, err := k.CalcTotalLiabilities(ctx, assets, ammPoolId, uusdcDenom)

	assert.NoError(t, err)
	assert.Equal(t, got, sdk.NewInt(1000))
}
