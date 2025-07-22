package ante

import (
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	parametertypes "github.com/elys-network/elys/v6/x/parameter/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// GaslessAddrs is the whitelist of Bech32 Elys addresses that may send txs
// with zero fees (in mempool).
var GaslessAddrs = []string{
  "elys16qgewtplahkqwqa0aqv4pxnxa58ulu48k6crhj", // Price Feeder 1 mainnet
  "elys1zwexzk6ns5ermvag5fc0gtyvrnxyaz9kzaflqf", // Price Feeder 2 mainnet
  "elys1nelawytdfdk4af3z0sy2p8vkrllk8zw9g32jmf", // Execution bot 1 mainnet
  "elys1gszp63euzm0ecs3qwu0j6mexjr97hjs7x5gvzk", // Execution bot 2 mainnet


  "elys1dae9z45ccetfwr208ghya6npntg75qvxgmg4p9", // Price Feeder testnet
  "elys18z3qtz4wag4mmv9f8ea25tpvkmkfk4vtyy6fpr", // Execution bot 1 testnet
  "elys135x29zpaph6sacf3kvu8p736jcue5qd30nf8mv", // Execution bot 2 testnet

  "elys1gy503ty29ydute5rksnkwgmtatelret9d455lt", // Execution bot devnet
  "elys1dae9z45ccetfwr208ghya6npntg75qvxgmg4p9", // Price Feeder devnet
}

// CheckTxFeeWithValidatorMinGasPrices checks fees in CheckTx (mempool).  If the
// fee-payer is in GaslessAddrs, we force minGasPrices to zero so 0uelys txs pass.
func CheckTxFeeWithValidatorMinGasPrices(
  ctx sdk.Context, tx sdk.Tx,
) (sdk.Coins, int64, error) {
  feeTx, ok := tx.(sdk.FeeTx)
  if !ok {
    return nil, 0, errorsmod.Wrap(sdkerrors.ErrTxDecode,
      "tx must implement FeeTx")
  }

  feeCoins, gas := feeTx.GetFee(), feeTx.GetGas()

  if ctx.IsCheckTx() {
    // start with the validator's configured minGasPrices
    minGasPrices := ctx.MinGasPrices()

    // Check if the fee payer is in the gasless whitelist
    feePayerBytes := feeTx.FeePayer()
    if len(feePayerBytes) > 0 {
      feePayer := sdk.AccAddress(feePayerBytes).String()
      
      for _, ga := range GaslessAddrs {
        if feePayer == ga {
          zero := sdkmath.LegacyZeroDec()
          minGasPrices = sdk.NewDecCoins(
            sdk.NewDecCoinFromDec(parametertypes.Elys, zero),
          )
          ctx.Logger().Info(
            "override minimum gas price to 0 for gasless address",
            "address", ga,
          )
          break
        }
      }
    }

    // now enforce: feeCoins >= minGasPrices * gasLimit
    if !minGasPrices.IsZero() {
      required := make(sdk.Coins, len(minGasPrices))
      gdec := sdkmath.LegacyNewDec(int64(gas))
      for i, gp := range minGasPrices {
        amt := gp.Amount.Mul(gdec).Ceil().RoundInt()
        required[i] = sdk.NewCoin(gp.Denom, amt)
      }
      if !feeCoins.IsAnyGTE(required) {
        return nil, 0, errorsmod.Wrapf(
          sdkerrors.ErrInsufficientFee,
          "insufficient fees; got: %s required: %s",
          feeCoins, required,
        )
      }
    }
  }

  // on DeliverTx/Simulate we just deduct whatever feeCoins was attached
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
