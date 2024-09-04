package keeper_test

import (
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
		items[i].BlockTime = int64(i * 10)

		curBlock++
		keeper.SetBorrowRate(ctx, uint64(curBlock), 1, items[i]) // Assuming pool ID 1
	}
	return items, curBlock
}

func TestBorrowRateGet(t *testing.T) {
	keeper, ctx := keepertest.PerpetualKeeper(t)
	_, lastBlock := createNBorrowRate(keeper, ctx, 10)
	ctx = ctx.WithBlockHeight(lastBlock)

	// 1st case: recent block
	res := keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()-2), 1, sdk.NewDec(1000))
	require.Equal(t, sdk.NewDec(19000), res) // 19 * 1000

	// 2nd case: older block
	res = keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()-8), 1, sdk.NewDec(1000))
	require.Equal(t, sdk.NewDec(52000), res) // 52 * 1000

	// 3rd case: future block (should return zero)
	res = keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()+10), 1, sdk.NewDec(1000))
	require.Equal(t, sdk.ZeroDec(), res)

	// 4th case: non-existent pool
	res = keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()-2), 2, sdk.NewDec(1000))
	require.Equal(t, sdk.ZeroDec(), res)
}
