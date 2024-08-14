package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Bond(goCtx context.Context, msg *types.MsgBond) (*types.MsgBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	creator := sdk.MustAccAddressFromBech32(msg.Creator)

	depositDenom := k.GetDepositDenom(ctx)
	depositCoin := sdk.NewCoin(depositDenom, msg.Amount)
	err := k.bk.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.Coins{depositCoin})
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetShareDenom()
	if params.RedemptionRate.IsZero() {
		return nil, types.ErrRedemptionRateIsZero
	}
	shareAmount := sdk.NewDecFromInt(depositCoin.Amount).Quo(params.RedemptionRate).RoundInt()
	shareCoins := sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

	err = k.bk.MintCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, shareCoins)
	if err != nil {
		return nil, err
	}

	_, found := k.assetProfileKeeper.GetEntry(ctx, shareDenom)
	if !found {
		// Set an entity to assetprofile
		entry := assetprofiletypes.Entry{
			Authority:                authtypes.NewModuleAddress(types.ModuleName).String(),
			BaseDenom:                shareDenom,
			Decimals:                 ptypes.BASE_DECIMAL,
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

	params.TotalValue = params.TotalValue.Add(msg.Amount)
	k.SetParams(ctx, params)

	if k.hooks != nil {
		err = k.hooks.AfterBond(ctx, creator, shareAmount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgBondResponse{}, nil
}
