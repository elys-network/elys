package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/transferhook module sentinel errors
var (
	ErrInvalidPacketMetadata        = errorsmod.Register(ModuleName, 1501, "invalid packet metadata")
	ErrUnsupportedAmmAction         = errorsmod.Register(ModuleName, 1502, "unsupported amm action")
	ErrInvalidModuleRoutes          = errorsmod.Register(ModuleName, 1504, "invalid number of module routes, only 1 module is allowed at a time")
	ErrUnsupportedTransferhookRoute = errorsmod.Register(ModuleName, 1505, "unsupported transferhook route")
	ErrInvalidReceiverAddress       = errorsmod.Register(ModuleName, 1506, "receiver address must be specified when using transferhook")
	ErrPacketForwardingInactive     = errorsmod.Register(ModuleName, 1507, "transferhook packet forwarding is disabled")
)
