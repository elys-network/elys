package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryRewardBalanceOfDenom(ctx sdk.Context, query *ammtypes.QueryBalanceRequest) ([]byte, error) {
	denom := query.Denom
	addr := query.Address
	if denom == paramtypes.Elys {
		return nil, errorsmod.Wrap(nil, "invalid denom")
	}

	var balance sdk.Coin
	commitment, found := oq.keeper.GetCommitments(ctx, addr)
	if !found {
		balance = sdk.NewCoin(denom, sdk.ZeroInt())
	} else {
		uncommittedToken, found := commitment.GetRewardsUnclaimedForDenom(denom)
		if !found {
			return nil, errorsmod.Wrap(nil, "invalid denom")
		}

		balance = sdk.NewCoin(denom, uncommittedToken.Amount)
	}

	res := incentivetypes.BalanceAvailable{
		Amount:    balance.Amount.Uint64(),
		UsdAmount: sdk.NewDecFromInt(balance.Amount),
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get reward balance response")
	}
	return responseBytes, nil
}
