package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tradeshield/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
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

	cmd.AddCommand(CmdCreatePendingSpotOrder())
	cmd.AddCommand(CmdUpdatePendingSpotOrder())
	cmd.AddCommand(CmdDeletePendingSpotOrder())
	cmd.AddCommand(CmdCreatePendingPerpetualOrder())
	cmd.AddCommand(CmdUpdatePendingPerpetualOrder())
	cmd.AddCommand(CmdDeletePendingPerpetualOrder())
	cmd.AddCommand(CmdUpdateParams())
	cmd.AddCommand(CmdExecuteOrders())
	// this line is used by starport scaffolding # 1

	return cmd
}
