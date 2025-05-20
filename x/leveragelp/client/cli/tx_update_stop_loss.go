package cli

import (
	sdkmath "cosmossdk.io/math"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateStopLoss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-stop-loss [position] [price]",
		Short: "Broadcast message update-stop-loss",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPosition, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argPrice, err := sdkmath.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateStopLoss(
				clientCtx.GetFromAddress().String(),
				argPosition,
				argPrice,
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
