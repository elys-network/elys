package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	mtp, err := k.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, err
	}

	var closedMtp *types.MTP
	var repayAmount sdk.Int
	closedMtp, repayAmount, err = k.CloseLong(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(closedMtp.Id), 10)),
		sdk.NewAttribute("address", closedMtp.Address),
		sdk.NewAttribute("collateral_asset", closedMtp.CollateralAssets[0]),
		sdk.NewAttribute("collateral_amount", closedMtp.CollateralAmounts[0].String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", closedMtp.Leverages[0].String()),
		sdk.NewAttribute("liabilities", closedMtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollaterals[0].String()),
		sdk.NewAttribute("health", closedMtp.MtpHealth.String()),
	))

	return &types.MsgCloseResponse{}, nil
}
