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
			Amount:  sdk.NewInt(100),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(types.MTP{}, types.ErrMTPDoesNotExist)

	_, _, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about the mtp not existing
	assert.True(t, errors.Is(err, types.ErrMTPDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_InvalidCloseSize(t *testing.T) {
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
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId: 2,
			Custody:   sdk.NewInt(0),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)

	_, _, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrInvalidCloseSize))
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
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId: 2,
			Custody:   sdk.NewInt(100),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(types.Pool{}, false)

	_, _, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

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
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:    2,
			CustodyAsset: "uatom",
			Custody:      sdk.NewInt(100),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(types.Pool{}, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammtypes.Pool{}, sdkerrors.Wrap(types.ErrPoolDoesNotExist, mtp.CustodyAsset))

	_, _, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestCloseLong_ErrorHandleBorrowInterest(t *testing.T) {
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
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:       2,
			CollateralAsset: ptypes.BaseCurrency,
			CustodyAsset:    "uatom",
			Collateral:      sdk.NewInt(0),
			Custody:         sdk.NewInt(100),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("HandleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(errors.New("error executing handle borrow interest"))

	_, _, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about handle borrow interest
	assert.Equal(t, errors.New("error executing handle borrow interest"), err)
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
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:       2,
			CollateralAsset: ptypes.BaseCurrency,
			CustodyAsset:    "uatom",
			Collateral:      sdk.NewInt(0),
			Custody:         sdk.NewInt(100),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("HandleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, msg.Amount).Return(errors.New("error executing take out custody"))

	_, _, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

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
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:       2,
			CollateralAsset: ptypes.BaseCurrency,
			CustodyAsset:    "uatom",
			Collateral:      sdk.NewInt(0),
			Custody:         sdk.NewInt(100),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("HandleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, msg.Amount).Return(nil)
	mockChecker.On("EstimateAndRepay", ctx, mtp, pool, ammPool, msg.Amount, ptypes.BaseCurrency).Return(sdk.Int{}, errors.New("error executing estimate and repay"))

	_, _, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

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
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:       2,
			CollateralAsset: ptypes.BaseCurrency,
			CustodyAsset:    "uatom",
			Collateral:      sdk.NewInt(0),
			Custody:         sdk.NewInt(100),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool     = ammtypes.Pool{}
		repayAmount = math.NewInt(100)
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, msg.Creator, msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("HandleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, msg.Amount).Return(nil)
	mockChecker.On("EstimateAndRepay", ctx, mtp, pool, ammPool, msg.Amount, ptypes.BaseCurrency).Return(repayAmount, nil)

	mtpOut, repayAmountOut, err := k.CloseLong(ctx, msg, ptypes.BaseCurrency)

	// Expect no error
	assert.Nil(t, err)
	assert.Equal(t, repayAmount, repayAmountOut)
	assert.Equal(t, mtp, *mtpOut)
	mockChecker.AssertExpectations(t)
}
