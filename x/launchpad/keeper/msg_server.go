package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/launchpad/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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

func (k Keeper) GenerateOrder(ctx sdk.Context, orderMaker string, spendingToken string, elysAmount math.Int, bonusRate sdk.Dec, price sdk.Dec) types.Purchase {
	params := k.GetParams(ctx)
	order := types.Purchase{}
	order.OrderId = k.LastOrderId(ctx) + 1
	order.OrderMaker = orderMaker
	order.SpendingToken = spendingToken
	order.ElysAmount = elysAmount
	order.BonusAmount = bonusRate.MulInt(elysAmount).TruncateInt()
	order.TokenAmount = sdk.NewDecFromInt(elysAmount).Mul(params.InitialPrice).Quo(sdk.NewDec(1000_000)).Quo(price).RoundInt()
	order.ReturnedElysAmount = sdk.ZeroInt()

	return order
}

func (k Keeper) CalcBuyElysResult(ctx sdk.Context, sender string, spendingToken string, tokenAmount math.Int) (math.Int, []types.Purchase, error) {
	params := k.GetParams(ctx)
	asset, found := k.assetProfileKeeper.GetEntry(ctx, spendingToken)
	if !found {
		return sdk.ZeroInt(), []types.Purchase{}, errorsmod.Wrapf(atypes.ErrAssetProfileNotFound, "denom: %s", spendingToken)
	}

	price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)
	if price.IsZero() {
		return sdk.ZeroInt(), []types.Purchase{}, oracletypes.ErrPriceNotSet
	}

	elysAmount := price.MulInt(tokenAmount).Mul(sdk.NewDec(1000_000)).Quo(params.InitialPrice).RoundInt()

	soldAmount := params.SoldAmount.Add(elysAmount)
	if soldAmount.GT(params.TotalReserve) {
		return sdk.ZeroInt(), []types.Purchase{}, types.ErrOverflowTotalReserve
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
	purchases := []types.Purchase{}
	for index, raisePercent := range params.BonusInfo.RaisePercents {
		roundMaxRaise := sdk.NewDecWithPrec(int64(raisePercent), 2).Mul(sdk.NewDecFromInt(params.TotalReserve)).RoundInt()
		if soldSoFar.LT(roundMaxRaise) {
			bonusPercent := params.BonusInfo.BonusPercents[index]
			bonusRate := sdk.NewDecWithPrec(int64(bonusPercent), 2)
			if roundMaxRaise.GTE(soldAmount) {
				roundSellAmount := soldAmount.Sub(soldSoFar)
				order := k.GenerateOrder(ctx, sender, spendingToken, roundSellAmount, bonusRate, price)
				purchases = append(purchases, order)
				break
			} else {
				roundSellAmount := roundMaxRaise.Sub(soldSoFar)
				order := k.GenerateOrder(ctx, sender, spendingToken, roundSellAmount, bonusRate, price)
				purchases = append(purchases, order)
				soldSoFar = roundMaxRaise
			}
		}
	}
	return elysAmount, purchases, nil
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

	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	spendingCoins := sdk.Coins{sdk.NewCoin(asset.Denom, msg.TokenAmount)}
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, spendingCoins)
	if err != nil {
		return nil, err
	}

	elysAmount, orders, err := k.CalcBuyElysResult(ctx, msg.Sender, msg.SpendingToken, msg.TokenAmount)
	if err != nil {
		return nil, err
	}

	// TODO: consider existing order on same bonus period - all should be combined into one
	for _, order := range orders {
		k.SetOrder(ctx, order)
	}

	params.SoldAmount = params.SoldAmount.Add(elysAmount)
	k.SetParams(ctx, params)

	elysCoins := sdk.Coins{sdk.NewCoin(ptypes.Elys, elysAmount)}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, elysCoins)
	if err != nil {
		return nil, err
	}

	return &types.MsgBuyElysResponse{}, nil
}

func (k Keeper) CalcReturnElysResult(ctx sdk.Context, orderId uint64, returnElysAmount math.Int) (math.Int, error) {
	params := k.GetParams(ctx)
	order := k.GetOrder(ctx, orderId)
	if order.OrderId == 0 {
		return sdk.ZeroInt(), types.ErrPurchaseOrderNotFound
	}

	maxReturnAmount := sdk.NewDecWithPrec(int64(params.MaxReturnPercent), 2).MulInt(order.ElysAmount).RoundInt()
	if order.ReturnedElysAmount.Add(returnElysAmount).GT(maxReturnAmount) {
		return sdk.ZeroInt(), types.ErrExceedMaxReturnAmount
	}

	returnTokenAmount := returnElysAmount.Mul(order.TokenAmount).Quo(order.ElysAmount)
	return returnTokenAmount, nil
}

