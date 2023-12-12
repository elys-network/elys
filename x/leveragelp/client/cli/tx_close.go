package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "close [position-id] [flags]",
		Short:   "Close leveragelp position",
		Example: `elysd tx leveragelp close 1 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			argPositionId, ok := strconv.ParseUint(args[0], 10, 64)
			if ok != nil {
				return errors.New("invalid position id")
			}

			msg := types.NewMsgClose(
				signer.String(),
				argPositionId,
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
