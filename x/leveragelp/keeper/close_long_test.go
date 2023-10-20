package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/elys-network/elys/x/leveragelp/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCloseLong_MtpNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "creator",
			Id:      1,
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(types.MTP{}, types.ErrMTPDoesNotExist)

	_, _, err := k.CloseLong(ctx, msg)

	// Expect an error about the mtp not existing
	assert.True(t, errors.Is(err, types.ErrMTPDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_PoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "creator",
			Id:      1,
		}
		mtp = types.MTP{
			AmmPoolId: 2,
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(types.Pool{}, false)

	_, _, err := k.CloseLong(ctx, msg)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrInvalidBorrowingAsset))
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_AmmPoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "creator",
			Id:      1,
		}
		mtp = types.MTP{
			AmmPoolId: 2,
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(types.Pool{}, true)

	_, _, err := k.CloseLong(ctx, msg)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_ErrorEstimateAndRepay(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "creator",
			Id:      1,
		}
		mtp = types.MTP{
			AmmPoolId: 2,
		}
		pool = types.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)

	_, _, err := k.CloseLong(ctx, msg)

	// Expect an error about estimate and repay
	assert.Equal(t, errors.New("error executing estimate and repay"), err)
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_SuccessfulClosingLongPosition(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "creator",
			Id:      1,
		}
		mtp = types.MTP{
			AmmPoolId: 2,
		}
		pool        = types.Pool{}
		repayAmount = math.NewInt(100)
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)

	mtpOut, repayAmountOut, err := k.CloseLong(ctx, msg)

	// Expect no error
	assert.Nil(t, err)
	assert.Equal(t, repayAmount, repayAmountOut)
	assert.Equal(t, mtp, *mtpOut)
	mockChecker.AssertExpectations(t)
}
