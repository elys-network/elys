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

	if msg.Refund {
		kol.Refunded = true
		k.SetKol(ctx, kol)
		return &types.MsgClaimKolResponse{
			ElysAmount:       math.ZeroInt(),
			VestedElysAmount: math.ZeroInt(),
		}, nil
	}

	currentHeight := uint64(ctx.BlockHeight())
	if currentHeight < params.StartAirdropClaimHeight {
		return nil, types.ErrAirdropNotStarted
	}

	total_elys := kol.Amount
	// 12.5% of total_amount will go directly to claimer
	// 87.5% of total_amount will be vested
	liquid_elys := math.LegacyNewDecFromInt(total_elys).Mul(math.LegacyMustNewDecFromStr("0.125")).TruncateInt()

	if liquid_elys.IsPositive() {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, kolWallet, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, liquid_elys)))
		if err != nil {
			return nil, err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, liquid_elys)))
		if err != nil {
			return nil, err
		}
	}

	kol.Claimed = true
	k.SetKol(ctx, kol)

	return &types.MsgClaimKolResponse{
		ElysAmount:       liquid_elys,
		VestedElysAmount: total_elys.Sub(liquid_elys),
	}, nil
}
