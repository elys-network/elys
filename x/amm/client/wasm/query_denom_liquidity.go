package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
)

func (oq *Querier) queryDenomLiquidity(ctx sdk.Context, denomLiquidity *types.QueryGetDenomLiquidityRequest) ([]byte, error) {
	res, err := oq.keeper.DenomLiquidity(ctx, denomLiquidity)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get denom liquidity")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize denom liquidity response")
	}
	return responseBytes, nil
}
