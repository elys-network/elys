package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "vaults",
		Short:                      "Transaction commands for the vaults module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdPerformActionJoinPool())
	cmd.AddCommand(CmdPerformActionExitPool())
	cmd.AddCommand(CmdPerformActionSwapByDenom())
	// this line is used by ignite scaffolding # 1

	return cmd
}
