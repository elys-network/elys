package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (oq *Querier) queryBorrowRatio(ctx sdk.Context, query *types.QueryBorrowRatioRequest) ([]byte, error) {
	res, err := oq.keeper.BorrowRatio(sdk.WrapSDKContext(ctx), query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get params")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize params response")
	}
	return responseBytes, nil
}
