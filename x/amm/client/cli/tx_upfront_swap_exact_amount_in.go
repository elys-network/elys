package cli

import (
	"errors"
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdUpFrontSwapExactAmountIn() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upfront-swap-exact-amount-in [token-in] [token-out-min-amount] [swap-route-pool-ids] [swap-route-denoms]",
		Short:   "Swap happens in the same Tx: Swap an exact amount of tokens for a minimum of another token, similar to swapping a token on the trade screen GUI.",
		Example: `elysd tx amm upfront-swap-exact-amount-in 100000uusdc 10000 0 uatom --from=bob --yes --gas=1000000`,
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

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpFrontSwapExactAmountIn(
				clientCtx.GetFromAddress().String(),
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

	return cmd
}
