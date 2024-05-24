package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) UpdateWasmConfig(goCtx context.Context, msg *types.MsgUpdateWasmConfig) (*types.MsgUpdateWasmConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	wasmMaxLabelSize, ok := sdk.NewIntFromString(msg.WasmMaxLabelSize)

	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wasm max label size must be a positive integer")
	}

	wasmMaxProposalWasmSize, ok := sdk.NewIntFromString(msg.WasmMaxProposalWasmSize)

	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wasm max proposal wasm size must be a positive integer")
	}

	wasmMaxSize, ok := sdk.NewIntFromString(msg.WasmMaxSize)

	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wasm max size must be a positive integer")
	}

	// increase wasm size limit
	wasm.MaxLabelSize = int(wasmMaxLabelSize.Int64())
	wasm.MaxWasmSize = int(wasmMaxProposalWasmSize.Int64())
	wasm.MaxProposalWasmSize = int(wasmMaxSize.Int64())

	params := k.GetParams(ctx)
	params.WasmMaxLabelSize = wasmMaxLabelSize
	params.WasmMaxProposalWasmSize = wasmMaxProposalWasmSize
	params.WasmMaxSize = wasmMaxSize
	k.SetParams(ctx, params)
	return &types.MsgUpdateWasmConfigResponse{}, nil
}
