package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

func (k msgServer) SwapTakerFeesAndBurn(goCtx context.Context, msg *types.MsgSwapTakerFeesAndBurn) (*types.MsgSwapTakerFeesAndBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	collectionAddressString := k.parameterKeeper.GetParams(ctx).TakerFeeCollectionAddress
	// Convert balances in taker address to elys
	collectionAddress, err := sdk.AccAddressFromBech32(collectionAddressString)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "invalid taker fee collection address")
	}

	balances := k.bankKeeper.GetAllBalances(ctx, collectionAddress)
	for _, balance := range balances {
		// need at least a certain amount to swap
		if balance.Denom == ptypes.Elys || balance.Amount.LT(sdkmath.NewInt(1000000)) {
			continue
		}
		_, err = k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
			Sender:    collectionAddressString,
			Recipient: collectionAddressString,
			Amount:    sdk.NewCoin(balance.Denom, balance.Amount),
			DenomIn:   balance.Denom,
			DenomOut:  ptypes.Elys,
			MinAmount: sdk.NewCoin(ptypes.Elys, sdkmath.ZeroInt()),
		})
		if err != nil {
			return nil, err
		}
	}

	elysBalance := k.bankKeeper.GetBalance(ctx, collectionAddress, ptypes.Elys)
	if elysBalance.IsPositive() {
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, collectionAddress, types.ModuleName, sdk.NewCoins(elysBalance))
		if err != nil {
			return nil, err
		}

		// burn elys token
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(elysBalance))
		if err != nil {
			return nil, err
		}
		// event for burning taker fees
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.TypeEvtTakerFeeBurn,
				sdk.NewAttribute("amount", elysBalance.String()),
			),
		})
	}

	return &types.MsgSwapTakerFeesAndBurnResponse{
		ElysBurnt: elysBalance.Amount,
	}, nil
}
