package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/vaults/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetVaultPnL(ctx sdk.Context, vaultId string, date string) types.VaultPnL {
	track := types.VaultPnL{}
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.VaultPnLPrefix)
	bz := store.Get(types.GetVaultPnLKey(vaultId, date))
	if len(bz) == 0 {
		return types.VaultPnL{
			VaultId: vaultId,
			Date:    date,
			PnlUsd:  math.LegacyZeroDec(),
		}
	}

	k.cdc.MustUnmarshal(bz, &track)
	return track
}

func (k Keeper) SetVaultPnL(ctx sdk.Context, track types.VaultPnL) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.VaultPnLPrefix)
	bz := k.cdc.MustMarshal(&track)
	store.Set(types.GetVaultPnLKey(track.VaultId, track.Date), bz)
}

func (k Keeper) DeleteVaultPnL(ctx sdk.Context, vaultId string, date string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.VaultPnLPrefix)
	store.Delete(types.GetVaultPnLKey(vaultId, date))
}

func (k Keeper) AddVaultPnL(ctx sdk.Context, track types.VaultPnL) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.VaultPnLPrefix)

	trackTotal := types.VaultPnL{}
	bz := store.Get(types.GetVaultPnLKey(track.VaultId, track.Date))
	if len(bz) == 0 {
		trackTotal = track
	} else {
		k.cdc.MustUnmarshal(bz, &trackTotal)
		trackTotal.PnlUsd = trackTotal.PnlUsd.Add(track.PnlUsd)
	}

	bz = k.cdc.MustMarshal(&trackTotal)
	store.Set(types.GetVaultPnLKey(track.VaultId, track.Date), bz)
}

func (k Keeper) TrackVaultPnL(ctx sdk.Context, vaultId string, pnlUsd math.LegacyDec) {
	track := types.VaultPnL{
		VaultId: vaultId,
		Date:    ctx.BlockTime().Format("2006-01-02"),
		PnlUsd:  pnlUsd,
	}
	k.AddVaultPnL(ctx, track)
}

// Returns last x days total for pnl
func (k Keeper) GetPnlTotal(ctx sdk.Context, vaultId string, days int) osmomath.BigDec {
	start := ctx.BlockTime()
	count := math.ZeroInt()
	total := osmomath.ZeroBigDec()

	for i := 0; i < days; i++ {
		date := start.AddDate(0, 0, i*-1).Format("2006-01-02")
		info := k.GetVaultPnL(ctx, vaultId, date)

		if info.PnlUsd.IsPositive() {
			total = total.Add(osmomath.BigDecFromDec(info.PnlUsd))
			count = count.Add(math.OneInt())
		}
	}

	if count.IsZero() {
		return osmomath.ZeroBigDec()
	}
	return total.Quo(osmomath.BigDecFromSDKInt(count))
}
