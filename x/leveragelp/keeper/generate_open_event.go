package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) GenerateOpenEvent(mtp *types.MTP) sdk.Event {
	collateralIndex := len(mtp.CollateralAssets) - 1
	mtpPosIndex := len(mtp.Leverages) - 1

	return sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAssets[collateralIndex]),
		sdk.NewAttribute("collateral_amount", mtp.CollateralAmounts[collateralIndex].String()),
		sdk.NewAttribute("leverage", mtp.Leverages[mtpPosIndex].String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollaterals[collateralIndex].String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	)
}
