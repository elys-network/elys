package cli

import (
	"errors"
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdSwapExactAmountIn() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swap-exact-amount-in [token-in] [token-out-min-amount] [swap-route-pool-ids] [swap-route-denoms]",
		Short:   "Swap an exact amount of tokens for a minimum of another token, similar to swapping a token on the trade screen GUI.",
		Example: `elysd tx amm swap-exact-amount-in 100000uusdc 10000 0 uatom --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTokenIn, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			argTokenOutMinAmount, ok := math.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid token-out-min-amount")
			}
			argCastSwapRoutePoolIds := strings.Split(args[2], listSeparator)
			argSwapRoutePoolIds := make([]uint64, len(argCastSwapRoutePoolIds))
			for i, arg := range argCastSwapRoutePoolIds {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argSwapRoutePoolIds[i] = value
			}
			argSwapRouteDenoms := strings.Split(args[3], listSeparator)

			recipient, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSwapExactAmountIn(
				clientCtx.GetFromAddress().String(),
				recipient,
				argTokenIn,
				argTokenOutMinAmount,
				argSwapRoutePoolIds,
				argSwapRouteDenoms,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagRecipient, "", "optional recipient field for the tokens swapped to be sent to")

	return cmd
}
