package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

var KolWallet = "elys1wk7jwkqt2h9cnpkst85j9n454e4y8znlgk842n"

func (k msgServer) ClaimKol(goCtx context.Context, msg *types.MsgClaimKol) (*types.MsgClaimKolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.ClaimAddress)
	airdropWallet := sdk.MustAccAddressFromBech32(AirdropWallet)
	params := k.GetParams(ctx)

	// if k.GetAirdropClaimed(ctx, sender).Claimed {
	// 	return nil, types.ErrAirdropAlreadyClaimed
	// }

	currentHeight := uint64(ctx.BlockHeight())
	if currentHeight < params.StartAirdropClaimHeight {
		return nil, types.ErrAirdropNotStarted
	}

	// if currentHeight > params.EndAirdropClaimHeight {
	// 	return nil, types.ErrAirdropEnded
	// }

	// Elys to be sent
	total_elys := k.GetNFTHolder(ctx, sender).Amount
	total_elys = total_elys.Add(k.GetCadet(ctx, sender).Amount)
	total_elys = total_elys.Add(k.GetGovernor(ctx, sender).Amount)

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, airdropWallet, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, total_elys)))
	if err != nil {
		return nil, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, total_elys)))
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimKolResponse{}, nil
}
