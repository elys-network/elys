package wasm

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (oq *Querier) queryBorrowRatio(ctx sdk.Context, query *types.QueryBorrowRatioRequest) ([]byte, error) {
	res, err := oq.keeper.BorrowRatio(sdk.WrapSDKContext(ctx), query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get borrow ratio")
	}

	resp := types.BalanceBorrowed{
		UsdAmount:  sdkmath.LegacyNewDecFromInt(res.TotalBorrow),
		Percentage: res.BorrowRatio,
	}

	responseBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize params response")
	}
	return responseBytes, nil
}
