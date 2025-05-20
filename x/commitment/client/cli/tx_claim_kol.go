package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v4/x/commitment/types"
	"github.com/spf13/cobra"
)

const FlagRefund = "refund"

func CmdClaimKol() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-kol",
		Short: "Broadcast message claim_kol",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			refund, err := cmd.Flags().GetBool(FlagRefund)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimKol(
				clientCtx.GetFromAddress().String(),
				refund,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Bool(FlagRefund, false, "flag to set refund to true")

	return cmd
}
