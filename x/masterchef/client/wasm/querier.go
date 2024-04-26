package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/masterchef/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
)

// Querier handles queries for the Masterchef module.
type Querier struct {
	keeper        *keeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
}

func NewQuerier(keeper *keeper.Keeper, stakingKeeper *stakingkeeper.Keeper) *Querier {
	return &Querier{
		keeper:        keeper,
		stakingKeeper: stakingKeeper,
	}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.MasterchefParams != nil:
		return oq.queryParams(ctx, query.MasterchefParams)
	case query.MasterchefExternalIncentive != nil:
		return oq.queryExternalIncentive(ctx, query.MasterchefExternalIncentive)
	case query.MasterchefPoolInfo != nil:
		return oq.queryPoolInfo(ctx, query.MasterchefPoolInfo)
	case query.MasterchefPoolRewardInfo != nil:
		return oq.queryPoolRewardInfo(ctx, query.MasterchefPoolRewardInfo)
	case query.MasterchefUserRewardInfo != nil:
		return oq.queryUserRewardInfo(ctx, query.MasterchefUserRewardInfo)
	case query.MasterchefUserPendingReward != nil:
		return oq.queryUserPendingReward(ctx, query.MasterchefUserPendingReward)
	case query.MasterchefStableStakeApr != nil:
		return oq.queryStableStakeApr(ctx, query.MasterchefStableStakeApr)
	case query.MasterchefPoolAprs != nil:
		return oq.queryPoolAprs(ctx, query.MasterchefPoolAprs)
	default:
		// This handler cannot handle the query
		return nil, wasmbindingstypes.ErrCannotHandleQuery
	}
}

func (oq *Querier) queryParams(ctx sdk.Context, query *types.QueryParamsRequest) ([]byte, error) {
	res, err := oq.keeper.Params(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get params")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize params response")
	}
	return responseBytes, nil
}

func (oq *Querier) queryExternalIncentive(ctx sdk.Context, query *types.QueryExternalIncentiveRequest) ([]byte, error) {
	res, err := oq.keeper.ExternalIncentive(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get external incentive")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize external incentive response")
	}
	return responseBytes, nil
}

func (oq *Querier) queryPoolInfo(ctx sdk.Context, query *types.QueryPoolInfoRequest) ([]byte, error) {
	res, err := oq.keeper.PoolInfo(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get pool info")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize pool info response")
	}
	return responseBytes, nil
}

func (oq *Querier) queryPoolRewardInfo(ctx sdk.Context, query *types.QueryPoolRewardInfoRequest) ([]byte, error) {
	res, err := oq.keeper.PoolRewardInfo(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get pool reward info")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize pool reward info response")
	}
	return responseBytes, nil
}

func (oq *Querier) queryUserRewardInfo(ctx sdk.Context, query *types.QueryUserRewardInfoRequest) ([]byte, error) {
	res, err := oq.keeper.UserRewardInfo(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get user reward info")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize user reward info response")
	}
	return responseBytes, nil
}

func (oq *Querier) queryUserPendingReward(ctx sdk.Context, query *types.QueryUserPendingRewardRequest) ([]byte, error) {
	res, err := oq.keeper.UserPendingReward(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get user pending reward")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize user pending reward response")
	}
	return responseBytes, nil
}

func (oq *Querier) queryStableStakeApr(ctx sdk.Context, query *types.QueryStableStakeAprRequest) ([]byte, error) {
	res, err := oq.keeper.StableStakeApr(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get user pending reward")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize user pending reward response")
	}
	return responseBytes, nil
}

func (oq *Querier) queryPoolAprs(ctx sdk.Context, query *types.QueryPoolAprsRequest) ([]byte, error) {
	res, err := oq.keeper.PoolAprs(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get user pending reward")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize user pending reward response")
	}
	return responseBytes, nil
}
