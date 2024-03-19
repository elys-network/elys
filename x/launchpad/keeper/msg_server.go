package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/launchpad/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k Keeper) isEnabledToken(ctx sdk.Context, spendingToken string) bool {
	params := k.GetParams(ctx)

	for _, token := range params.SpendingTokens {
		if token == spendingToken {
			return true
		}
	}
	return false
}

func (k Keeper) CreateOrder(ctx sdk.Context, orderMaker string, spendingToken string, elysAmount math.Int, bonusRate sdk.Dec) {
	params := k.GetParams(ctx)
	order := types.Purchase{}
	order.OrderId = k.LastOrderId(ctx) + 1
	order.OrderMaker = orderMaker
	order.SpendingToken = spendingToken
	order.ElysAmount = elysAmount
	order.BonusAmount = bonusRate.MulInt(elysAmount).TruncateInt()
	order.TokenAmount = sdk.NewDecFromInt(elysAmount).Mul(sdk.NewDec(1000_000)).Quo(params.InitialPrice).RoundInt()
	order.ReturnedElysAmount = sdk.ZeroInt()

	k.SetPurchase(ctx, order)
}

func (k msgServer) BuyElys(goCtx context.Context, msg *types.MsgBuyElys) (*types.MsgBuyElysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if params.LaunchpadStarttime > uint64(ctx.BlockTime().Unix()) {
		return nil, types.ErrLaunchpadNotStarted
	}

	if params.LaunchpadStarttime+params.LaunchpadDuration < uint64(ctx.BlockTime().Unix()) {
		return nil, types.ErrLaunchpadAlreadyFinished
	}

	if !k.isEnabledToken(ctx, msg.SpendingToken) {
		return nil, types.ErrNotEnabledSpendingToken
	}

	asset, found := k.assetProfileKeeper.GetEntry(ctx, msg.SpendingToken)
	if !found {
		return nil, errorsmod.Wrapf(atypes.ErrAssetProfileNotFound, "denom: %s", msg.SpendingToken)
	}

	price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)

	elysAmount := price.MulInt(msg.TokenAmount).Mul(sdk.NewDec(1000_000)).Quo(params.InitialPrice).RoundInt()

	soldAmount := params.SoldAmount.Add(elysAmount)
	if soldAmount.GT(params.TotalReserve) {
		return nil, types.ErrOverflowTotalReserve
	}

	// 0-20% raise 100% bonus
	// 20-30% raise 90% bonus
	// 30-40% raise bonus 80%
	// 40-50% raise bonus 70%
	// 50-60% raise bonus 60%
	// 60-70% raise bonus 50%
	// 70-80% raise bonus 40%
	// 80-100% raise bonus 30%
	soldSoFar := params.SoldAmount
	for index, raisePercent := range params.BonusInfo.RaisePercents {
		roundMaxRaise := sdk.NewDecWithPrec(int64(raisePercent), 2).Mul(sdk.NewDecFromInt(params.TotalReserve)).RoundInt()
		if soldSoFar.LT(roundMaxRaise) {
			bonusPercent := params.BonusInfo.BonusPercents[index]
			bonusRate := sdk.NewDecWithPrec(int64(bonusPercent), 2)
			order := types.Purchase{
				OrderId:            0,
				OrderMaker:         msg.Sender,
				SpendingToken:      msg.SpendingToken,
				TokenAmount:        msg.TokenAmount,
				ElysAmount:         sdk.ZeroInt(),
				ReturnedElysAmount: sdk.ZeroInt(),
				BonusAmount:        sdk.ZeroInt(),
			}
			if roundMaxRaise.GTE(soldAmount) {
				roundSellAmount := soldAmount.Sub(soldSoFar)
				order.ElysAmount = roundSellAmount
				order.BonusAmount = bonusRate.MulInt(roundSellAmount).TruncateInt()
				order.TokenAmount = sdk.NewDecFromInt(roundSellAmount).Mul(sdk.NewDec(1000_000)).Quo(params.InitialPrice).RoundInt()
				k.SetPurchase(ctx, order)
				break
			} else {
				roundSellAmount := roundMaxRaise.Sub(soldSoFar)
				order.ElysAmount = roundSellAmount
				order.BonusAmount = bonusRate.MulInt(roundSellAmount).TruncateInt()
				order.TokenAmount = sdk.NewDecFromInt(roundSellAmount).Mul(sdk.NewDec(1000_000)).Quo(params.InitialPrice).RoundInt()
				k.SetPurchase(ctx, order)
				soldSoFar = roundMaxRaise
			}
		}
	}

	params.SoldAmount = soldAmount
	k.SetParams(ctx, params)

	return &types.MsgBuyElysResponse{}, nil
}

func (k msgServer) ReturnElys(goCtx context.Context, msg *types.MsgReturnElys) (*types.MsgReturnElysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx
	// TODO: implement
	// params := k.GetParams(ctx)
	// params.MaxReturnPercent
	// params.ReturnDuration

	return &types.MsgReturnElysResponse{}, nil
}

func (k msgServer) WithdrawRaised(goCtx context.Context, msg *types.MsgWithdrawRaised) (*types.MsgWithdrawRaisedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if params.WithdrawAddress != msg.Sender {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "expected %s, got %s", params.WithdrawAddress, msg.Sender)
	}

	// TODO: implement
	// params := k.GetParams(ctx)
	// params.ReturnDuration
	// params.WithdrawAddress
	// params.WithdrawnAmount

	return &types.MsgWithdrawRaisedResponse{}, nil
}
