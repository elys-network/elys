package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/spf13/cobra"
)

func CmdSetPriceFeeder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-price-feeder [feeder] [isActive]",
		Short: "Update a priceFeeder",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			isActive, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgSetPriceFeeder(
				clientCtx.GetFromAddress().String(),
				args[0],
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
		Use:   "delete-price-feeder [feeder]",
		Short: "Delete a priceFeeder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeletePriceFeeder(
				clientCtx.GetFromAddress().String(),
				args[0],
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
