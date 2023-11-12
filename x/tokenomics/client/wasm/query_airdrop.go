package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (oq *Querier) queryAirdrop(ctx sdk.Context, query *types.QueryGetAirdropRequest) ([]byte, error) {
	res, err := oq.keeper.Airdrop(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get airdrop")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize airdrop response")
	}
	return responseBytes, nil
}
