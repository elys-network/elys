package wasm

import (
	"encoding/json"
	"math"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (oq *Querier) queryAllValidators(ctx sdk.Context, query *commitmenttypes.QueryValidatorsRequest) ([]byte, error) {
	totalBonded := oq.stakingKeeper.TotalBondedTokens(ctx)
	delegatorAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)

	allValidators := oq.stakingKeeper.GetAllValidators(ctx)
	validators := make([]stakingtypes.Validator, 0)
	if err == nil {
		validators = oq.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr, math.MaxInt16)
	}

	validatorsCW := oq.BuildAllValidatorsResponseCW(ctx, allValidators, validators, totalBonded, query.DelegatorAddress)
	res := commitmenttypes.QueryDelegatorValidatorsResponse{
		Validators: validatorsCW,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get delegation validator response")
	}

	return responseBytes, nil
}

func (oq *Querier) BuildAllValidatorsResponseCW(ctx sdk.Context, allValidators []stakingtypes.Validator, validators []stakingtypes.Validator, totalBonded cosmos_sdk_math.Int, delegatorAddress string) []commitmenttypes.ValidatorDetail {
	var validatorsCW []commitmenttypes.ValidatorDetail
	for _, validator := range allValidators {
		isDelegated := false
		for _, v := range validators {
			if validator.OperatorAddress == v.OperatorAddress {
				isDelegated = true
				break
			}
		}

		var validatorCW commitmenttypes.ValidatorDetail
		validatorCW.Id = validator.Description.Identity
		validatorCW.Address = validator.OperatorAddress
		validatorCW.Name = validator.Description.Moniker
		validatorCW.Staked = commitmenttypes.BalanceAvailable{
			Amount:    sdk.ZeroInt(),
			UsdAmount: sdk.ZeroDec(),
		}
		validatorCW.Commission = validator.GetCommission()
		validatorCW.Jailed = strconv.FormatBool(validator.Jailed)
		validatorCW.Inactive = strconv.FormatBool(!validator.IsBonded())
		// if there is delegation,
		if isDelegated {
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

			validatorCW.Staked.Amount = tokens.TruncateInt()
			validatorCW.Staked.UsdAmount = tokens
		}

		votingPower := sdk.NewDecFromInt(validator.Tokens).QuoInt(totalBonded).MulInt(sdk.NewInt(100))
		validatorCW.VotingPower = votingPower

		validatorsCW = append(validatorsCW, validatorCW)
	}

	return validatorsCW
}
