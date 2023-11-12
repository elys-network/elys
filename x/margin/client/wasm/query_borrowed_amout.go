package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryBorrowedAmount(ctx sdk.Context, query *wasmbindingstypes.QueryBorrowAmountRequest) ([]byte, error) {
	res := wasmbindingstypes.QueryBalanceResponse{
		Balance: sdk.NewCoin(paramtypes.BaseCurrency, sdk.NewInt(0)),
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}

	return responseBytes, nil
}
