package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/commitment/types"
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

	cmd.AddCommand(CmdCommitClaimedRewards())
	cmd.AddCommand(CmdUncommitTokens())
	cmd.AddCommand(CmdVest())
	cmd.AddCommand(CmdCancelVest())
	cmd.AddCommand(CmdVestNow())
	cmd.AddCommand(CmdUpdateVestingInfo())
	cmd.AddCommand(CmdVestLiquid())
	cmd.AddCommand(CmdStake())
	cmd.AddCommand(CmdUnstake())
	cmd.AddCommand(CmdClaimAirdrop())
	// this line is used by starport scaffolding # 1

	return cmd
}
