package ante

import (
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	parametertypes "github.com/elys-network/elys/v6/x/parameter/types"
)

// GaslessAddrs contains only the governance module address.
// Governance can grant feegrants to operational addresses (price feeders, liquidation bots, etc.)
// When those addresses use feegrants, the fee payer becomes governance, enabling gasless transactions.
var GaslessAddrs = []string{
	authtypes.NewModuleAddress(govtypes.ModuleName).String(), // Governance module address
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
  isGasless := false

  if ctx.IsCheckTx() {
    // start with the validator's configured minGasPrices
    minGasPrices := ctx.MinGasPrices()

    // Check if the fee payer is in the gasless whitelist
    feePayerBytes := feeTx.FeePayer()
    if len(feePayerBytes) > 0 {
      feePayer := sdk.AccAddress(feePayerBytes).String()
      
      for _, ga := range GaslessAddrs {
        if feePayer == ga {
          isGasless = true
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
  // Give gasless transactions the highest priority
  var priority int64
  if isGasless {
    priority = math.MaxInt64
    ctx.Logger().Info(
      "gasless transaction given highest priority",
      "priority", priority,
    )
  } else {
    priority = getTxPriority(feeCoins, int64(gas))
  }
  
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
