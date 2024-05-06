package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (oq *Querier) queryNumberOfCommitments(ctx sdk.Context, query *commitmenttypes.QueryNumberOfCommitmentsRequest) ([]byte, error) {
	res, err := oq.keeper.NumberOfCommitments(ctx, query)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to number of commitments response")
	}
	return responseBytes, nil
}
