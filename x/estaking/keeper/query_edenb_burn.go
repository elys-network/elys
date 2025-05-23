package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) EdenBBurnAmount(goCtx context.Context, req *types.QueryEdenBBurnAmountRequest) (*types.QueryEdenBBurnAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegator, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	// Get commitments
	commitments := k.commKeeper.GetCommitments(ctx, delegator)

	// Get previous amount
	prevElysStaked := k.GetElysStaked(ctx, delegator)
	if prevElysStaked.Amount.IsZero() {
		return nil, status.Error(codes.NotFound, "no staked Elys")
	}

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)
	totalEdenB := edenBCommitted.Add(edenBClaimed)

	edenBToBurn := osmomath.ZeroBigDec()

	if req.TokenType == types.TokenType_TOKEN_TYPE_ELYS {
		if req.Amount.GT(prevElysStaked.Amount) {
			return nil, status.Error(codes.InvalidArgument, "amount is greater than staked Elys")
		}
		// Unstaked
		unstakedElysDec := osmomath.BigDecFromSDKInt(req.Amount)
		edenCommittedAndElysStakedDec := osmomath.BigDecFromSDKInt(edenCommitted.Add(prevElysStaked.Amount))

		if edenCommittedAndElysStakedDec.GT(osmomath.ZeroBigDec()) {
			edenBToBurn = unstakedElysDec.Quo(edenCommittedAndElysStakedDec).Mul(osmomath.BigDecFromSDKInt(totalEdenB))
		}
		if edenCommittedAndElysStakedDec.IsZero() {
			edenBToBurn = osmomath.BigDecFromSDKInt(totalEdenB)
		}
	} else if req.TokenType == types.TokenType_TOKEN_TYPE_EDEN {
		if req.Amount.GT(edenCommitted) {
			return nil, status.Error(codes.InvalidArgument, "amount is greater than eden committed")
		}

		unclaimedAmtDec := osmomath.BigDecFromSDKInt(req.Amount)
		// This formula should be applied before eden uncommitted or elys staked is removed from eden committed amount and elys staked amount respectively
		// So add uncommitted amount to committed eden bucket in calculation.
		edenCommittedAndElysStakedDec := osmomath.BigDecFromSDKInt(edenCommitted.Add(prevElysStaked.Amount).Add(req.Amount))
		if edenCommittedAndElysStakedDec.IsZero() {
			return nil, status.Error(codes.NotFound, "no eden committed")
		}

		if edenCommittedAndElysStakedDec.GT(osmomath.ZeroBigDec()) {
			edenBToBurn = unclaimedAmtDec.Quo(edenCommittedAndElysStakedDec).Mul(osmomath.BigDecFromSDKInt(totalEdenB))
		}
		if edenCommittedAndElysStakedDec.IsZero() {
			edenBToBurn = osmomath.BigDecFromSDKInt(totalEdenB)
		}
	}

	return &types.QueryEdenBBurnAmountResponse{BurnEdenbAmount: edenBToBurn.Dec().TruncateInt()}, nil
}
