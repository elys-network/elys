package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	closedMtp, repayAmount, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, initialCollateral, initialCustody, initialLiabilities, err := k.ClosePosition(ctx, msg)
	if err != nil {
		return nil, err
	}

	err = k.EmitClose(ctx, "tx", closedMtp, repayAmount, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, initialCollateral, initialCustody, initialLiabilities)
	if err != nil {
		return nil, err
	}

	return &types.MsgCloseResponse{
		Id:     closedMtp.Id,
		Amount: repayAmount,
	}, nil
}
