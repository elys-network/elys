package wasm

import (
	"encoding/json"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (oq *Querier) queryDelegations(ctx sdk.Context, query *wasmbindingstypes.QueryDelegatorDelegationsRequest) ([]byte, error) {
	if query.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	delAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	delegations := oq.stakingKeeper.GetDelegatorDelegations(ctx, delAddr, math.MaxInt16)
	delegationResps, err := stakingkeeper.DelegationsToDelegationResponses(ctx, oq.stakingKeeper, delegations)

	res := wasmbindingstypes.QueryDelegatorDelegationsResponse{
		DelegationResponses: delegationResps,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}

	return responseBytes, nil
}
