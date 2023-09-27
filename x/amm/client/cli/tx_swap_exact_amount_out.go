package cli

import (
	"errors"
	"strconv"

	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSwapExactAmountOut() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swap-exact-amount-out [token-out] [token-out-max-amount] [swap-route-pool-ids] [swap-route-denoms]",
		Short:   "Swap a maximum amount of tokens for an exact amount of another token, similar to swapping a token on the trade screen GUI.",
		Example: `elysd tx amm swap-exact-amount-out 100000uatom 200000 0 uusdc --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTokenOut, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			argTokenOutMaxAmount, ok := math.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid token-out-max-amount")
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

			msg := types.NewMsgSwapExactAmountOut(
				clientCtx.GetFromAddress().String(),
				argTokenOut,
				argTokenOutMaxAmount,
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
