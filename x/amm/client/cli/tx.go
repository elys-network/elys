package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/elys-network/elys/v4/x/amm/types"
)

const (
	FlagAddress   = "address"
	FlagRecipient = "recipient"
	listSeparator = ","
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

	cmd.AddCommand(CmdCreatePool())
	cmd.AddCommand(CmdJoinPool())
	cmd.AddCommand(CmdExitPool())
	cmd.AddCommand(CmdUpFrontSwapExactAmountIn())
	cmd.AddCommand(CmdSwapExactAmountIn())
	cmd.AddCommand(CmdSwapExactAmountOut())
	cmd.AddCommand(CmdSwapByDenom())
	cmd.AddCommand(CmdUpdatePoolParams())
	// this line is used by starport scaffolding # 1

	return cmd
}
