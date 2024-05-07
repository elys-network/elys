package wasm

import (
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryVestingInfo(ctx sdk.Context, query *commitmenttypes.QueryVestingInfoRequest) ([]byte, error) {
	addr := query.Address

	commitment := oq.keeper.GetCommitments(ctx, addr)
	vestingTokens := commitment.GetVestingTokens()

	baseCurrency, found := oq.assetKeeper.GetUsdcDenom(ctx)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	edenDenomPrice := oq.ammKeeper.GetEdenDenomPrice(ctx, baseCurrency)

	totalVesting := sdk.ZeroInt()
	vestingDetails := make([]commitmenttypes.VestingDetail, 0)
	for i, vesting := range vestingTokens {
		vestingDetail := commitmenttypes.VestingDetail{
			Id: fmt.Sprintf("%d", i),
			TotalVesting: commitmenttypes.BalanceAvailable{
				Amount:    vesting.TotalAmount,
				UsdAmount: edenDenomPrice.MulInt(vesting.TotalAmount),
			},
			Claimed: commitmenttypes.BalanceAvailable{
				Amount:    vesting.ClaimedAmount,
				UsdAmount: edenDenomPrice.MulInt(vesting.ClaimedAmount),
			},
			VestedSoFar: commitmenttypes.BalanceAvailable{
				Amount:    vesting.VestedSoFar(ctx),
				UsdAmount: edenDenomPrice.MulInt(vesting.VestedSoFar(ctx)),
			},
			RemainingBlocks: vesting.NumBlocks - (ctx.BlockHeight() - vesting.StartBlock),
		}

		vestingDetails = append(vestingDetails, vestingDetail)
		totalVesting = totalVesting.Add(vesting.TotalAmount.Sub(vesting.ClaimedAmount))
	}

	res := commitmenttypes.QueryVestingInfoResponse{
		Vesting: commitmenttypes.BalanceAvailable{
			Amount:    totalVesting,
			UsdAmount: edenDenomPrice.MulInt(totalVesting),
		},
		VestingDetails: vestingDetails,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get reward balance response")
	}
	return responseBytes, nil
}
