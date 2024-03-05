package wasm

import (
	"encoding/json"
	"math"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (oq *Querier) queryDelegatorValidators(ctx sdk.Context, query *commitmenttypes.QueryValidatorsRequest) ([]byte, error) {
	totalBonded := oq.stakingKeeper.TotalBondedTokens(ctx)
	delegatorAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid delegator address")
	}

	validators := make([]stakingtypes.Validator, 0)
	validators = oq.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr, math.MaxInt16)

	validatorsCW := oq.BuildDelegatorValidatorsResponseCW(ctx, validators, totalBonded, query.DelegatorAddress)
	res := commitmenttypes.QueryDelegatorValidatorsResponse{
		Validators: validatorsCW,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get delegation validator response")
	}

	return responseBytes, nil
}

func (oq *Querier) BuildDelegatorValidatorsResponseCW(ctx sdk.Context, validators []stakingtypes.Validator, totalBonded cosmos_sdk_math.Int, delegatorAddress string) []commitmenttypes.ValidatorDetail {
	var validatorsCW []commitmenttypes.ValidatorDetail
	for _, validator := range validators {
		var validatorCW commitmenttypes.ValidatorDetail
		validatorCW.Id = validator.Description.Identity
		validatorCW.Address = validator.OperatorAddress
		validatorCW.Name = validator.Description.Moniker
		validatorCW.Commission = validator.GetCommission()

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

		validatorCW.Staked = commitmenttypes.BalanceAvailable{
			Amount:    tokens.TruncateInt(),
			UsdAmount: tokens,
		}

		votingPower := sdk.NewDecFromInt(validator.Tokens).QuoInt(totalBonded).MulInt(sdk.NewInt(100))
		validatorCW.VotingPower = votingPower

		validatorsCW = append(validatorsCW, validatorCW)
	}

	return validatorsCW
}
