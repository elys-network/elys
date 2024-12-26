package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group masterchef queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdChainTVL())
	cmd.AddCommand(CmdAllLiquidityPoolTVL())
	cmd.AddCommand(CmdQueryExternalIncentive())
	cmd.AddCommand(CmdQueryPoolInfo())
	cmd.AddCommand(CmdQueryPoolRewardInfo())
	cmd.AddCommand(CmdQueryUserRewardInfo())
	cmd.AddCommand(CmdQueryUserPendingReward())
	cmd.AddCommand(CmdQueryStableStakeApr())
	cmd.AddCommand(CmdQueryPoolAprs())
	cmd.AddCommand(CmdShowFeeInfo())

	cmd.AddCommand(CmdListFeeInfo())

	cmd.AddCommand(CmdApr())

	cmd.AddCommand(CmdAprs())

	cmd.AddCommand(CmdPoolRewards())

	// this line is used by starport scaffolding # 1

	return cmd
}
