package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (k Keeper) CheckPoolADL(ctx sdk.Context, leveragePool types.Pool) (bool, error) {
	ammPool, err := k.GetAmmPool(ctx, leveragePool.AmmPoolId)
	if err != nil {
		return false, err
	}
	ratio := leveragePool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	params := k.GetParams(ctx)
	maxRatio := leveragePool.MaxLeveragelpRatio.Add(params.ExitBuffer)
	if ratio.GT(maxRatio) {
		return true, nil
	}
	return false, nil
}

func (k Keeper) ClosePositionsOnADL(ctx sdk.Context, leveragePool types.Pool) {

}
