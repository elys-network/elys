package querier

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (oq *Querier) queryBalanceOfDenom(ctx sdk.Context, query *ammtypes.QueryBalanceRequest) ([]byte, error) {
	res, err := oq.keeper.Balance(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize balance response")
	}
	return responseBytes, nil
}
