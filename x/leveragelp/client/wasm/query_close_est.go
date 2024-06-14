package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (oq *Querier) queryCloseEst(ctx sdk.Context, query *types.QueryCloseEstRequest) ([]byte, error) {
	res, err := oq.keeper.CloseEst(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get leveragelp close estimation")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize leveragelp close estimation response")
	}
	return responseBytes, nil
}
