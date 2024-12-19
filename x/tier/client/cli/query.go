package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/tier/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group membershiptier queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListPortfolio())
	cmd.AddCommand(CmdShowPortfolio())
	cmd.AddCommand(CmdCalculateDiscount())

	cmd.AddCommand(CmdLeverageLpTotal())

	cmd.AddCommand(CmdRewardsTotal())

	cmd.AddCommand(CmdStakedPool())

	cmd.AddCommand(CmdPerpetual())

	cmd.AddCommand(CmdLiquidTotal())

	cmd.AddCommand(CmdGetAmmPrice())

	cmd.AddCommand(CmdGetConsolidatedPrice())

	cmd.AddCommand(CmdStaked())

	cmd.AddCommand(CmdGetUsersPoolData())

	cmd.AddCommand(CmdLockedOrder())

	cmd.AddCommand(CmdGetAllPrices())

	// this line is used by starport scaffolding # 1

	return cmd
}
