package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

const (
	FlagMinAmount = "min-amount"
	FlagMaxAmount = "max-amount"
)

func CmdSwapByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swap-by-denom [amount] [denom-in] [denom-out] [recipient]",
		Short:   "Swap an exact amount of tokens for a minimum of another token or a maximum amount of tokens for an exact amount on another token, similar to swapping a token on the trade screen GUI.",
		Example: "elysd tx amm swap-by-denom 1000000000uatom uatom uusd --min-amount=1000000000uatom --max-amount=1000000000uatom --discount=0.1 --from jack --keyring-backend test",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			argDenomIn := args[1]
			argDenomOut := args[2]
			argRecipient := args[3]

			minAmountStr, err := cmd.Flags().GetString(FlagMinAmount)
			if err != nil {
				return err
			}
			minAmount := sdk.Coin{}
			if minAmountStr != "" {
				minAmount, err = sdk.ParseCoinNormalized(minAmountStr)
				if err != nil {
					return err
				}
			}

			maxAmountStr, err := cmd.Flags().GetString(FlagMaxAmount)
			if err != nil {
				return err
			}
			maxAmount := sdk.Coin{}
			if maxAmountStr == "" {
				maxAmount, err = sdk.ParseCoinNormalized(maxAmountStr)
				if err != nil {
					return err
				}
			}

			discountStr, err := cmd.Flags().GetString(FlagDiscount)
			if err != nil {
				return err
			}
			discount, err := sdk.NewDecFromStr(discountStr)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSwapByDenom(
				clientCtx.GetFromAddress().String(),
				argRecipient,
				argAmount,
				minAmount,
				maxAmount,
				argDenomIn,
				argDenomOut,
				discount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagMinAmount, "", "minimum amount of tokens to receive")
	cmd.Flags().String(FlagMaxAmount, "", "maximum amount of tokens to send")
	cmd.Flags().String(FlagDiscount, "0.0", "discount to apply to the swap fee (only smart contract broker can apply the discount)")

	return cmd
}
