package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (oq *Querier) queryAirdropAll(ctx sdk.Context, query *types.QueryAllAirdropRequest) ([]byte, error) {
	res, err := oq.keeper.AirdropAll(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get airdrop all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize airdrop all response")
	}
	return responseBytes, nil
}
