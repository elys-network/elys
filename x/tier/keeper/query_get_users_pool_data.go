package keeper

import (
	"context"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	stablestaketypes "github.com/elys-network/elys/v7/x/stablestake/types"
	"github.com/elys-network/elys/v7/x/tier/types"
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

	listCommitments, pagination, err := k.commitement.GetAllCommitmentsWithPagination(ctx, req.Pagination)

	if err != nil {
		return nil, err
	}

	usersData := []*types.UserData{}

	pools := map[uint64]detailPool{}

	for _, user := range listCommitments {

		u := types.UserData{
			User:  user.Creator,
			Pools: []*types.Pool{},
		}

		for _, commitment := range user.CommittedTokens {
			if strings.HasPrefix(commitment.Denom, "stablestake/share") {
				stableId, err := stablestaketypes.GetPoolIDFromPath(commitment.Denom)
				if err != nil {
					continue
				}
				borrowPool, found := k.stablestakeKeeper.GetPool(ctx, stableId)
				if !found {
					continue
				}
				redemptionRate := k.stablestakeKeeper.CalculateRedemptionRateForPool(ctx, borrowPool)
				tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, borrowPool.GetDepositDenom())
				fiatValue := commitment.GetBigDecAmount().Mul(redemptionRate).Mul(tokenPrice)

				u.Pools = append(u.Pools, &types.Pool{
					Pool:      borrowPool.DepositDenom,
					PoolId:    strconv.FormatUint(borrowPool.Id, 10),
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
				var found bool
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

				info := k.amm.PoolExtraInfo(ctx, pool, types.OneDay)
				fiatValue := commitment.GetBigDecAmount().Mul(info.GetBigDecLpTokenPrice()).Quo(ammtypes.OneShareBigDec)

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
