package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/estaking/keeper"
	"github.com/elys-network/elys/x/estaking/types"
)

// Querier handles queries for the Estaking module.
type Querier struct {
	keeper *keeper.Keeper
}

func NewQuerier(keeper *keeper.Keeper) *Querier {
	return &Querier{
		keeper: keeper,
	}
}

func (oq *Querier) HandleQuery(ctx sdk.Context, query wasmbindingstypes.ElysQuery) ([]byte, error) {
	switch {
	case query.EstakingParams != nil:
		return oq.queryParams(ctx, query.EstakingParams)
	case query.EstakingRewards != nil:
		return oq.queryRewards(ctx, query.EstakingRewards)
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

func (oq *Querier) queryRewards(ctx sdk.Context, query *types.QueryRewardsRequest) ([]byte, error) {
	res, err := oq.keeper.Rewards(ctx, query)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get estaking rewards")
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize estaking rewards response")
	}
	return responseBytes, nil
}
