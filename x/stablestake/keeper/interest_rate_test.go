package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

func createNInterest(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.InterestBlock, int64) {
	items := make([]types.InterestBlock, n)
	ctx = ctx.WithBlockHeight(1000)
	curBlock := ctx.BlockHeight()
	for i := range items {
		items[i].InterestRate = sdk.NewDec(int64(i))
		items[i].BlockTime = int64(i * 10)

		curBlock++
		keeper.SetInterest(ctx, uint64(curBlock), items[i])
	}
	return items, curBlock
}

func TestInterestGet(t *testing.T) {
	keeper, ctx := keepertest.StablestakeKeeper(t)
	_, lastBlock := createNInterest(keeper, ctx, 10)
	ctx = ctx.WithBlockHeight(lastBlock)

	// 1st case
	res := keeper.GetInterest(ctx, uint64(ctx.BlockHeight()-2), uint64(ctx.BlockTime().Unix()-1), sdk.NewDec(86400*365))
	require.Equal(t, res.Int64(), int64(8))

	// 2nd case
	res = keeper.GetInterest(ctx, uint64(ctx.BlockHeight()-20), uint64(ctx.BlockTime().Unix()-1), sdk.NewDec(86400*365))
	require.Equal(t, res.Int64(), int64(2))

	// 3rd case
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1000)
	res = keeper.GetInterest(ctx, uint64(ctx.BlockHeight()-20), uint64(ctx.BlockTime().Unix()-1), sdk.NewDec(86400*365))
	require.Equal(t, res.Int64(), int64(0))
}
