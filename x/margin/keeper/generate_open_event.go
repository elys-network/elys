package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) GenerateOpenEvent(mtp *types.MTP) sdk.Event {
	collateralIndex := len(mtp.Collaterals) - 1
	custodyIndex := len(mtp.Custodies) - 1
	mtpPosIndex := len(mtp.Leverages) - 1

	return sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral", mtp.Collaterals[collateralIndex].String()),
		sdk.NewAttribute("custody", mtp.Custodies[custodyIndex].String()),
		sdk.NewAttribute("leverage", mtp.Leverages[mtpPosIndex].String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollaterals[collateralIndex].String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustodies[custodyIndex].String()),
		sdk.NewAttribute("interest_unpaid_collateral", mtp.InterestUnpaidCollaterals[collateralIndex].String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	)
}
