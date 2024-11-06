package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
			User: user.Creator,
			Pool: []*types.Pool{},
		}

		for _, commitment := range user.CommittedTokens {
			if strings.HasPrefix(commitment.Denom, "stablestake/share") {
				usdValue := commitment.Amount.ToLegacyDec().Mul(params.RedemptionRate).Mul(tokenPrice)
				u.Pool = append(u.Pool, &types.Pool{
					Pool:  "USDC",
					Value: usdValue.String(),
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
