package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
)

func (oq *Querier) queryBalance(ctx sdk.Context, query *ammtypes.QueryBalanceRequest) ([]byte, error) {
	res, err := oq.keeper.Balance(sdk.WrapSDKContext(ctx), query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance")
	}

	balance := res.Balance
	resp := incentivetypes.BalanceAvailable{
		Amount:    balance.Amount.Uint64(),
		UsdAmount: sdk.NewDecFromInt(balance.Amount),
	}

	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize balance response")
	}
	return responseBytes, nil
}
