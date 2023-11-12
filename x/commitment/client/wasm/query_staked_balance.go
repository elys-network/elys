package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryStakedBalanceOfDenom(ctx sdk.Context, query *wasmbindingstypes.QueryBalanceRequest) ([]byte, error) {
	denom := query.Denom
	addr := query.Address
	address, err := sdk.AccAddressFromBech32(query.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	bondedAmt := oq.stakingKeeper.GetDelegatorBonded(ctx, address)
	balance := sdk.NewCoin(denom, bondedAmt)
	if denom != paramtypes.Elys {
		commitment, found := oq.keeper.GetCommitments(ctx, addr)
		if !found {
			balance = sdk.NewCoin(denom, sdk.ZeroInt())
		} else {
			committedToken, found := commitment.GetCommittedTokensForDenom(denom)
			if !found {
				return nil, errorsmod.Wrap(nil, "invalid denom")
			}

			balance = sdk.NewCoin(denom, committedToken.Amount)
		}
	}

	res := wasmbindingstypes.QueryBalanceResponse{
		Balance: balance,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get staked balance response")
	}
	return responseBytes, nil
}
