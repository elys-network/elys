package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/spf13/cobra"
)

func CmdSetAssetInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-asset-info [denom] [display] [bandTicker] [binanceTicker] [osmosisTicker]",
		Short: "Update an asset info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetAssetInfo(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
				args[2],
				args[3],
				args[4],
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

func CmdDeleteAssetInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-asset-info [denom]",
		Short: "Delete an asset info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteAssetInfo(
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
