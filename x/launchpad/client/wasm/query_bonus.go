package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
)

func (oq *Querier) queryBonus(ctx sdk.Context, query *types.QueryBonusRequest) ([]byte, error) {
	res, err := oq.keeper.Bonus(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get user bonus")
	}
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize user bonus")
	}
	return responseBytes, nil
}
