package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryBalanceOfDenom(ctx sdk.Context, query *wasmbindingstypes.QueryBalanceRequest) ([]byte, error) {
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
			rewardUnclaimed, found := commitment.GetRewardsUnclaimedForDenom(denom)
			if !found {
				return nil, errorsmod.Wrap(nil, "invalid denom")
			}

			balance = sdk.NewCoin(denom, rewardUnclaimed.Amount)
		}
	}

	res := wasmbindingstypes.QueryBalanceResponse{
		Balance: balance,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}
	return responseBytes, nil
}