func (k msgServer) ReturnElys(goCtx context.Context, msg *types.MsgReturnElys) (*types.MsgReturnElysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	launchpadEndTime := params.LaunchpadStarttime + params.LaunchpadDuration
	if launchpadEndTime > uint64(ctx.BlockTime().Unix()) {
		return nil, types.ErrLaunchpadNotFinished
	}

	if params.LaunchpadStarttime+params.ReturnDuration < uint64(ctx.BlockTime().Unix()) {
		return nil, types.ErrLaunchpadReturnPeriodFinished
	}

	order := k.GetOrder(ctx, msg.OrderId)
	if order.OrderId == 0 {
		return nil, types.ErrPurchaseOrderNotFound
	}

	if order.OrderMaker != msg.Sender {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "expected %s, got %s", order.OrderMaker, msg.Sender)
	}

	coins := sdk.Coins{sdk.NewCoin(ptypes.Elys, msg.ReturnElysAmount)}
	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	returnTokenAmount, err := k.CalcReturnElysResult(ctx, msg.OrderId, msg.ReturnElysAmount)
	if err != nil {
		return nil, err
	}

	asset, found := k.assetProfileKeeper.GetEntry(ctx, order.SpendingToken)
	if !found {
		return nil, errorsmod.Wrapf(atypes.ErrAssetProfileNotFound, "denom: %s", order.SpendingToken)
	}

	coins = sdk.Coins{sdk.NewCoin(asset.Denom, returnTokenAmount)}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgReturnElysResponse{}, nil
}

func (k msgServer) DepositElysToken(goCtx context.Context, msg *types.MsgDepositElysToken) (*types.MsgDepositElysTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{msg.Coin})
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositElysTokenResponse{}, nil
}

func (k msgServer) WithdrawRaised(goCtx context.Context, msg *types.MsgWithdrawRaised) (*types.MsgWithdrawRaisedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if params.WithdrawAddress != msg.Sender {
		return nil, errorsmod.Wrapf(types.ErrInvalidWithrawAccount, "expected %s, got %s", params.WithdrawAddress, msg.Sender)
	}

	// Once return period is over, can withdraw all
	returnEndTime := params.LaunchpadStarttime + params.LaunchpadDuration + params.ReturnDuration
	if returnEndTime < uint64(ctx.BlockTime().Unix()) {
		addr := sdk.MustAccAddressFromBech32(msg.Sender)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins(msg.Coins))
		if err != nil {
			return nil, err
		}
	}

	totalWithdrawAmount := sdk.Coins(params.WithdrawnAmount).Add(msg.Coins...)
	maxWithdrawAmount := sdk.Coins{}
	orders := k.GetAllOrders(ctx)
	for _, order := range orders {
		orderWithdrawable := sdk.NewDecWithPrec(int64(100-params.MaxReturnPercent), 2).MulInt(order.TokenAmount).RoundInt()
		fmt.Println("order.TokenAmount", order.TokenAmount.String())
		fmt.Println("sdk.NewDecWithPrec(int64(100-params.MaxReturnPercent), 2)", sdk.NewDecWithPrec(int64(100-params.MaxReturnPercent), 2).String())
		fmt.Println("orderWithdrawable", orderWithdrawable.String())

		asset, found := k.assetProfileKeeper.GetEntry(ctx, order.SpendingToken)
		if !found {
			return nil, errorsmod.Wrapf(atypes.ErrAssetProfileNotFound, "denom: %s", order.SpendingToken)
		}
		coin := sdk.NewCoin(asset.Denom, orderWithdrawable)
		maxWithdrawAmount = maxWithdrawAmount.Add(coin)
	}

	fmt.Println("maxWithdrawAmount", maxWithdrawAmount.String())
	fmt.Println("totalWithdrawAmount", totalWithdrawAmount.String())
	fmt.Println("params.WithdrawnAmount", sdk.Coins(params.WithdrawnAmount).String())
	if totalWithdrawAmount.IsAnyGT(maxWithdrawAmount) {
		return nil, types.ErrExceedMaxWithdrawableAmount
	}

	coins := sdk.Coins(msg.Coins)
	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	if err != nil {
		return nil, err
	}

	params.WithdrawnAmount = sdk.Coins(params.WithdrawnAmount).Add(msg.Coins...)
	k.SetParams(ctx, params)

	return &types.MsgWithdrawRaisedResponse{}, nil
}

// Update params through gov proposal
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// store params
	k.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
