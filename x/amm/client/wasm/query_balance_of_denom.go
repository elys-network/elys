package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryBalanceOfDenom(ctx sdk.Context, query *types.QueryBalanceRequest) ([]byte, error) {
	denom := query.Denom
	addr := query.Address
	address, err := sdk.AccAddressFromBech32(query.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}
	balance := oq.bankKeeper.GetBalance(ctx, address, denom)
	if denom != paramtypes.Elys {
		commitment, found := oq.commitmentKeeper.GetCommitments(ctx, addr)
		if !found {
			balance = sdk.NewCoin(denom, sdk.ZeroInt())
		} else {
			rewardUnclaimed := commitment.GetRewardUnclaimedForDenom(denom)
			balance = sdk.NewCoin(denom, rewardUnclaimed)
		}
	}

	res := types.QueryBalanceResponse{
		Balance: balance,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}
	return responseBytes, nil
}
