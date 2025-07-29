package keeper

import (
	"github.com/elys-network/elys/v7/x/amm/types"
)

// To be used in migration

func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}

func (k Keeper) GetAssetProfileKeeper() types.AssetProfileKeeper {
	return k.assetProfileKeeper
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}
