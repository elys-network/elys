package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (oq *Querier) queryBorrowedAmount(ctx sdk.Context, query *commitmenttypes.QueryBorrowAmountRequest) ([]byte, error) {
	res := types.BalanceBorrowed{
		UsdAmount:  sdk.ZeroDec(),
		Percentage: sdk.ZeroDec(),
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}

	return responseBytes, nil
}
