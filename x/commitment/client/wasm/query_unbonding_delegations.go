package wasm

import (
	"encoding/json"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (oq *Querier) queryUnbondingDelegations(ctx sdk.Context, query *wasmbindingstypes.QueryDelegatorUnbondingDelegationsRequest) ([]byte, error) {
	if query.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	delAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	unbonding_delegations := oq.stakingKeeper.GetUnbondingDelegations(ctx, delAddr, math.MaxInt16)
	res := wasmbindingstypes.QueryDelegatorUnbondingDelegationsResponse{
		UnbondingResponses: unbonding_delegations,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}

	return responseBytes, nil
}
