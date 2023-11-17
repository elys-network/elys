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

func CmdBrokerClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "broker-close [mtp-id] [mtp-owner] [flags]",
		Short:   "Close margin position as a broker",
		Example: `elysd tx margin broker-close 1 sif123 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
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

			argMtpId, ok := strconv.ParseUint(args[0], 10, 64)
			if ok != nil {
				return errors.New("invalid mtp id")
			}

			argMtpOwner := args[1]
			if argMtpOwner == "" {
				return errors.New("invalid mtp owner address")
			}

			msg := types.NewMsgBrokerClose(
				signer.String(),
				argMtpId,
				argMtpOwner,
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
