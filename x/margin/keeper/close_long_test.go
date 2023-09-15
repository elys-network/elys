package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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
			AmmPoolId:     2,
			CustodyAssets: []string{"uatom"},
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(types.Pool{}, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAssets[0]).Return(ammtypes.Pool{}, sdkerrors.Wrap(types.ErrPoolDoesNotExist, mtp.CustodyAssets[0]))

	_, _, err := k.CloseLong(ctx, msg)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_ErrorHandleInterest(t *testing.T) {
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
			AmmPoolId:        2,
			CustodyAssets:    []string{"uatom"},
			CollateralAssets: []string{ptypes.BaseCurrency},
		}
		pool = types.Pool{
			InterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAssets[0]).Return(ammPool, nil)
	mockChecker.On("HandleInterest", ctx, &mtp, &pool, ammPool, mtp.CollateralAssets[0], mtp.CustodyAssets[0]).Return(errors.New("error executing handle interest"))

	_, _, err := k.CloseLong(ctx, msg)

	// Expect an error about handle interest
	assert.Equal(t, errors.New("error executing handle interest"), err)
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_ErrorTakeOutCustody(t *testing.T) {
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
			AmmPoolId:        2,
			CustodyAssets:    []string{"uatom"},
			CollateralAssets: []string{ptypes.BaseCurrency},
		}
		pool = types.Pool{
			InterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAssets[0]).Return(ammPool, nil)
	mockChecker.On("HandleInterest", ctx, &mtp, &pool, ammPool, mtp.CollateralAssets[0], mtp.CustodyAssets[0]).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, mtp.CustodyAssets[0]).Return(errors.New("error executing take out custody"))

	_, _, err := k.CloseLong(ctx, msg)

	// Expect an error about take out custody
	assert.Equal(t, errors.New("error executing take out custody"), err)
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
			AmmPoolId:        2,
			CustodyAssets:    []string{"uatom"},
			CollateralAssets: []string{ptypes.BaseCurrency},
		}
		pool = types.Pool{
			InterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAssets[0]).Return(ammPool, nil)
	mockChecker.On("HandleInterest", ctx, &mtp, &pool, ammPool, mtp.CollateralAssets[0], mtp.CustodyAssets[0]).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, mtp.CustodyAssets[0]).Return(nil)
	mockChecker.On("EstimateAndRepay", ctx, mtp, pool, ammPool, mtp.CollateralAssets[0], mtp.CustodyAssets[0]).Return(sdk.Int{}, errors.New("error executing estimate and repay"))

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
			AmmPoolId:        2,
			CustodyAssets:    []string{"uatom"},
			CollateralAssets: []string{ptypes.BaseCurrency},
		}
		pool = types.Pool{
			InterestRate: math.LegacyNewDec(2),
		}
		ammPool     = ammtypes.Pool{}
		repayAmount = math.NewInt(100)
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAssets[0]).Return(ammPool, nil)
	mockChecker.On("HandleInterest", ctx, &mtp, &pool, ammPool, mtp.CollateralAssets[0], mtp.CustodyAssets[0]).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, mtp.CustodyAssets[0]).Return(nil)
	mockChecker.On("EstimateAndRepay", ctx, mtp, pool, ammPool, mtp.CollateralAssets[0], mtp.CustodyAssets[0]).Return(repayAmount, nil)

	mtpOut, repayAmountOut, err := k.CloseLong(ctx, msg)

	// Expect no error
	assert.Nil(t, err)
	assert.Equal(t, repayAmount, repayAmountOut)
	assert.Equal(t, mtp, *mtpOut)
	mockChecker.AssertExpectations(t)
}
