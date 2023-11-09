package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func (oq *Querier) queryAssetInfo(ctx sdk.Context, assetInfo *wasmbindingstypes.AssetInfo) ([]byte, error) {
	denom := assetInfo.Denom

	AssetInfoResp, err := oq.keeper.AssetInfo(ctx, &oracletypes.QueryGetAssetInfoRequest{Denom: denom})
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query asset info")
	}

	res := wasmbindingstypes.AssetInfoResponse{
		AssetInfo: &wasmbindingstypes.AssetInfoType{
			Denom:      AssetInfoResp.AssetInfo.Denom,
			Display:    AssetInfoResp.AssetInfo.Display,
			BandTicker: AssetInfoResp.AssetInfo.BandTicker,
			ElysTicker: AssetInfoResp.AssetInfo.ElysTicker,
			Decimal:    AssetInfoResp.AssetInfo.Decimal,
		},
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize asset info response")
	}
	return responseBytes, nil
}
