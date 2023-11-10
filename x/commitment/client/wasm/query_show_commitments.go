package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (oq *Querier) queryShowCommitments(ctx sdk.Context, req *types.QueryShowCommitmentsRequest) ([]byte, error) {
	res, err := oq.keeper.ShowCommitments(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to show commitments")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize show commitments response")
	}
	return responseBytes, nil
}
