package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
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
	lockups := make([]commitmenttypes.Lockup, 0)

	if denom != paramtypes.Elys {
		commitment := oq.keeper.GetCommitments(ctx, addr)
		if denom == paramtypes.BaseCurrency {
			denom = stabletypes.GetShareDenom()
		}

		committedToken := commitment.GetCommittedAmountForDenom(denom)
		lockups = commitment.GetCommittedLockUpsForDenom(denom)
		balance = sdk.NewCoin(denom, committedToken)
	}

	resp := commitmenttypes.StakedAvailable{
		Amount:    balance.Amount,
		UsdAmount: sdk.NewDecFromInt(balance.Amount),
		Lockups:   lockups,
	}

	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get staked balance response")
	}

	return responseBytes, nil
}
