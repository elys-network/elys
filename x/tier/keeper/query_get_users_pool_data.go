package keeper

import (
	"context"
	"strconv"
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

type detailPool struct {
	Title string
	Pool  ammtypes.Pool
}

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

	tokenPrice, _ := k.oracleKeeper.GetAssetPriceFromDenom(ctx, usdcDenom)
	params := k.stablestakeKeeper.GetParams(ctx)

	usersData := []*types.UserData{}

	pools := map[uint64]detailPool{}

	for _, user := range listCommitments {

		u := types.UserData{
			User:  user.Creator,
			Pools: []*types.Pool{},
		}

		for _, commitment := range user.CommittedTokens {
			if strings.HasPrefix(commitment.Denom, "stablestake/share") {
				fiatValue := tokenPrice.MulLegacyDec(params.RedemptionRate).MulInt(commitment.Amount)
				u.Pools = append(u.Pools, &types.Pool{
					Pool:      "USDC",
					PoolId:    commitment.Denom,
					FiatValue: fiatValue.String(),
					Amount:    commitment.Amount,
				})

				continue
			}

			if strings.HasPrefix(commitment.Denom, "amm/pool") {

				poolId, err := GetPoolIdFromShareDenom(commitment.Denom)
				if err != nil {
					continue
				}

				poolTitle := ""
				var pool ammtypes.Pool
				if p, ok := pools[poolId]; ok {
					poolTitle = p.Title
					pool = p.Pool
				} else {
					pool, found = k.amm.GetPool(ctx, poolId)
					if !found {
						continue
					}

					for _, asset := range pool.PoolAssets {
						entry, found := k.assetProfileKeeper.GetEntryByDenom(ctx, asset.Token.Denom)
						if !found {
							continue
						}
						poolTitle += entry.DisplayName
					}

					pools[poolId] = detailPool{
						Title: poolTitle,
						Pool:  pool,
					}
				}

				info := k.amm.PoolExtraInfo(ctx, pool)
				fiatValue := commitment.Amount.ToLegacyDec().Mul(info.LpTokenPrice).QuoInt(ammtypes.OneShare)

				poolID := strconv.FormatUint(pool.PoolId, 10)

				u.Pools = append(u.Pools, &types.Pool{
					Pool:      poolTitle,
					PoolId:    poolID,
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
