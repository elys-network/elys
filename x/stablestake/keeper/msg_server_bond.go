package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Bond(goCtx context.Context, msg *types.MsgBond) (*types.MsgBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	sender := sdk.MustAccAddressFromBech32(msg.Creator)

	entry, found := k.assetProfileKeeper.GetEntry(ctx, params.DepositDenom)
	if !found {
		return nil, errors.New("invalid denom")
	}

	depositDenom := entry.Denom
	depositCoin := sdk.NewCoin(depositDenom, msg.Amount)
	err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{depositCoin})
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

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, shareCoins)
	if err != nil {
		return nil, err
	}

	entry, found = k.assetProfileKeeper.GetEntry(ctx, shareDenom)
	if !found {
		// Set an entity to assetprofile
		entry = assetprofiletypes.Entry{
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

	// Create a commit LP token message
	msgLiquidCommitLPToken := &ctypes.MsgCommitLiquidTokens{
		Creator:   sender.String(),
		Denom:     shareDenom,
		Amount:    shareAmount,
		LockUntil: uint64(ctx.BlockTime().Unix()),
	}

	// Commit LP token
	msgServer := commitmentkeeper.NewMsgServerImpl(*k.commitmentKeeper)
	_, err = msgServer.CommitLiquidTokens(sdk.WrapSDKContext(ctx), msgLiquidCommitLPToken)
	if err != nil {
		return nil, err
	}

	params.TotalValue = params.TotalValue.Add(msg.Amount)
	k.SetParams(ctx, params)

	return &types.MsgBondResponse{}, nil
}
