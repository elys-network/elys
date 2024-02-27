package keeper_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestOpen_ErrorCheckUserAuthorization(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Position:     types.Position_LONG,
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(errorsmod.Wrap(types.ErrUnauthorised, "unauthorised"))

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrUnauthorised))
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorCheckMaxOpenPositions(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Position:     types.Position_LONG,
			Leverage:     sdk.NewDec(10),
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions"))

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrMaxOpenPositions))
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorPreparePools(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:      "creator",
			Position:     types.Position_LONG,
			Leverage:     sdk.NewDec(10),
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.Collateral.Denom, msg.TradingAsset).Return(uint64(0), ammtypes.Pool{}, types.Pool{}, errors.New("error executing prepare pools"))

	_, err := k.Open(ctx, msg)

	assert.Equal(t, errors.New("error executing prepare pools"), err)
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorCheckPoolHealth(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Position:     types.Position_LONG,
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
		poolId = uint64(1)
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.Collateral.Denom, msg.TradingAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset"))

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrInvalidBorrowingAsset))
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorInvalidPosition(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrInvalidPosition))
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorOpenLong(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Position:     types.Position_LONG,
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
		poolId = uint64(1)
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.Collateral.Denom, msg.TradingAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(nil)
	mockChecker.On("OpenLong", ctx, poolId, msg, ptypes.BaseCurrency).Return(&types.MTP{}, errors.New("error executing open long"))

	_, err := k.Open(ctx, msg)

	assert.Equal(t, errors.New("error executing open long"), err)
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorOpenShort(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Position:     types.Position_SHORT,
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
		poolId = uint64(1)
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.Collateral.Denom, msg.TradingAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(nil)
	mockChecker.On("OpenShort", ctx, poolId, msg, ptypes.BaseCurrency).Return(&types.MTP{}, errors.New("error executing open short"))

	_, err := k.Open(ctx, msg)

	assert.Equal(t, errors.New("error executing open short"), err)
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}

func TestOpen_Successful(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile, nil)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Position:     types.Position_SHORT,
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, sdk.OneInt()),
		}
		poolId = uint64(1)
		mtp    = &types.MTP{}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.Collateral.Denom, msg.TradingAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(nil)
	mockChecker.On("OpenShort", ctx, poolId, msg, ptypes.BaseCurrency).Return(mtp, nil)
	mockChecker.On("UpdateOpenPrice", ctx, mtp, ammtypes.Pool{}, ptypes.BaseCurrency).Return(nil)
	mockChecker.On("EmitOpenEvent", ctx, mtp).Return()

	_, err := k.Open(ctx, msg)

	assert.Nil(t, err)
	mockAssetProfile.AssertExpectations(t)
	mockChecker.AssertExpectations(t)
}
