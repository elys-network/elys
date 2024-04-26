package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "close [position-id] [amount] [flags]",
		Short:   "Close leveragelp position",
		Example: `elysd tx leveragelp close 1 10000000 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			argPositionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errors.New("invalid position id")
			}

			argAmount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := types.NewMsgClose(
				signer.String(),
				argPositionId,
				argAmount,
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
