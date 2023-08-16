package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestOpenLong_PoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1),
		}
		poolId = uint64(42)
	)

	// Mock behavior
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, false)

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_PoolDisabled(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1),
		}
		poolId = uint64(42)
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(false)

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect an error about the pool being disabled
	assert.True(t, errors.Is(err, types.ErrMTPDisabled))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_InsufficientAmmPoolBalanceForLeveragedAmount(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(2),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  "uusdc",
			BorrowAsset:      "uatom",
			Position:         types.Position_LONG,
		}
		poolId = uint64(42)
	)

	// Mock the behaviors to get to the HasSufficientPoolBalance check
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, poolId, msg.BorrowAsset).Return(ammtypes.Pool{}, nil) // Assuming a valid pool is returned

	// Mock the behavior where HasSufficientPoolBalance returns false
	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.BorrowAsset, sdk.NewInt(2000)).Return(false) // Example value for sdk.NewInt(100)

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect an error about the borrow amount being too high
	assert.True(t, errors.Is(err, types.ErrBorrowTooHigh))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_InsufficientLiabilities(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(2),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  "uusdc",
			BorrowAsset:      "uatom",
			Position:         types.Position_LONG,
		}
		poolId = uint64(42)
	)

	// Mock the behaviors to get to the CheckMinLiabilities check
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, poolId, msg.BorrowAsset).Return(ammtypes.Pool{}, nil)                          // Assuming a valid pool is returned
	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.BorrowAsset, sdk.NewInt(2000)).Return(true) // Example value for sdk.NewInt(100)

	// Mock the behavior where CheckMinLiabilities returns an error indicating insufficient liabilities
	liabilityError := errors.New("insufficient liabilities")
	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, sdk.NewDec(1), types.Pool{}, ammtypes.Pool{}, msg.BorrowAsset).Return(liabilityError)

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect the custom error indicating insufficient liabilities
	assert.True(t, errors.Is(err, liabilityError))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_InsufficientAmmPoolBalanceForCustody(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  "uusdc",
			BorrowAsset:      "uatom",
			Position:         types.Position_LONG,
		}
		poolId = uint64(42)
	)
	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, poolId, msg.BorrowAsset).Return(ammtypes.Pool{}, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())

	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.BorrowAsset, leveragedAmount).Return(true)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, types.Pool{}, ammtypes.Pool{}, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.CollateralAsset, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.BorrowAsset, ammtypes.Pool{}).Return(custodyAmount, nil)
	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.CollateralAsset, custodyAmount).Return(false)

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect an error about custody amount being too high
	assert.True(t, errors.Is(err, types.ErrCustodyTooHigh))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_ErrorsDuringOperations(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  "uusdc",
			BorrowAsset:      "uatom",
			Position:         types.Position_LONG,
		}
		poolId = uint64(42)
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, poolId, msg.BorrowAsset).Return(ammtypes.Pool{}, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())

	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.BorrowAsset, leveragedAmount).Return(true)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, types.Pool{}, ammtypes.Pool{}, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.CollateralAsset, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.BorrowAsset, ammtypes.Pool{}).Return(custodyAmount, nil)
	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.CollateralAsset, custodyAmount).Return(true)

	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, msg.Leverage, poolId)

	borrowError := errors.New("borrow error")
	mockChecker.On("Borrow", ctx, msg.CollateralAsset, msg.CollateralAmount, custodyAmount, mtp, &ammtypes.Pool{}, &types.Pool{}, eta).Return(borrowError)

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect the borrow error
	assert.True(t, errors.Is(err, borrowError))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_LeverageRatioLessThanSafetyFactor(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  "uusdc",
			BorrowAsset:      "uatom",
			Position:         types.Position_LONG,
		}
		poolId = uint64(42)
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, poolId, msg.BorrowAsset).Return(ammtypes.Pool{}, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())

	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.BorrowAsset, leveragedAmount).Return(true)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, types.Pool{}, ammtypes.Pool{}, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.CollateralAsset, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.BorrowAsset, ammtypes.Pool{}).Return(custodyAmount, nil)
	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.CollateralAsset, custodyAmount).Return(true)

	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, msg.Leverage, poolId)

	mockChecker.On("Borrow", ctx, msg.CollateralAsset, msg.CollateralAmount, custodyAmount, mtp, &ammtypes.Pool{}, &types.Pool{}, eta).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &types.Pool{}).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &types.Pool{}).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("UpdateMTPHealth", ctx, *mtp, ammtypes.Pool{}).Return(lr, nil)
	mockChecker.On("GetSafetyFactor", ctx).Return(sdk.NewDec(100))

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect an error indicating MTP is unhealthy
	assert.True(t, errors.Is(err, types.ErrMTPUnhealthy))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_Success(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  "uusdc",
			BorrowAsset:      "uatom",
			Position:         types.Position_LONG,
		}
		poolId = uint64(42)
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, poolId, msg.BorrowAsset).Return(ammtypes.Pool{}, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())

	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.BorrowAsset, leveragedAmount).Return(true)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, types.Pool{}, ammtypes.Pool{}, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.CollateralAsset, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.BorrowAsset, ammtypes.Pool{}).Return(custodyAmount, nil)
	mockChecker.On("HasSufficientPoolBalance", ctx, ammtypes.Pool{}, msg.CollateralAsset, custodyAmount).Return(true)

	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, msg.Leverage, poolId)

	mockChecker.On("Borrow", ctx, msg.CollateralAsset, msg.CollateralAmount, custodyAmount, mtp, &ammtypes.Pool{}, &types.Pool{}, eta).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &types.Pool{}).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &types.Pool{}).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("UpdateMTPHealth", ctx, *mtp, ammtypes.Pool{}).Return(lr, nil)

	safetyFactor := math.LegacyNewDec(10)

	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor)

	_, err := k.OpenLong(ctx, poolId, msg)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
