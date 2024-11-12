package wasm

import (
	"encoding/json"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (oq *Querier) queryDelegatorValidators(ctx sdk.Context, query *commitmenttypes.QueryValidatorsRequest) ([]byte, error) {
	totalBonded, err := oq.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return nil, err
	}
	delegatorAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid delegator address")
	}

	validators, err := oq.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr, math.MaxInt16)
	if err != nil {
		return nil, err
	}

	validatorsCW := oq.BuildDelegatorValidatorsResponseCW(ctx, validators.Validators, totalBonded, query.DelegatorAddress)
	res := commitmenttypes.QueryDelegatorValidatorsResponse{
		Validators: validatorsCW,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get delegation validator response")
	}

	return responseBytes, nil
}

func (oq *Querier) BuildDelegatorValidatorsResponseCW(ctx sdk.Context, validators []stakingtypes.Validator, totalBonded sdkmath.Int, delegatorAddress string) []commitmenttypes.ValidatorDetail {
	edenDenomPrice := sdkmath.LegacyZeroDec()
	baseCurrency, found := oq.assetKeeper.GetUsdcDenom(ctx)
	if found {
		edenDenomPrice = oq.ammKeeper.GetEdenDenomPrice(ctx, baseCurrency)
	}

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
			UsdAmount: edenDenomPrice.Mul(tokens),
		}

		votingPower := sdkmath.LegacyNewDecFromInt(validator.Tokens).QuoInt(totalBonded).MulInt(sdkmath.NewInt(100))
		validatorCW.VotingPower = votingPower

		validatorsCW = append(validatorsCW, validatorCW)
	}

	return validatorsCW
}
