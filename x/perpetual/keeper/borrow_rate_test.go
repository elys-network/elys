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

	params := keeper.GetParams(ctx)
	params.FixedFundingRate = sdk.NewDec(0)
	keeper.SetParams(ctx, &params)

	for i := range items {
		items[i].InterestRate = sdk.NewDec(int64(i + 1)) // Start from 1 to avoid zero interest
		items[i].BlockHeight = int64(i * 10)
		items[i].BlockTime = int64(i)

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
	res := keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()-2), uint64(ctx.BlockTime().Unix())-(86400*365), 1, sdk.NewDec(1000))
	require.Equal(t, sdk.NewDec(9500), res) // 19 * 1000 / 2

	// 2nd case: older block
	res = keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()-8), uint64(ctx.BlockTime().Unix())-(86400*365), 1, sdk.NewDec(1000))
	require.Equal(t, sdk.NewDec(6500), res) // 52 * 1000 / 8

	// 3rd case: future block (should return zero)
	res = keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()+10), uint64(ctx.BlockTime().Unix())-(86400*365), 1, sdk.NewDec(1000))
	require.Equal(t, sdk.ZeroDec(), res)

	// 4th case: non-existent pool
	res = keeper.GetBorrowRate(ctx, uint64(ctx.BlockHeight()-2), uint64(ctx.BlockTime().Unix())-(86400*365), 2, sdk.NewDec(1000))
	require.Equal(t, sdk.ZeroDec(), res)
}
