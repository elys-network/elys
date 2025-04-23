package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/vaults/types"
)

func (k msgServer) Deposit(goCtx context.Context, req *types.MsgDeposit) (*types.MsgDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgDepositResponse{}, nil
}

func (k Keeper) VaultUsdValue(ctx sdk.Context, vaultId uint64) (sdkmath.LegacyDec, error) {
	vault, found := k.GetVault(ctx, vaultId)
	if !found {
		return sdkmath.LegacyZeroDec(), err
	}
	totalValue := sdkmath.LegacyZeroDec()
	for _, coin := range vault.AllowedCoins {
		totalValue = totalValue.Add(k.tierKeeper.)
	}
	return totalValue,nil
}
