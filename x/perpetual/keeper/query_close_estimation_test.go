package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCloseEstimation_InvalidRequest(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	// Mock behavior
	// No mock behavior

	_, err := k.CloseEstimation(sdk.Context{}, nil)
	assert.Error(t, err)
}

func TestCloseEstimation_InvalidAddress(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	var (
		ctx   = sdk.Context{} // Mock or setup a context
		query = &types.QueryCloseEstimationRequest{
			Address:    "invalid_address",
			PositionId: 1,
		}
	)

	// Mock behavior
	// No mock behavior

	_, err := k.CloseEstimation(ctx, query)
	assert.Error(t, err)
}

func TestCloseEstimation_MTPNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	address := sdk.AccAddress([]byte("address"))

	var (
		ctx   = sdk.Context{} // Mock or setup a context
		query = &types.QueryCloseEstimationRequest{
			Address:    address.String(),
			PositionId: 1,
		}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(query.Address), query.PositionId).Return(types.MTP{}, types.ErrMTPDoesNotExist)

	_, err := k.CloseEstimation(ctx, query)
	assert.Error(t, err)

	mockChecker.AssertExpectations(t)
}

func TestCloseEstimation_ExistingMTP(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseLongChecker: mockChecker,
	}

	address := sdk.AccAddress([]byte("address"))

	var (
		ctx   = sdk.Context{} // Mock or setup a context
		query = &types.QueryCloseEstimationRequest{
			Address:    address.String(),
			PositionId: 1,
		}
		mtp = types.MTP{
			AmmPoolId:        2,
			CollateralAsset:  ptypes.BaseCurrency,
			Collateral:       sdk.NewInt(100),
			CustodyAsset:     "uatom",
			Custody:          sdk.NewInt(50),
			LiabilitiesAsset: ptypes.BaseCurrency,
			Liabilities:      sdk.NewInt(400),
			TradingAsset:     "uatom",
			Position:         types.Position_LONG,
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool = ammtypes.Pool{}
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(query.Address), query.PositionId).Return(mtp, nil)
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true)
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil)

	res, err := k.CloseEstimation(ctx, query)
	assert.NoError(t, err)

	mockChecker.AssertExpectations(t)

	assert.Equal(t, mtp.Position, res.Position)
	assert.Equal(t, mtp.Custody, res.PositionSize.Amount)
	assert.Equal(t, mtp.Liabilities, res.Liabilities.Amount)
	assert.Equal(t, sdk.ZeroDec(), res.PriceImpact)
	assert.Equal(t, sdk.ZeroDec(), res.SwapFee)
	assert.Equal(t, sdk.Coin{}, res.RepayAmount)
}
