package cli

import (
	"strconv"

	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdExecuteOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-orders [spot-order-ids] [perpetual-order-ids]",
		Short: "Verify that submitted orders meet the criteria for execution and process those that do, while skipping those that don't.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCastSpotOrderIds := strings.Split(args[0], listSeparator)
			if len(argCastSpotOrderIds) == 1 && argCastSpotOrderIds[0] == "" {
				argCastSpotOrderIds = []string{}
			}
			argSpotOrderIds := make([]uint64, len(argCastSpotOrderIds))
			for i, arg := range argCastSpotOrderIds {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argSpotOrderIds[i] = value
			}
			argCastPerpetualOrderIds := strings.Split(args[1], listSeparator)
			if len(argCastPerpetualOrderIds) == 1 && argCastPerpetualOrderIds[0] == "" {
				argCastPerpetualOrderIds = []string{}
			}
			argPerpetualOrderIds := make([]uint64, len(argCastPerpetualOrderIds))
			for i, arg := range argCastPerpetualOrderIds {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argPerpetualOrderIds[i] = value
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgExecuteOrders(
				clientCtx.GetFromAddress().String(),
				argSpotOrderIds,
				argPerpetualOrderIds,
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
