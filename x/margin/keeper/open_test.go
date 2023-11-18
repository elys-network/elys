package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestOpen_ErrorCheckUserAuthorization(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
			Position:        types.Position_LONG,
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(sdkerrors.Wrap(types.ErrUnauthorised, "unauthorised"))

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrUnauthorised))
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorCheckMaxOpenPositions(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         "creator",
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
			Position:        types.Position_LONG,
			Leverage:        sdk.NewDec(10),
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(sdkerrors.Wrap(types.ErrMaxOpenPositions, "cannot open new positions"))

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrMaxOpenPositions))
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorPreparePools(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         "creator",
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
			Position:        types.Position_LONG,
			Leverage:        sdk.NewDec(10),
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.BorrowAsset).Return(uint64(0), ammtypes.Pool{}, types.Pool{}, errors.New("error executing prepare pools"))

	_, err := k.Open(ctx, msg)

	assert.Equal(t, errors.New("error executing prepare pools"), err)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorCheckPoolHealth(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
			Position:        types.Position_LONG,
		}
		poolId = uint64(1)
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.BorrowAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset"))

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrInvalidBorrowingAsset))
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorInvalidPosition(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
		}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)

	_, err := k.Open(ctx, msg)

	assert.True(t, errors.Is(err, types.ErrInvalidPosition))
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorOpenLong(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
			Position:        types.Position_LONG,
		}
		poolId = uint64(1)
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.BorrowAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(nil)
	mockChecker.On("OpenLong", ctx, poolId, msg, ptypes.BaseCurrency).Return(&types.MTP{}, errors.New("error executing open long"))

	_, err := k.Open(ctx, msg)

	assert.Equal(t, errors.New("error executing open long"), err)
	mockChecker.AssertExpectations(t)
}

func TestOpen_ErrorOpenShort(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
			Position:        types.Position_SHORT,
		}
		poolId = uint64(1)
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.BorrowAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(nil)
	mockChecker.On("OpenShort", ctx, poolId, msg, ptypes.BaseCurrency).Return(&types.MTP{}, errors.New("error executing open short"))

	_, err := k.Open(ctx, msg)

	assert.Equal(t, errors.New("error executing open short"), err)
	mockChecker.AssertExpectations(t)
}

func TestOpen_Successful(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenChecker)
	mockAssetProfile := new(mocks.AssetProfileKeeper)

	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", nil, nil, nil, mockAssetProfile)
	k.OpenChecker = mockChecker

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			CollateralAsset: "aaa",
			BorrowAsset:     "bbb",
			Position:        types.Position_SHORT,
		}
		poolId = uint64(1)
		mtp    = &types.MTP{}
	)

	// Mock behavior
	mockAssetProfile.On("GetEntry", ctx, ptypes.BaseCurrency).Return(assetprofiletypes.Entry{BaseDenom: ptypes.BaseCurrency, Denom: ptypes.BaseCurrency}, true)
	mockChecker.On("CheckUserAuthorization", ctx, msg).Return(nil)
	mockChecker.On("CheckSameAssetPosition", ctx, msg).Return(nil)
	mockChecker.On("CheckMaxOpenPositions", ctx).Return(nil)
	mockChecker.On("PreparePools", ctx, msg.BorrowAsset).Return(poolId, ammtypes.Pool{}, types.Pool{}, nil)
	mockChecker.On("CheckPoolHealth", ctx, poolId).Return(nil)
	mockChecker.On("OpenShort", ctx, poolId, msg, ptypes.BaseCurrency).Return(mtp, nil)
	mockChecker.On("EmitOpenEvent", ctx, mtp).Return()

	_, err := k.Open(ctx, msg)

	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
