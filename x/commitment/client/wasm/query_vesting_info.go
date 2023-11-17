package wasm

import (
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryVestingInfo(ctx sdk.Context, query *commitmenttypes.QueryVestingInfoRequest) ([]byte, error) {
	addr := query.Address
	denom := paramtypes.Eden

	commitment := oq.keeper.GetCommitments(ctx, addr)
	vestingTokens := commitment.GetVestingTokens()

	totalVesting := sdk.ZeroInt()
	vestingDetails := make([]commitmenttypes.VestingDetail, 0)
	for i, vesting := range vestingTokens {
		// we only support Eden vesting, so manually put the denom here.
		if vesting.Denom != denom {
			continue
		}

		vested := vesting.TotalAmount.Sub(vesting.UnvestedAmount)
		vestingDetail := commitmenttypes.VestingDetail{
			Id: fmt.Sprintf("%d", i),
			// The total vest for the current vest
			TotalVest: commitmenttypes.BalanceAvailable{
				Amount:    vesting.TotalAmount,
				UsdAmount: sdk.NewDecFromInt(vesting.TotalAmount),
			},
			// The balance that's already vested
			BalanceVested: commitmenttypes.BalanceAvailable{
				Amount:    vested,
				UsdAmount: sdk.NewDecFromInt(vested),
			},
			// The remaining amount to vest
			RemainingVest: commitmenttypes.BalanceAvailable{
				Amount:    vesting.UnvestedAmount,
				UsdAmount: sdk.NewDecFromInt(vesting.UnvestedAmount),
			},
			// Remaining time to vest. Javascript timestamp.
			// ToDo: should convert this into proper javascript timestamp
			RemainingTime: vesting.CurrentEpoch,
		}

		vestingDetails = append(vestingDetails, vestingDetail)
		totalVesting = totalVesting.Add(vesting.TotalAmount)
	}

	res := commitmenttypes.QueryVestingInfoResponse{
		Vesting: commitmenttypes.BalanceAvailable{
			Amount:    totalVesting,
			UsdAmount: sdk.NewDecFromInt(totalVesting),
		},
		VestingDetails: vestingDetails,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get reward balance response")
	}
	return responseBytes, nil
}
