package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNotFound_GetNetOpenInterest(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenDefineAssetsChecker = mockChecker

	ctx := sdk.Context{}
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, false)

	got := k.GetNetOpenInterest(ctx, types.Pool{})

	assert.Equal(t, got, sdk.ZeroInt())
}

func TestSuccess_GetNetOpenInterest(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)
	mockAmm := new(mocks.AmmKeeper)

	ctx := sdk.Context{}

	uusdcDenom := "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65"

	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: uusdcDenom}, true)

	pool := ammtypes.Pool{
		PoolId: 2,
	}

	mockAmm.On("GetPool", ctx, pool.PoolId).Return(pool, true)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, mockAssetProfile, nil)

	coin := sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", sdk.NewInt(2000))
	mockChecker.On("EstimateSwapGivenOut", ctx, coin, uusdcDenom, pool).Return(sdk.NewInt(1000), nil)

	k.OpenDefineAssetsChecker = mockChecker

	got := k.GetNetOpenInterest(ctx, types.Pool{
		AmmPoolId: 2,
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.NewInt(2000),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.ZeroInt(),
			},
		},
		PoolAssetsShort: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.ZeroInt(),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.NewInt(2000),
			},
		},
	})

	assert.Equal(t, got, sdk.NewInt(1000))
}

func TestSuccess2_GetNetOpenInterest(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)
	mockAmm := new(mocks.AmmKeeper)

	ctx := sdk.Context{}

	uusdcDenom := "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65"

	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: uusdcDenom}, true)

	pool := ammtypes.Pool{
		PoolId: 2,
	}

	mockAmm.On("GetPool", ctx, pool.PoolId).Return(pool, true)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, mockAssetProfile, nil)

	coin := sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", sdk.NewInt(2000))
	mockChecker.On("EstimateSwapGivenOut", ctx, coin, uusdcDenom, pool).Return(sdk.NewInt(1000), nil)

	k.OpenDefineAssetsChecker = mockChecker

	got := k.GetNetOpenInterest(ctx, types.Pool{
		AmmPoolId: 2,
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.ZeroInt(),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.NewInt(2000),
			},
		},
		PoolAssetsShort: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.NewInt(500),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.ZeroInt(),
			},
		},
	})

	assert.Equal(t, got, sdk.NewInt(500))
}

func TestErrorLongPoolNotFound_GetNetOpenInterest(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)
	mockAmm := new(mocks.AmmKeeper)

	ctx := sdk.Context{}

	uusdcDenom := "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65"

	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: uusdcDenom}, true)

	pool := ammtypes.Pool{
		PoolId: 2,
	}

	mockAmm.On("GetPool", ctx, pool.PoolId).Return(pool, false)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, mockAssetProfile, nil)

	coin := sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", sdk.NewInt(2000))
	mockChecker.On("EstimateSwapGivenOut", ctx, coin, uusdcDenom, pool).Return(sdk.NewInt(1000), nil)

	k.OpenDefineAssetsChecker = mockChecker

	got := k.GetNetOpenInterest(ctx, types.Pool{
		AmmPoolId: 2,
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.ZeroInt(),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.NewInt(2000),
			},
		},
		PoolAssetsShort: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.NewInt(500),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.ZeroInt(),
			},
		},
	})

	assert.Equal(t, got, sdk.ZeroInt())
}

func TestErrorLongEstimateSwapGivenOut_GetNetOpenInterest(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)
	mockAmm := new(mocks.AmmKeeper)

	ctx := sdk.Context{}

	uusdcDenom := "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65"

	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: uusdcDenom}, true)

	pool := ammtypes.Pool{
		PoolId: 2,
	}

	mockAmm.On("GetPool", ctx, pool.PoolId).Return(pool, true)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, mockAssetProfile, nil)

	coin := sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", sdk.NewInt(2000))
	mockChecker.On("EstimateSwapGivenOut", ctx, coin, uusdcDenom, pool).Return(sdk.NewInt(1000), types.ErrInvalidAmount)

	k.OpenDefineAssetsChecker = mockChecker

	got := k.GetNetOpenInterest(ctx, types.Pool{
		AmmPoolId: 2,
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.ZeroInt(),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.NewInt(2000),
			},
		},
		PoolAssetsShort: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.NewInt(500),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.ZeroInt(),
			},
		},
	})

	assert.Equal(t, got, sdk.ZeroInt())
}

func TestErrorShortPoolNotFound_GetNetOpenInterest(t *testing.T) {
	mockChecker := new(mocks.OpenDefineAssetsChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)
	mockAmm := new(mocks.AmmKeeper)

	ctx := sdk.Context{}

	uusdcDenom := "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65"

	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: uusdcDenom}, true)

	pool := ammtypes.Pool{
		PoolId: 2,
	}

	mockAmm.On("GetPool", ctx, pool.PoolId).Return(pool, false)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, mockAssetProfile, nil)

	coin := sdk.NewCoin("ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953", sdk.NewInt(2000))
	mockChecker.On("EstimateSwapGivenOut", ctx, coin, uusdcDenom, pool).Return(sdk.NewInt(1000), nil)

	k.OpenDefineAssetsChecker = mockChecker

	got := k.GetNetOpenInterest(ctx, types.Pool{
		AmmPoolId: 2,
		PoolAssetsLong: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.NewInt(2000),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.ZeroInt(),
			},
		},
		PoolAssetsShort: []types.PoolAsset{
			{
				AssetDenom:  "ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65",
				Liabilities: sdk.ZeroInt(),
			},
			{
				AssetDenom:  "ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
				Liabilities: sdk.NewInt(2000),
			},
		},
	})

	assert.Equal(t, got, sdk.ZeroInt())
}
