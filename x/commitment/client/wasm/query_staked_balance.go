package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryStakedBalanceOfDenom(ctx sdk.Context, query *ammtypes.QueryBalanceRequest) ([]byte, error) {
	denom := query.Denom
	addr := query.Address
	address, err := sdk.AccAddressFromBech32(query.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	bondedAmt := oq.stakingKeeper.GetDelegatorBonded(ctx, address)
	balance := sdk.NewCoin(denom, bondedAmt)
	if denom != paramtypes.Elys {
		commitment := oq.keeper.GetCommitments(ctx, addr)
		committedToken := commitment.GetCommittedAmountForDenom(denom)

		balance = sdk.NewCoin(denom, committedToken)
	}

	res := commitmenttypes.BalanceAvailable{
		Amount:    balance.Amount,
		UsdAmount: sdk.NewDecFromInt(balance.Amount),
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get staked balance response")
	}
	return responseBytes, nil
}
