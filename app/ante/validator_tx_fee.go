package ante

import (
	"math"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	oracletypes "github.com/elys-network/elys/v5/x/oracle/types"
	parametertypes "github.com/elys-network/elys/v5/x/parameter/types"
)

// CheckTxFeeWithValidatorMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
func CheckTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() {
		minGasPrices := ctx.MinGasPrices()

		// Check for specific message types to adjust gas price
		msgs := tx.GetMsgs()
		if len(msgs) == 1 {
			msgType := strings.ToLower(sdk.MsgTypeURL(msgs[0]))
			if strings.Contains(msgType, sdk.MsgTypeURL(&oracletypes.MsgFeedPrice{})) || strings.Contains(msgType, sdk.MsgTypeURL(&oracletypes.MsgFeedMultiplePrices{})) {
				// set the minimum gas price to 0 ELYS if the message is a feed price
				minGasPrice := sdk.DecCoin{
					Denom:  parametertypes.Elys,
					Amount: sdkmath.LegacyZeroDec(),
				}
				if !minGasPrice.IsValid() {
					return nil, 0, errorsmod.Wrap(sdkerrors.ErrLogic, "invalid gas price")
				}
				minGasPrices = sdk.NewDecCoins(minGasPrice)

				// print minGasPrices
				ctx.Logger().Info("Override minimum gas prices: " + minGasPrices.String())
			}
		}

		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdkmath.LegacyNewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return nil, 0, errorsmod.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	priority := getTxPriority(feeCoins, int64(gas))
	return feeCoins, priority, nil
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
