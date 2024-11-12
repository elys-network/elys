package cli

import (
	"errors"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdBrokerClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "broker-close [mtp-id] [amount] [owner] [flags]",
		Short:   "Broker Closes perpetual position",
		Example: `elysd tx perpetual broker-close 1 10000000 elys1w9uac4zrf9z7qd604qxk2y4n74568lfl8vutz4 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			_, err = sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return errors.New("invalid owner address")
			}

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

			msg := types.NewMsgBrokerClose(
				signer.String(),
				argMtpId,
				argAmount,
				args[2],
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
