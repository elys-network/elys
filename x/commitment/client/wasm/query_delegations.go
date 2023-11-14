package wasm

import (
	"encoding/json"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/commitment/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (oq *Querier) queryDelegations(ctx sdk.Context, query *types.QueryDelegatorDelegationsRequest) ([]byte, error) {
	if query.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	delAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	delegations := oq.stakingKeeper.GetDelegatorDelegations(ctx, delAddr, math.MaxInt16)
	delegationResps, err := oq.DelegationsToDelegationResponses(ctx, delegations)

	res := types.QueryDelegatorDelegationsResponse{
		DelegationResponses: delegationResps,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}

	return responseBytes, nil
}

func (oq *Querier) DelegationsToDelegationResponses(ctx sdk.Context, delegations stakingtypes.Delegations) ([]types.DelegationResponse, error) {
	resp := make([]types.DelegationResponse, len(delegations))

	for i, del := range delegations {
		delResp, err := oq.DelegationToDelegationResponse(ctx, del)
		if err != nil {
			return nil, err
		}

		resp[i] = delResp
	}

	return resp, nil
}

func (oq *Querier) DelegationToDelegationResponse(ctx sdk.Context, del stakingtypes.Delegation) (types.DelegationResponse, error) {
	val, found := oq.stakingKeeper.GetValidator(ctx, del.GetValidatorAddr())
	if !found {
		return types.DelegationResponse{}, stakingtypes.ErrNoValidatorFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
	if err != nil {
		return types.DelegationResponse{}, err
	}

	return types.DelegationResponse{
		Delegation: types.Delegation{
			DelegatorAddress: delegatorAddress.String(),
			ValidatorAddress: del.GetValidatorAddr().String(),
			Shares:           del.Shares,
		},
		Balance: sdk.NewCoin(oq.stakingKeeper.BondDenom(ctx), val.TokensFromShares(del.Shares).TruncateInt()),
	}, nil
}
