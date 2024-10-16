package keeper_test

import (
	"cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func createNBorrowRate(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.InterestBlock, int64) {
	items := make([]types.InterestBlock, n)
	ctx = ctx.WithBlockHeight(1000)
	curBlock := ctx.BlockHeight()
	for i := range items {
		items[i].InterestRate = sdk.NewDec(int64(i + 1)) // Start from 1 to avoid zero interest
		items[i].BlockHeight = int64(i * 10)

		curBlock++
		keeper.SetBorrowRate(ctx, uint64(curBlock), 1, items[i]) // Assuming pool ID 1
	}
	return items, curBlock
}

func TestBorrowRateGet(t *testing.T) {
	keeper, ctx, _ := keepertest.PerpetualKeeper(t)
	_, lastBlock := createNBorrowRate(keeper, ctx, 10)
	ctx = ctx.WithBlockHeight(lastBlock)

	// 1st case: recent block
	res := keeper.GetBorrowInterestRate(ctx, uint64(ctx.BlockHeight()-2), 1, math.LegacyOneDec())
	require.Equal(t, sdk.NewDec(19), res) // 19

	// 2nd case: older block
	res = keeper.GetBorrowInterestRate(ctx, uint64(ctx.BlockHeight()-8), 1, math.LegacyOneDec())
	require.Equal(t, sdk.NewDec(52), res) // 52

	// 3rd case: future block (should return zero)
	res = keeper.GetBorrowInterestRate(ctx, uint64(ctx.BlockHeight()+10), 1, math.LegacyOneDec())
	require.Equal(t, sdk.ZeroDec(), res)

	// 4th case: non-existent pool
	res = keeper.GetBorrowInterestRate(ctx, uint64(ctx.BlockHeight()-2), 2, math.LegacyOneDec())
	require.Equal(t, sdk.ZeroDec(), res)
}
