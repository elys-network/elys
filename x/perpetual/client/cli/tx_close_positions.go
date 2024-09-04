package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClosePositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-positions [liquidate] [stoploss]",
		Short: "Broadcast message close-positions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argLiquidate := args[0]
			argStoploss := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgClosePositions(
				clientCtx.GetFromAddress().String(),
				argLiquidate,
				argStoploss,
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
