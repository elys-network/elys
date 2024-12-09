package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_AtomStakers(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	val := types.AtomStaker{
		Address: addr.String(),
		Amount:  math.OneInt().MulRaw(21),
	}

	keeper.SetAtomStaker(ctx, val)

	retrieved := keeper.GetAtomStaker(ctx, addr)
	assert.Equal(t, val, retrieved)

	retrievedNone := keeper.GetAtomStaker(ctx, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()))
	assert.Equal(t, math.ZeroInt(), retrievedNone.Amount)

	all := keeper.GetAllAtomStakers(ctx)
	assert.Equal(t, len(all), 1)
}

func TestKeeper_NftHolders(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	val := types.NftHolder{
		Address: addr.String(),
		Amount:  math.OneInt().MulRaw(21),
	}

	keeper.SetNFTHodler(ctx, val)

	retrieved := keeper.GetNFTHolder(ctx, addr)
	assert.Equal(t, val, retrieved)

	retrievedNone := keeper.GetNFTHolder(ctx, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()))
	assert.Equal(t, math.ZeroInt(), retrievedNone.Amount)

	all := keeper.GetAllNFTHolders(ctx)
	assert.Equal(t, len(all), 1)
}

func TestKeeper_Governors(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	val := types.Governor{
		Address: addr.String(),
		Amount:  math.OneInt().MulRaw(21),
	}

	keeper.SetGovernor(ctx, val)

	retrieved := keeper.GetGovernor(ctx, addr)
	assert.Equal(t, val, retrieved)

	retrievedNone := keeper.GetGovernor(ctx, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()))
	assert.Equal(t, math.ZeroInt(), retrievedNone.Amount)

	all := keeper.GetAllGovernors(ctx)
	assert.Equal(t, len(all), 1)
}

func TestKeeper_Cadets(t *testing.T) {
	keeper, ctx := keepertest.CommitmentKeeper(t)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	val := types.Cadet{
		Address: addr.String(),
		Amount:  math.OneInt().MulRaw(21),
	}

	keeper.SetCadet(ctx, val)

	retrieved := keeper.GetCadet(ctx, addr)
	assert.Equal(t, val, retrieved)

	retrievedNone := keeper.GetCadet(ctx, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()))
	assert.Equal(t, math.ZeroInt(), retrievedNone.Amount)

	all := keeper.GetAllCadets(ctx)
	assert.Equal(t, len(all), 1)
}
