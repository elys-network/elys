package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
)

func (oq *Querier) queryDenomLiquidityAll(ctx sdk.Context, denomLiquidityAll *types.QueryAllDenomLiquidityRequest) ([]byte, error) {
	res, err := oq.keeper.DenomLiquidityAll(ctx, denomLiquidityAll)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get all denom liquidity")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize all denom liquidity response")
	}
	return responseBytes, nil
}
