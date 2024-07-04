package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/elys-network/elys/x/perpetual/types"
)

const (
	FlagDiscount        = "discount"
	FlagTakeProfitPrice = "take-profit"
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
	cmd.AddCommand(CmdAddCollateral())
	// this line is used by starport scaffolding # 1

	return cmd
}
