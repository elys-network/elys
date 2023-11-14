package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryRewardBalanceOfDenom(ctx sdk.Context, query *ammtypes.QueryBalanceRequest) ([]byte, error) {
	denom := query.Denom
	addr := query.Address
	if denom == paramtypes.Elys {
		return nil, errorsmod.Wrap(nil, "invalid denom")
	}

	var balance sdk.Coin
	commitment := oq.keeper.GetCommitments(ctx, addr)
	uncommittedToken := commitment.GetRewardUnclaimedForDenom(denom)
	balance = sdk.NewCoin(denom, uncommittedToken)

	res := commitmenttypes.BalanceAvailable{
		Amount:    balance.Amount.Uint64(),
		UsdAmount: sdk.NewDecFromInt(balance.Amount),
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get reward balance response")
	}
	return responseBytes, nil
}
