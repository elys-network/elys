package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
)

func (oq *Querier) queryValidators(ctx sdk.Context, query *wasmbindingstypes.QueryValidatorsRequest) ([]byte, error) {
	validators := oq.stakingKeeper.GetAllValidators(ctx)
	res := wasmbindingstypes.QueryDelegatorValidatorsResponse{
		Validators: validators,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}

	return responseBytes, nil
}
