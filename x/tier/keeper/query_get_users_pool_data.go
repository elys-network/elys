package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
"pool": "ATOM/USDC",
"value": "2.013869980978454599" //USDC
"poolId"://
"fiat_value":
*/
func (k Keeper) GetUsersPoolData(goCtx context.Context, req *types.QueryGetUsersPoolDataRequest) (*types.QueryGetUsersPoolDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	usdcDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)

	if !found {
		return nil, errors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	listCommitments, pagination, err := k.commitement.GetAllCommitmentsWithPagination(ctx, req.Pagination)

	if err != nil {
		return nil, err
	}

	tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, usdcDenom)
	params := k.stablestakeKeeper.GetParams(ctx)

	usersData := []*types.UserData{}

	for _, user := range listCommitments {

		u := types.UserData{
			User:  user.Creator,
			Pools: []*types.Pool{},
		}

		for _, commitment := range user.CommittedTokens {
			if strings.HasPrefix(commitment.Denom, "stablestake/share") {
				fiatValue := commitment.Amount.ToLegacyDec().Mul(params.RedemptionRate).Mul(tokenPrice)
				u.Pools = append(u.Pools, &types.Pool{
					Pool:      "USDC",
					PoolId:    commitment.Denom,
					FiatValue: fiatValue.String(),
					Amount:    commitment.Amount,
				})
			}

			if strings.HasPrefix(commitment.Denom, "amm/pool") {

				poolId, err := GetPoolIdFromShareDenom(commitment.Denom)

				if err != nil {
					continue
				}

				pool, found := k.amm.GetPool(ctx, poolId)
				if !found {
					continue
				}

				info := k.amm.PoolExtraInfo(ctx, pool)
				fiatValue := commitment.Amount.ToLegacyDec().Mul(info.LpTokenPrice).QuoInt(ammtypes.OneShare)

				u.Pools = append(u.Pools, &types.Pool{
					Pool:      commitment.Denom,
					PoolId:    string(pool.PoolId),
					FiatValue: fiatValue.String(),
					Amount:    commitment.Amount,
				})
			}

		}

		usersData = append(usersData, &u)
	}

	return &types.QueryGetUsersPoolDataResponse{
		Users:      usersData,
		Pagination: pagination,
	}, nil
}
