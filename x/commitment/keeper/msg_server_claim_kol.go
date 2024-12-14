package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

var KolWallet = "elys1wk7jwkqt2h9cnpkst85j9n454e4y8znlgk842n"

func (k msgServer) ClaimKol(goCtx context.Context, msg *types.MsgClaimKol) (*types.MsgClaimKolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.ClaimAddress)
	kolWallet := sdk.MustAccAddressFromBech32(KolWallet)
	params := k.GetParams(ctx)

	kol := k.GetKol(ctx, sender)
	if kol.Amount.IsZero() {
		return nil, types.ErrKolNotFound
	}

	if kol.Claimed {
		return nil, types.ErrKolAlreadyClaimed
	}

	if kol.Refunded {
		return nil, types.ErrKolRefunded
	}

	//TODO: If price touches 1, set permanent block on refunds

	currentHeight := uint64(ctx.BlockHeight())
	if currentHeight < params.StartAirdropClaimHeight {
		return nil, types.ErrAirdropNotStarted
	}

	total_elys := kol.Amount

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, kolWallet, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, total_elys)))
	if err != nil {
		return nil, err
	}

	// 15% of total_amount will go directly to claimer
	// 85% of total_amount will be vested
	liquid_elys := math.LegacyNewDecFromInt(total_elys).Mul(math.LegacyMustNewDecFromStr("0.15")).TruncateInt()
	// TODO: Handle vesting
	//vested_elys := total_elys.Sub(liquid_elys)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, liquid_elys)))
	if err != nil {
		return nil, err
	}

	// Create a new vesting schedule

	kol.Claimed = true
	k.SetKol(ctx, kol)

	return &types.MsgClaimKolResponse{}, nil
}
