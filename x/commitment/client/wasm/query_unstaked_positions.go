package wasm

import (
	"encoding/json"
	"fmt"
	"math"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (oq *Querier) queryUnStakedPositions(ctx sdk.Context, query *commitmenttypes.QueryValidatorsRequest) ([]byte, error) {
	totalBonded := oq.stakingKeeper.TotalBondedTokens(ctx)
	delegatorAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid delegator address")
	}

	unbonding_delegations := oq.stakingKeeper.GetUnbondingDelegations(ctx, delegatorAddr, math.MaxInt16)

	unstakedPositionsCW := oq.BuildUnStakedPositionResponseCW(ctx, unbonding_delegations, totalBonded, query.DelegatorAddress)
	res := commitmenttypes.QueryUnstakedPositionResponse{
		UnstakedPosition: unstakedPositionsCW,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get delegation validator response")
	}

	return responseBytes, nil
}

func (oq *Querier) BuildUnStakedPositionResponseCW(ctx sdk.Context, unbondingDelegations []stakingtypes.UnbondingDelegation, totalBonded cosmos_sdk_math.Int, delegatorAddress string) []commitmenttypes.UnstakedPosition {
	var unstakedPositions []commitmenttypes.UnstakedPosition
	i := 1
	for _, undelegation := range unbondingDelegations {
		for _, entity := range undelegation.Entries {
			var unstakedPosition commitmenttypes.UnstakedPosition
			unstakedPosition.Id = fmt.Sprintf("%d", i)

			valAddress, err := sdk.ValAddressFromBech32(undelegation.ValidatorAddress)
			if err != nil {
				continue
			}

			// Get validator
			val := oq.stakingKeeper.Validator(ctx, valAddress).(stakingtypes.Validator)
			votingPower := sdk.NewDecFromInt(val.GetBondedTokens()).QuoInt(totalBonded).MulInt(sdk.NewInt(100))

			website := val.Description.Website
			if len(website) < 1 {
				website = " "
			}

			unstakedPosition.Validator = commitmenttypes.StakingValidator{
				// The validator address.
				Address: val.OperatorAddress,
				// The validator name.
				Name: val.Description.Moniker,
				// Voting power percentage for this validator.
				VotingPower: votingPower,
				// Comission percentage for the validator.
				Commission: val.GetCommission(),
				// The url of the validator profile picture
				ProfilePictureSrc: website,
			}
			unstakedPosition.RemainingTime = uint64(entity.CompletionTime.Unix())
			unstakedPosition.Unstaked = commitmenttypes.BalanceAvailable{
				Amount:    entity.Balance,
				UsdAmount: sdk.NewDecFromInt(entity.Balance),
			}

			unstakedPositions = append(unstakedPositions, unstakedPosition)
			i++
		}

	}

	return unstakedPositions
}
