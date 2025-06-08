package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/elys-network/elys/v6/x/vaults/types"
)

func CmdPerformAction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perform-action [vault-id]",
		Short: "Perform an action on a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get vault ID
			vaultId, err := parseUint64(args[0])
			if err != nil {
				return fmt.Errorf("invalid vault ID: %w", err)
			}

			// Get action JSON
			actionJSON, err := cmd.Flags().GetString("action")
			if err != nil {
				return fmt.Errorf("failed to get action: %w", err)
			}

			// Parse action JSON
			var action types.Action
			if err := json.Unmarshal([]byte(actionJSON), &action); err != nil {
				return fmt.Errorf("invalid action JSON: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgPerformAction{
				Creator: clientCtx.GetFromAddress().String(),
				VaultId: vaultId,
				Action:  &action,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("action", "", "The action to perform (JSON string)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseUint64(s string) (uint64, error) {
	var i uint64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
