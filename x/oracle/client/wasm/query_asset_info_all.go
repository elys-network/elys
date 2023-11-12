package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryAssetInfoAll(ctx sdk.Context, req *oracletypes.QueryAllAssetInfoRequest) ([]byte, error) {
	res, err := oq.keeper.AssetInfoAll(ctx, req)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query asset info all")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize asset info all response")
	}
	return responseBytes, nil
}
