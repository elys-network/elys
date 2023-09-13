package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestCheckSameAssets_NewPosition(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PositionChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		PositionChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context
	mtp := types.NewMTP("creator", ptypes.USDC, ptypes.ATOM, types.Position_LONG, sdk.NewDec(5), 1)
	k.SetMTP(ctx, mtp)

	msg := &types.MsgOpen{
		Creator:          "creator",
		CollateralAsset:  ptypes.ATOM,
		CollateralAmount: sdk.NewInt(100),
		BorrowAsset:      ptypes.ATOM,
		Position:         types.Position_SHORT,
		Leverage:         sdk.NewDec(1),
	}

	mtp = k.CheckSamePosition(ctx, msg)

	// Expect no error
	assert.Nil(t, mtp)
	mockChecker.AssertExpectations(t)
}
