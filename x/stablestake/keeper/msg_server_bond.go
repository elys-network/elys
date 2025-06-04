package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	"github.com/elys-network/elys/v6/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k msgServer) Bond(goCtx context.Context, msg *types.MsgBond) (*types.MsgBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}
	updatedNetAmount := pool.NetAmount.Add(msg.Amount)
	maxBondable := k.GetMaxBondableAmount(ctx, pool.DepositDenom)
	if updatedNetAmount.GT(maxBondable) {
		return nil, fmt.Errorf("vault cannot have more than the max bondable amount %s (current %s)", maxBondable.String(), pool.NetAmount.String())
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	redemptionRate := k.CalculateRedemptionRateForPool(ctx, pool)

	depositDenom := pool.GetDepositDenom()
	depositCoin := sdk.NewCoin(depositDenom, msg.Amount)
	err := k.bk.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.Coins{depositCoin})
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetShareDenomForPool(pool.Id)
	// Initial case
	if redemptionRate.IsZero() {
		redemptionRate = osmomath.OneBigDec()
	}
	shareAmount := osmomath.BigDecFromSDKInt(depositCoin.Amount).Quo(redemptionRate).Dec().RoundInt()
	shareCoins := sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

	err = k.bk.MintCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, shareCoins)
	if err != nil {
		return nil, err
	}

	_, found = k.assetProfileKeeper.GetEntry(ctx, shareDenom)
	if !found {
		depositDenomProfile, found := k.assetProfileKeeper.GetEntryByDenom(ctx, depositDenom)
		if !found {
			return nil, fmt.Errorf("deposit denom (%s) profile not found", depositCoin)
		}
		// Set an entity to assetprofile
		entry := assetprofiletypes.Entry{
			Authority:                authtypes.NewModuleAddress(types.ModuleName).String(),
			BaseDenom:                shareDenom,
			Decimals:                 depositDenomProfile.Decimals,
			Denom:                    shareDenom,
			Path:                     "",
			IbcChannelId:             "",
			IbcCounterpartyChannelId: "",
			DisplayName:              shareDenom,
			DisplaySymbol:            "",
			Network:                  "",
			Address:                  "",
			ExternalSymbol:           "",
			TransferLimit:            "",
			Permissions:              make([]string, 0),
			UnitDenom:                "",
			IbcCounterpartyDenom:     "",
			IbcCounterpartyChainId:   "",
			CommitEnabled:            true,
			WithdrawEnabled:          true,
		}

		k.assetProfileKeeper.SetEntry(ctx, entry)
	}

	// Commit LP token
	err = k.commitmentKeeper.CommitLiquidTokens(ctx, creator, shareDenom, shareAmount, 0)
	if err != nil {
		return nil, err
	}

	pool.NetAmount = updatedNetAmount
	k.SetPool(ctx, pool)

	if k.hooks != nil {
		err = k.hooks.AfterBond(ctx, creator, shareAmount, pool.Id)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventBond,
		sdk.NewAttribute("address", msg.Creator),
		sdk.NewAttribute("amount", msg.Amount.String()),
		sdk.NewAttribute("shares", sdk.NewCoin(shareDenom, shareAmount).String()),
	))

	return &types.MsgBondResponse{}, nil
}
