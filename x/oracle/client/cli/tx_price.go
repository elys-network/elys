package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/oracle/types"
	"github.com/spf13/cobra"
)

func CmdFeedPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-price [asset] [price] [source]",
		Short: "Feed a new price",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get value arguments
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			price, err := sdkmath.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgFeedPrice(
				clientCtx.GetFromAddress().String(),
				args[0],
				price,
				args[2],
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
