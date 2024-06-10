package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/membershiptier/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreatePortfolio() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-portfolio [user]",
		Short: "Broadcast message create-portfolio",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argUser := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePortfolio(
				clientCtx.GetFromAddress().String(),
				argUser,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
