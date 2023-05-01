package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSetWithdrawAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-withdraw-address [delegator-address] [withdraw-address]",
		Short: "Broadcast message set-withdraw-address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argWithdrawAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetWithdrawAddress(
				clientCtx.GetFromAddress().String(),
				argWithdrawAddress,
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
