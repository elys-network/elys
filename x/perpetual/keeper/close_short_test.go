package keeper_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCloseShort_MtpNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(types.MTP{}, types.ErrMTPDoesNotExist)

	_, _, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about the mtp not existing
	assert.True(t, errors.Is(err, types.ErrMTPDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestCloseShort_InvalidCloseSize(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId: 2,
			Custody:   sdk.NewInt(0),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(mtp, nil)

	_, _, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrInvalidCloseSize))
	mockChecker.AssertExpectations(t)
}

func TestCloseShort_PoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId: 2,
			Custody:   sdk.NewInt(100),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(types.Pool{}, false)

	_, _, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrInvalidBorrowingAsset))
	mockChecker.AssertExpectations(t)
}

func TestCloseShort_AmmPoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId: 2,
			Custody:   sdk.NewInt(100),
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(types.Pool{}, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammtypes.Pool{}, errorsmod.Wrap(types.ErrPoolDoesNotExist, mtp.CustodyAsset))

	_, _, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestCloseShort_ErrorSettleBorrowInterest(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:  2,
			Custody:    sdk.NewInt(100),
			Collateral: sdk.NewInt(0),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("SettleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(errors.New("error executing handle borrow interest"))

	_, _, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about handle borrow interest
	assert.Equal(t, errors.New("error executing handle borrow interest"), err)
	mockChecker.AssertExpectations(t)
}

func TestCloseShort_ErrorTakeOutCustody(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:  2,
			Custody:    sdk.NewInt(100),
			Collateral: sdk.NewInt(0),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("SettleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, msg.Amount).Return(errors.New("error executing take out custody"))

	_, _, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about take out custody
	assert.Equal(t, errors.New("error executing take out custody"), err)
	mockChecker.AssertExpectations(t)
}

func TestCloseShort_ErrorEstimateAndRepay(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:  2,
			Custody:    sdk.NewInt(100),
			Collateral: sdk.NewInt(0),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("SettleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, msg.Amount).Return(nil)
	mockChecker.On("EstimateAndRepay", ctx, mtp, pool, ammPool, msg.Amount, ptypes.BaseCurrency).Return(math.Int{}, errors.New("error executing estimate and repay"))

	_, _, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect an error about estimate and repay
	assert.Equal(t, errors.New("error executing estimate and repay"), err)
	mockChecker.AssertExpectations(t)
}

func TestCloseShort_SuccessfulClosingLongPosition(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgClose{
			Creator: "cosmos10duudma7ef9849ee42zhe5q4t4fmk0z99uuh92",
			Id:      1,
			Amount:  sdk.NewInt(100),
		}
		mtp = types.MTP{
			AmmPoolId:  2,
			Custody:    sdk.NewInt(100),
			Collateral: sdk.NewInt(0),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool     = ammtypes.Pool{}
		repayAmount = math.NewInt(100)
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Id).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)
	mockChecker.On("SettleBorrowInterest", ctx, &mtp, &pool, ammPool).Return(nil)
	mockChecker.On("TakeOutCustody", ctx, mtp, &pool, msg.Amount).Return(nil)
	mockChecker.On("EstimateAndRepay", ctx, mtp, pool, ammPool, msg.Amount, ptypes.BaseCurrency).Return(repayAmount, nil)

	mtpOut, repayAmountOut, err := k.CloseShort(ctx, msg, ptypes.BaseCurrency)

	// Expect no error
	assert.Nil(t, err)
	assert.Equal(t, repayAmount, repayAmountOut)
	assert.Equal(t, mtp, *mtpOut)
	mockChecker.AssertExpectations(t)
}
