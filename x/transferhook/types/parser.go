package types

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

type RawPacketMetadata struct {
	Transferhook *struct {
		Receiver string             `json:"receiver"`
		Amm      *AmmPacketMetadata `json:"amm,omitempty"`
	} `json:"transferhook"`
}

type PacketForwardMetadata struct {
	Receiver    string
	RoutingInfo ModuleRoutingInfo
}

type ModuleRoutingInfo interface {
	Validate() error
}

// Packet metadata info specific to Amm (e.g. 1-click swap)
type AmmPacketMetadata struct {
	Action string                       `json:"action"`
	Routes []ammtypes.SwapAmountInRoute `json:"routes"`
}

// Validate amm packet metadata fields
// including the elys address and action type
func (m AmmPacketMetadata) Validate() error {
	if m.Action != "Swap" {
		return errorsmod.Wrapf(ErrUnsupportedAmmAction, "action %s is not supported", m.Action)
	}

	return nil
}

// Parse packet metadata intended for transferhook
// In the ICS-20 packet, the metadata can optionally indicate a module to route to (e.g. amm)
// The PacketForwardMetadata returned from this function contains attributes for each transferhook supported module
// It can only be forward to one module per packet
// Returns nil if there was no metadata found
func ParsePacketMetadata(metadata string) (*PacketForwardMetadata, error) {
	// If we can't unmarshal the metadata into a PacketMetadata struct,
	// assume packet forwarding was no used and pass back nil so that transferhook is ignored
	var raw RawPacketMetadata
	if err := json.Unmarshal([]byte(metadata), &raw); err != nil {
		return nil, nil
	}

	// If no forwarding logic was used for transferhook, return the metadata with each disabled
	if raw.Transferhook == nil {
		return nil, nil
	}

	// Confirm a receiver address was supplied
	if _, err := sdk.AccAddressFromBech32(raw.Transferhook.Receiver); err != nil {
		return nil, errorsmod.Wrapf(ErrInvalidPacketMetadata, ErrInvalidReceiverAddress.Error())
	}

	// Parse the packet info into the specific module type
	// We increment the module count to ensure only one module type was provided
	moduleCount := 0
	var routingInfo ModuleRoutingInfo
	if raw.Transferhook.Amm != nil {
		moduleCount++
		routingInfo = *raw.Transferhook.Amm
	}
	if moduleCount != 1 {
		return nil, errorsmod.Wrapf(ErrInvalidPacketMetadata, ErrInvalidModuleRoutes.Error())
	}

	// Validate the packet info according to the specific module type
	if err := routingInfo.Validate(); err != nil {
		return nil, errorsmod.Wrapf(err, ErrInvalidPacketMetadata.Error())
	}

	return &PacketForwardMetadata{
		Receiver:    raw.Transferhook.Receiver,
		RoutingInfo: routingInfo,
	}, nil
}
