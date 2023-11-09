package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkquery "github.com/cosmos/cosmos-sdk/types/query"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryPriceAll(ctx sdk.Context, priceAll *wasmbindingstypes.PriceAll) ([]byte, error) {
	pagination := priceAll.Pagination

	// Calling the PriceAll function and handling its response
	priceResponse, err := oq.keeper.PriceAll(ctx, &oracletypes.QueryAllPriceRequest{Pagination: pagination})
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get all prices")
	}

	// copy array priceResponse.Price
	price := make([]oracletypes.Price, len(priceResponse.Price))
	copy(price, priceResponse.Price)

	res := wasmbindingstypes.AllPriceResponse{
		Price: price,
		Pagination: &sdkquery.PageResponse{
			NextKey: priceResponse.Pagination.NextKey,
		},
	}

	// Serializing the response to a JSON byte array
	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize price response")
	}

	return responseBytes, nil
}
