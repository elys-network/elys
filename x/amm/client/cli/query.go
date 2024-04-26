package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group amm queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListPool())
	cmd.AddCommand(CmdShowPool())
	cmd.AddCommand(CmdListDenomLiquidity())
	cmd.AddCommand(CmdShowDenomLiquidity())
	cmd.AddCommand(CmdSwapEstimation())
	cmd.AddCommand(CmdBalance())
	cmd.AddCommand(CmdInRouteByDenom())
	cmd.AddCommand(CmdOutRouteByDenom())
	cmd.AddCommand(CmdSwapEstimationByDenom())
	cmd.AddCommand(CmdJoinPoolEstimation())
	cmd.AddCommand(CmdExitPoolEstimation())
	cmd.AddCommand(CmdTrackedSlippage())
	cmd.AddCommand(CmdTrackedSlippageAll())

	// this line is used by starport scaffolding # 1

	return cmd
}
