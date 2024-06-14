package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
)

func (oq *Querier) queryCalculateDiscount(ctx sdk.Context, query *types.QueryCalculateDiscountRequest) ([]byte, error) {
	res, err := oq.keeper.CalculateDiscount(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to calculate discount")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize calculate discount response")
	}
	return responseBytes, nil
}
