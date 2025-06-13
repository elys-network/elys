package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdOpen())
	cmd.AddCommand(CmdClose())
	cmd.AddCommand(CmdUpdateParams())
	cmd.AddCommand(CmdWhitelist())
	cmd.AddCommand(CmdDewhitelist())
	cmd.AddCommand(CmdClaimRewards())
	cmd.AddCommand(CmdClaimAllUserRewards())
	cmd.AddCommand(CmdUpdateStopLoss())
	cmd.AddCommand(CmdClosePositions())
	// this line is used by starport scaffolding # 1

	return cmd
}
