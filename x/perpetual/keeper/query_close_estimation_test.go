package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCloseEstimation_InvalidRequest(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseEstimationChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseEstimationChecker: mockChecker,
	}

	// Mock behavior
	// No mock behavior

	_, err := k.CloseEstimation(sdk.Context{}, nil)
	assert.Error(t, err)
}

func TestCloseEstimation_InvalidAddress(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.CloseEstimationChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseEstimationChecker: mockChecker,
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
	mockChecker := new(mocks.CloseEstimationChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		CloseEstimationChecker: mockChecker,
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
	// Setup the perpetual keeper
	k, ctx, assetProfileKeeper := keepertest.PerpetualKeeper(t)

	// Setup the mock checker
	mockChecker := new(mocks.CloseEstimationChecker)

	// assign the mock checker to the keeper
	k.CloseEstimationChecker = mockChecker

	address := sdk.AccAddress([]byte("address"))

	// get swap fee param
	swapFee := k.GetSwapFee(ctx)

	var (
		query = &types.QueryCloseEstimationRequest{
			Address:    address.String(),
			PositionId: 1,
		}
		mtp = types.MTP{
			AmmPoolId:                      2,
			CollateralAsset:                ptypes.BaseCurrency,
			Collateral:                     sdk.NewInt(100),
			CustodyAsset:                   "uatom",
			Custody:                        sdk.NewInt(50),
			LiabilitiesAsset:               ptypes.BaseCurrency,
			Liabilities:                    sdk.NewInt(400),
			TradingAsset:                   "uatom",
			Position:                       types.Position_LONG,
			BorrowInterestUnpaidCollateral: sdk.NewInt(10),
			OpenPrice:                      sdk.MustNewDecFromStr("1.5"),
		}
		pool = types.Pool{
			BorrowInterestRate: math.LegacyNewDec(2),
		}
		ammPool      = ammtypes.Pool{}
		baseCurrency = "usdc"
	)

	// Mock behavior
	mockChecker.On("GetMTP", ctx, sdk.MustAccAddressFromBech32(query.Address), query.PositionId).Return(mtp, nil).Once()
	mockChecker.On("GetPool", ctx, mtp.AmmPoolId).Return(pool, true).Once()
	mockChecker.On("GetAmmPool", ctx, mtp.AmmPoolId, mtp.CustodyAsset).Return(ammPool, nil).Once()
	mockChecker.On("EstimateSwap", ctx, sdk.NewCoin(mtp.CustodyAsset, mtp.Custody), mtp.CollateralAsset, ammPool).Return(math.NewInt(10000), nil).Once()
	mockChecker.On("EstimateSwapGivenOut", ctx, sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral), baseCurrency, ammPool).Return(math.NewInt(200), nil).Once()
	mockChecker.On("EstimateSwapGivenOut", ctx, sdk.NewCoin(baseCurrency, sdk.NewInt(9400)), mtp.CollateralAsset, ammPool).Return(math.NewInt(9400), nil).Once()
	mockChecker.On("EstimateSwapGivenOut", ctx, sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral), baseCurrency, ammPool).Return(math.NewInt(1111), nil).Once()

	assetProfileKeeper.On("GetEntry", ctx, ptypes.BaseCurrency).Return(atypes.Entry{
		Denom: baseCurrency,
	}, true).Once()

	res, err := k.CloseEstimation(ctx, query)
	assert.NoError(t, err)

	mockChecker.AssertExpectations(t)
	assetProfileKeeper.AssertExpectations(t)

	assert.Equal(t, mtp.Position, res.Position)
	assert.Equal(t, mtp.Custody, res.PositionSize.Amount)
	assert.Equal(t, mtp.Liabilities, res.Liabilities.Amount)
	assert.Equal(t, sdk.ZeroDec(), res.PriceImpact)
	assert.Equal(t, swapFee, res.SwapFee)
	assert.Equal(t, sdk.NewCoin(mtp.CollateralAsset, sdk.NewInt(9400)), res.ReturnAmount)
}
