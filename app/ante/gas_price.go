package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	parametertypes "github.com/elys-network/elys/x/parameter/types"
)

// AdjustGasPriceDecorator is a custom decorator to reduce fee prices .
type AdjustGasPriceDecorator struct {
}

// NewAdjustGasPriceDecorator create a new instance of AdjustGasPriceDecorator
func NewAdjustGasPriceDecorator() AdjustGasPriceDecorator {
	return AdjustGasPriceDecorator{}
}

// AnteHandle adjusts the gas price based on the tx type.
func (r AdjustGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	// if there are no messages or there are multiple messages, return
	if len(msgs) != 1 {
		return next(ctx, tx, simulate)
	}

	// retrieve the first message
	m := msgs[0]

	switch m.(type) {
	case *oracletypes.MsgFeedPrice:
	case *oracletypes.MsgFeedMultiplePrices:
		minGasPrice := sdk.DecCoin{
			Denom:  parametertypes.Elys,
			Amount: sdk.MustNewDecFromStr("0.00000001"),
		}
		if !minGasPrice.IsValid() {
			return ctx, errorsmod.Wrap(sdkerrors.ErrLogic, "invalid gas price")
		}
		ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(minGasPrice))
	}

	return next(ctx, tx, simulate)
}
