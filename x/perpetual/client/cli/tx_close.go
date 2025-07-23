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
		Use:     "close [mtp-id] [closingRatio] [pool-id]",
		Short:   "Close perpetual position",
		Example: `elysd tx perpetual close 1 0.37 2 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(3),
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

			argClosingRatio, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgClose(
				signer.String(),
				argMtpId,
				math.ZeroInt(),
				poolId,
				argClosingRatio,
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
