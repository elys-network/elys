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

func (oq *Querier) queryStakedPositions(ctx sdk.Context, query *commitmenttypes.QueryValidatorsRequest) ([]byte, error) {
	totalBonded := oq.stakingKeeper.TotalBondedTokens(ctx)
	delegatorAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid delegator address")
	}

	validators := make([]stakingtypes.Validator, 0)
	validators = oq.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr, math.MaxInt16)

	stakedPositionsCW := oq.BuildStakedPositionResponseCW(ctx, validators, totalBonded, query.DelegatorAddress)
	res := commitmenttypes.QueryStakedPositionResponse{
		StakedPosition: stakedPositionsCW,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get staked position response")
	}

	return responseBytes, nil
}

func (oq *Querier) BuildStakedPositionResponseCW(ctx sdk.Context, validators []stakingtypes.Validator, totalBonded cosmos_sdk_math.Int, delegatorAddress string) []commitmenttypes.StakedPosition {
	var stakedPositions []commitmenttypes.StakedPosition
	for i, validator := range validators {
		var stakedPosition commitmenttypes.StakedPosition
		stakedPosition.Id = fmt.Sprintf("%d", i+1)

		valAddress, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			continue
		}
		delAddress, err := sdk.AccAddressFromBech32(delegatorAddress)
		if err != nil {
			continue
		}

		delegations, _ := oq.stakingKeeper.GetDelegation(ctx, delAddress, valAddress)

		shares := delegations.GetShares()
		tokens := validator.TokensFromSharesTruncated(shares)
		delegatedAmt := tokens.TruncateInt()
		votingPower := sdk.NewDecFromInt(validator.Tokens).QuoInt(totalBonded).MulInt(sdk.NewInt(100))
		stakedPosition.Validator = commitmenttypes.StakingValidator{
			// The validator address.
			Address: validator.OperatorAddress,
			// The validator name.
			Name: validator.Description.Moniker,
			// Voting power percentage for this validator.
			VotingPower: votingPower,
			// Comission percentage for the validator.
			Commission: validator.GetCommission(),
			// The url of the validator profile picture
			ProfilePictureSrc: validator.Description.Website,
		}
		stakedPosition.Staked = commitmenttypes.BalanceAvailable{
			Amount:    delegatedAmt.Uint64(),
			UsdAmount: tokens,
		}

		stakedPositions = append(stakedPositions, stakedPosition)
	}

	return stakedPositions
}
