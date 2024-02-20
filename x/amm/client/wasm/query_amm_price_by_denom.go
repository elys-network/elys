package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (oq *Querier) queryAmmPriceByDenom(ctx sdk.Context, query *ammtypes.QueryAMMPriceRequest) ([]byte, error) {
	denom := ptypes.BaseCurrency
	assetProfile, found := oq.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	// If amount is zero
	if query.TokenIn.Amount.IsZero() {
		responseBytes, err := json.Marshal(sdk.ZeroDec())
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to serialize in route by denom response")
		}
		return responseBytes, nil
	}

	// uses asset profile denom
	usdcDenom := assetProfile.Denom

	resp, err := oq.keeper.InRouteByDenom(sdk.WrapSDKContext(ctx), &ammtypes.QueryInRouteByDenomRequest{DenomIn: query.TokenIn.Denom, DenomOut: usdcDenom})
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get in route by denom")
	}

	routes := resp.InRoute
	tokenIn := query.TokenIn
	discount := query.Discount

	spotPrice, _, _, _, _, _, _, err := oq.keeper.CalcInRouteSpotPrice(ctx, tokenIn, routes, discount, sdk.ZeroDec())
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get in route by denom")
	}

	responseBytes, err := json.Marshal(spotPrice)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize in route by denom response")
	}
	return responseBytes, nil
}
