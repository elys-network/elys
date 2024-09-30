package wasm

import (
	"encoding/json"
	"fmt"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (oq *Querier) queryStakedPositions(ctx sdk.Context, query *commitmenttypes.QueryValidatorsRequest) ([]byte, error) {
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

	stakedPositionsCW := oq.BuildStakedPositionResponseCW(ctx, validators.Validators, totalBonded, query.DelegatorAddress)
	res := commitmenttypes.QueryStakedPositionResponse{
		StakedPosition: stakedPositionsCW,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get staked position response")
	}

	return responseBytes, nil
}

func (oq *Querier) BuildStakedPositionResponseCW(ctx sdk.Context, validators []stakingtypes.Validator, totalBonded sdkmath.Int, delegatorAddress string) []commitmenttypes.StakedPosition {
	edenDenomPrice := sdkmath.LegacyZeroDec()

	baseCurrency, found := oq.assetKeeper.GetUsdcDenom(ctx)
	if found {
		edenDenomPrice = oq.ammKeeper.GetEdenDenomPrice(ctx, baseCurrency)
	}

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
		delAmount := tokens.TruncateInt()
		votingPower := sdkmath.LegacyNewDecFromInt(validator.Tokens).QuoInt(totalBonded).MulInt(sdkmath.NewInt(100))

		website := validator.Description.Website
		if len(website) < 1 {
			website = " "
		}

		stakedPosition.Validator = commitmenttypes.StakingValidator{
			Id: validator.Description.Identity,
			// The validator address.
			Address: validator.OperatorAddress,
			// The validator name.
			Name: validator.Description.Moniker,
			// Voting power percentage for this validator.
			VotingPower: votingPower,
			// Comission percentage for the validator.
			Commission: validator.GetCommission(),
		}
		stakedPosition.Staked = commitmenttypes.BalanceAvailable{
			Amount:    delAmount,
			UsdAmount: edenDenomPrice.Mul(tokens),
		}

		stakedPositions = append(stakedPositions, stakedPosition)
	}

	return stakedPositions
}
