package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/perpetual/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group perpetual queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdGetPositions())
	cmd.AddCommand(CmdGetPositionsByPool())
	cmd.AddCommand(CmdGetStatus())
	cmd.AddCommand(CmdGetPositionsForAddress())
	cmd.AddCommand(CmdGetWhitelist())
	cmd.AddCommand(CmdIsWhitelisted())
	cmd.AddCommand(CmdListPool())
	cmd.AddCommand(CmdShowPool())
	cmd.AddCommand(CmdMtp())
	cmd.AddCommand(CmdOpenEstimation())

	cmd.AddCommand(CmdGetAllToPay())

	cmd.AddCommand(CmdCloseEstimation())

	// this line is used by starport scaffolding # 1

	return cmd
}
