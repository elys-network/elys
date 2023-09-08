package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, ammPool ammtypes.Pool) (sdk.Dec, error) {
	xl := mtp.Liabilities

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}
	// include unpaid interest in debt (from disabled incremental pay)
	if mtp.InterestUnpaidCollateral.GT(sdk.ZeroInt()) {
		xl = xl.Add(mtp.InterestUnpaidCollateral)
	}

	custodyTokenIn := sdk.NewCoin(mtp.CustodyAsset, mtp.CustodyAmount)
	// All liabilty is in usdc
	C, err := k.EstimateSwapGivenOut(ctx, custodyTokenIn, ptypes.USDC, ammPool)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	lr := sdk.NewDecFromBigInt(C.BigInt()).Quo(sdk.NewDecFromBigInt(xl.BigInt()))

	return lr, nil
}
