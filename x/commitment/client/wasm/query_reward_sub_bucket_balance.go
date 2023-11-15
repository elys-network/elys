package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryRewardSubBucketBalanceOfDenom(ctx sdk.Context, query *commitmenttypes.QuerySubBucketBalanceRequest) ([]byte, error) {
	denom := query.Denom
	addr := query.Address
	if denom == paramtypes.Elys {
		return nil, errorsmod.Wrap(nil, "invalid denom")
	}

	var balance sdk.Coin
	commitment := oq.keeper.GetCommitments(ctx, addr)
	uncommittedToken := commitment.GetRewardUnclaimedForDenom(denom)
	switch query.Program {
	case int64(commitmenttypes.EarnType_USDC_PROGRAM):
		uncommittedToken = commitment.GetUsdcSubBucketRewardUnclaimedForDenom(denom)
	case int64(commitmenttypes.EarnType_ELYS_PROGRAM):
		uncommittedToken = commitment.GetElysSubBucketRewardUnclaimedForDenom(denom)
	case int64(commitmenttypes.EarnType_EDEN_PROGRAM):
		uncommittedToken = commitment.GetEdenSubBucketRewardUnclaimedForDenom(denom)
	case int64(commitmenttypes.EarnType_EDENB_PROGRAM):
		uncommittedToken = commitment.GetEdenBSubBucketRewardUnclaimedForDenom(denom)
	}

	balance = sdk.NewCoin(denom, uncommittedToken)

	res := commitmenttypes.BalanceAvailable{
		Amount:    balance.Amount.Uint64(),
		UsdAmount: sdk.NewDecFromInt(balance.Amount),
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get reward balance response")
	}
	return responseBytes, nil
}
