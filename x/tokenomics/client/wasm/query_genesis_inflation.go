package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (oq *Querier) queryGenesisInflation(ctx sdk.Context, query *types.QueryGetGenesisInflationRequest) ([]byte, error) {
	res, err := oq.keeper.GenesisInflation(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get genesis inflation")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize genesis inflation response")
	}
	return responseBytes, nil
}
