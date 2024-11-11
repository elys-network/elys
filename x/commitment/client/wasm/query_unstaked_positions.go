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

func (oq *Querier) queryUnStakedPositions(ctx sdk.Context, query *commitmenttypes.QueryValidatorsRequest) ([]byte, error) {
	totalBonded, err := oq.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return nil, err
	}
	delegatorAddr, err := sdk.AccAddressFromBech32(query.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid delegator address")
	}

	unbonding_delegations, err := oq.stakingKeeper.GetUnbondingDelegations(ctx, delegatorAddr, math.MaxInt16)
	if err != nil {
		return nil, err
	}

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

func (oq *Querier) BuildUnStakedPositionResponseCW(ctx sdk.Context, unbondingDelegations []stakingtypes.UnbondingDelegation, totalBonded sdkmath.Int, delegatorAddress string) []commitmenttypes.UnstakedPosition {
	edenDenomPrice := sdkmath.LegacyZeroDec()
	baseCurrency, found := oq.assetKeeper.GetUsdcDenom(ctx)
	if found {
		edenDenomPrice = oq.ammKeeper.GetEdenDenomPrice(ctx, baseCurrency)
	}

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
			validator, err := oq.stakingKeeper.Validator(ctx, valAddress)
			if err != nil {
				continue
			}
			val := validator.(stakingtypes.Validator)
			votingPower := sdkmath.LegacyNewDecFromInt(val.GetBondedTokens()).QuoInt(totalBonded).MulInt(sdkmath.NewInt(100))

			website := val.Description.Website
			if len(website) < 1 {
				website = " "
			}

			unstakedPosition.Validator = commitmenttypes.StakingValidator{
				Id: val.Description.Identity,
				// The validator address.
				Address: val.OperatorAddress,
				// The validator name.
				Name: val.Description.Moniker,
				// Voting power percentage for this validator.
				VotingPower: votingPower,
				// Comission percentage for the validator.
				Commission: val.GetCommission(),
			}
			unstakedPosition.RemainingTime = uint64(entity.CompletionTime.Unix())
			unstakedPosition.Unstaked = commitmenttypes.BalanceAvailable{
				Amount:    entity.Balance,
				UsdAmount: edenDenomPrice.MulInt(entity.Balance),
			}

			unstakedPositions = append(unstakedPositions, unstakedPosition)
			i++
		}
	}

	return unstakedPositions
}
