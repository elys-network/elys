package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/x/commitment/types"
)

// UpdateVestingInfo add/update specific vesting info by denom on Params
func (k msgServer) UpdateVestingInfo(goCtx context.Context, msg *types.MsgUpdateVestingInfo) (*types.MsgUpdateVestingInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	params := k.GetParams(ctx)
	vestingInfo, index := k.GetVestingInfo(ctx, msg.BaseDenom)
	if vestingInfo == nil {
		vestingInfo = &types.VestingInfo{
			BaseDenom:      msg.BaseDenom,
			VestingDenom:   msg.VestingDenom,
			NumBlocks:      msg.NumBlocks,
			VestNowFactor:  sdkmath.NewInt(msg.VestNowFactor),
			NumMaxVestings: msg.NumMaxVestings,
		}
		params.VestingInfos = append(params.VestingInfos, *vestingInfo)
	} else {
		params.VestingInfos[index].BaseDenom = msg.BaseDenom
		params.VestingInfos[index].VestingDenom = msg.VestingDenom
		params.VestingInfos[index].NumBlocks = msg.NumBlocks
		params.VestingInfos[index].VestNowFactor = sdkmath.NewInt(msg.VestNowFactor)
		params.VestingInfos[index].NumMaxVestings = msg.NumMaxVestings
	}

	err := params.Validate()
	if err != nil {
		return nil, fmt.Errorf("params validation failed: %s", err.Error())
	}

	// store params
	k.SetParams(ctx, params)
	return &types.MsgUpdateVestingInfoResponse{}, nil
}
