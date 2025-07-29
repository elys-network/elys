package cli

import (
	"errors"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdUpdateStopLoss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-stop-loss [amount] [id] [poolId]",
		Short: "Broadcast message update-stop-loss",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPrice, err := math.LegacyNewDecFromStr(args[0])
			if err != nil {
				return errors.New("invalid stoploss amount")
			}
			positionId, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateStopLoss(
				clientCtx.GetFromAddress().String(),
				uint64(positionId),
				argPrice,
				poolId,
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
