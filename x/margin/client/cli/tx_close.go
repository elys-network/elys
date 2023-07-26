package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close margin position",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			id, err := cmd.Flags().GetUint64("id")
			if err != nil {
				return err
			}

			msg := types.NewMsgClose(
				signer.String(),
				id,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Uint64("id", 0, "id of the position")
	_ = cmd.MarkFlagRequired("id")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
