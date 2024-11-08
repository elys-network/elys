package migrations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/amm/utils"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPool(ctx)
	for _, pool := range pools {
		newPoolAddress := types.NewPoolAddress(pool.PoolId)
		poolAccountModuleName := types.GetPoolIdModuleName(pool.PoolId)
		if err := utils.CreateModuleAccount(ctx, m.keeper.GetAccountKeeper(), newPoolAddress, poolAccountModuleName); err != nil {
			panic(fmt.Errorf("error creating new pool account for %d: %w", pool.PoolId, err))
		}

		// Bank: Transfer funds from prevPoolAddress to new newPoolAddress
		prevPoolAddressBalances := m.keeper.GetBankKeeper().GetAllBalances(ctx, sdk.AccAddress(pool.Address))
		m.keeper.GetBankKeeper().SendCoins(ctx.Context(), sdk.AccAddress(pool.GetAddress()), newPoolAddress, prevPoolAddressBalances)

		// AssetProfile: Update authority in assetprofile entry
		poolBaseDenom := types.GetPoolShareDenom(pool.PoolId)

		entry, found := m.keeper.GetAssetProfileKeeper().GetEntry(ctx, poolBaseDenom)
		// Should not happen
		if !found {
			panic(fmt.Errorf("assetprofile not found for basedenom: %s", poolBaseDenom))
		}

		entry.Authority = newPoolAddress.String()
		m.keeper.GetAssetProfileKeeper().SetEntry(ctx, entry)

		oldPoolAccount := m.keeper.GetAccountKeeper().GetAccount(ctx.Context(), sdk.AccAddress(pool.Address))
		m.keeper.GetAccountKeeper().RemoveAccount(ctx.Context(), oldPoolAccount)
		pool.Address = newPoolAddress.String()

		m.keeper.SetPool(ctx, pool)

		// Update the name and Symbol of pool share token metadata
		metadata, found := m.keeper.GetBankKeeper().GetDenomMetaData(ctx, poolBaseDenom)
		// Should not happen
		if !found {
			panic(fmt.Errorf("denom metadata for poolshare denom not found in bank denom metadata: %s", poolBaseDenom))
		}
		metadata.Name = metadata.Base
		metadata.Symbol = metadata.Display

		m.keeper.GetBankKeeper().SetDenomMetaData(ctx, metadata)
	}

	return nil
}
