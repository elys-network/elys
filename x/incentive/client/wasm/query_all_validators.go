package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
)

func (oq *Querier) queryAllValidators(ctx sdk.Context, query *wasmbindingstypes.QueryValidatorsRequest) ([]byte, error) {
	totalBonded := oq.stakingKeeper.TotalBondedTokens(ctx)
	delegatorAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)

	allValidators := oq.stakingKeeper.GetAllValidators(ctx)
	validators := make([]stakingtypes.Validator, 0)
	if err == nil {
		validators = oq.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr, wasmbindingstypes.MAX_RETRY_VALIDATORS)
	}

	validatorsCW := oq.BuildAllValidatorsResponseCW(ctx, allValidators, validators, totalBonded, query.DelegatorAddress)
	res := wasmbindingstypes.QueryDelegatorValidatorsResponse{
		Validators: validatorsCW,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get delegation validator response")
	}

	return responseBytes, nil
}

func (oq *Querier) BuildAllValidatorsResponseCW(ctx sdk.Context, allValidators []stakingtypes.Validator, validators []stakingtypes.Validator, totalBonded cosmos_sdk_math.Int, delegatorAddress string) []wasmbindingstypes.ValidatorDetail {
	var validatorsCW []wasmbindingstypes.ValidatorDetail
	for _, validator := range allValidators {
		isDelegated := false
		for _, v := range validators {
			if validator.OperatorAddress == v.OperatorAddress {
				isDelegated = true
				break
			}
		}

		var validatorCW wasmbindingstypes.ValidatorDetail
		validatorCW.Address = validator.OperatorAddress
		validatorCW.Name = validator.Description.Moniker
		validatorCW.Staked = wasmbindingstypes.BalanceAvailable{
			Amount:    0,
			UsdAmount: sdk.NewDec(0),
		}
		validatorCW.Commission = validator.GetCommission()

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
			// Get validator
			val := oq.stakingKeeper.Validator(ctx, valAddress)

			shares := delegations.GetShares()
			tokens := val.TokensFromSharesTruncated(shares)
			delegatedAmt := tokens.TruncateInt()

			validatorCW.Staked.Amount = delegatedAmt.Uint64()
			validatorCW.Staked.UsdAmount = tokens
		}

		votingPower := sdk.NewDecFromInt(validator.Tokens).QuoInt(totalBonded).MulInt(sdk.NewInt(100))
		validatorCW.VotingPower = votingPower
		validatorCW.ProfilePictureSrc = validator.Description.Website

		validatorsCW = append(validatorsCW, validatorCW)
	}

	return validatorsCW
}
