package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/tokenomics/types"
)

func (k msgServer) UpdateGenesisInflation(goCtx context.Context, msg *types.MsgUpdateGenesisInflation) (*types.MsgUpdateGenesisInflationResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	genesisInflation := types.GenesisInflation{
		Authority:             msg.Authority,
		Inflation:             msg.Inflation,
		SeedVesting:           msg.SeedVesting,
		StrategicSalesVesting: msg.StrategicSalesVesting,
	}

	k.SetGenesisInflation(ctx, genesisInflation)

	return &types.MsgUpdateGenesisInflationResponse{}, nil
}
