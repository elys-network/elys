package keeper_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/x/assetprofile/keeper"
	"github.com/elys-network/elys/v4/x/assetprofile/types"
)

func createEntries(keeper *keeper.Keeper, ctx sdk.Context, entries []string) {
	var entry types.Entry
	for _, e := range entries {
		entry.BaseDenom = e
		keeper.SetEntry(ctx, entry)
	}
}

func TestFixEntries(t *testing.T) {
	keeper, ctx := keepertest.AssetprofileKeeper(t)
	list := []string{"uatom", "uusdc", "ibc/45D6B52CAD911A15BD9C2F5FFDA80E26AFCB05C7CD520070790ABC86D2B24229", "ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349"}
	createEntries(keeper, ctx, list)
	allEntries := keeper.GetAllEntry(ctx)
	require.Equal(t, 4, len(allEntries))
	for _, e := range allEntries {
		require.True(t, slices.Contains(list, e.BaseDenom))
	}
	keeper.FixEntries(ctx)
	finalList := []string{"uatom", "uusdc"}
	allEntries = keeper.GetAllEntry(ctx)
	require.Equal(t, 2, len(allEntries))
	for _, e := range allEntries {
		require.True(t, slices.Contains(finalList, e.BaseDenom))
	}
}
