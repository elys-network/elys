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

func (oq *Querier) queryUnbondingDelegations(ctx sdk.Context, query *types.QueryDelegatorUnbondingDelegationsRequest) ([]byte, error) {
	if query.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	delAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	unbonding_delegations, err := oq.stakingKeeper.GetUnbondingDelegations(ctx, delAddr, math.MaxInt16)
	if err != nil {
		return nil, err
	}
	unbonding_delegations_cw := BuildUnbondingDelegationResponseCW(unbonding_delegations)
	res := types.QueryDelegatorUnbondingDelegationsResponse{
		UnbondingResponses: unbonding_delegations_cw,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get balance response")
	}

	return responseBytes, nil
}

func BuildUnbondingDelegationResponseCW(unbondingDelegations []stakingtypes.UnbondingDelegation) []types.UnbondingDelegation {
	var unbondingDelegationsCW []types.UnbondingDelegation
	for _, unbondingDelegation := range unbondingDelegations {
		var unbondingDelegationCW types.UnbondingDelegation
		unbondingDelegationCW.DelegatorAddress = unbondingDelegation.DelegatorAddress
		unbondingDelegationCW.ValidatorAddress = unbondingDelegation.ValidatorAddress

		for _, entity := range unbondingDelegation.Entries {
			newEntity := types.UnbondingDelegationEntry{
				// creation_height is the height which the unbonding took place.
				CreationHeight: entity.CreationHeight,
				// completion_time is the unix time for unbonding completion.
				CompletionTime: entity.CompletionTime.Unix(),
				// initial_balance defines the tokens initially scheduled to receive at completion.
				InitialBalance: entity.InitialBalance,
				// balance defines the tokens to receive at completion.
				Balance: entity.Balance,
				// Incrementing id that uniquely identifies this entry
				UnbondingId: entity.UnbondingId,
			}

			unbondingDelegationCW.Entries = append(unbondingDelegationCW.Entries, newEntity)
		}

		unbondingDelegationsCW = append(unbondingDelegationsCW, unbondingDelegationCW)
	}

	return unbondingDelegationsCW
}
