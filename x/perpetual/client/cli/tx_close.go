package cli

import (
	"errors"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "close [mtp-id] [amount] [flags]",
		Short:   "Close perpetual position",
		Example: `elysd tx perpetual close 1 10000000 --from=bob --yes --gas=1000000`,
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

			argMtpId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errors.New("invalid mtp id")
			}

			argAmount, ok := math.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := types.NewMsgClose(
				signer.String(),
				argMtpId,
				argAmount,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
