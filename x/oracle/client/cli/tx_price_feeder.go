package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v5/x/oracle/types"
	"github.com/spf13/cobra"
)

func CmdSetPriceFeeder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-price-feeder [isActive]",
		Short: "Set a price feeder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			isActive, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgSetPriceFeeder(
				clientCtx.GetFromAddress().String(),
				isActive,
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

func CmdDeletePriceFeeder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-price-feeder",
		Short: "Delete a priceFeeder",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeletePriceFeeder(
				clientCtx.GetFromAddress().String(),
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
