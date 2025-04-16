package keeper

import (
	"context"
	"github.com/elys-network/elys/utils"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetConsolidatedPrice(goCtx context.Context, req *types.QueryGetConsolidatedPriceRequest) (*types.QueryGetConsolidatedPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	oracle, amm, oracleDec := k.RetrieveConsolidatedPrice(ctx, req.Denom)

	return &types.QueryGetConsolidatedPriceResponse{
		AmmPrice:       amm.Dec(),
		OraclePrice:    oracle.Dec(),
		OraclePriceDec: oracleDec.Dec(),
	}, nil
}

func (k Keeper) GetAllPrices(goCtx context.Context, req *types.QueryGetAllPricesRequest) (*types.QueryGetAllPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var prices []*types.Price

	assetentries := k.assetProfileKeeper.GetAllEntry(ctx)
	for _, assetEntry := range assetentries {
		if strings.HasPrefix(assetEntry.Denom, "amm/pool") || strings.HasPrefix(assetEntry.Denom, "stablestake") {
			continue
		}
		denom := assetEntry.Denom
		if assetEntry.Denom == ptypes.Eden {
			denom = ptypes.Elys
		}
		tokenPriceOracle := k.oracleKeeper.GetDenomPrice(ctx, denom).Mul(utils.Pow10(assetEntry.Decimals))
		tokenPriceAmm := k.amm.CalcAmmPrice(ctx, denom, assetEntry.Decimals).Mul(utils.Pow10(assetEntry.Decimals))
		prices = append(prices, &types.Price{
			Denom:       assetEntry.Denom,
			OraclePrice: tokenPriceOracle.Dec(),
			AmmPrice:    tokenPriceAmm.Dec(),
		})
	}

	return &types.QueryGetAllPricesResponse{
		Prices: prices,
	}, nil
}
