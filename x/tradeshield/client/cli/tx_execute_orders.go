package cli

import (
	"encoding/json"
	"os"
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
		Use:   "execute-orders [spot-order-ids] [perpetual-orders.json]",
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

			argPerpetualOrderKeys, err := readPerpetualOrderKeyJSON(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgExecuteOrders(
				clientCtx.GetFromAddress().String(),
				argSpotOrderIds,
				argPerpetualOrderKeys,
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

func readPerpetualOrderKeyJSON(filename string) ([]types.PerpetualOrderKey, error) {
	var perpetualOrderKeys []types.PerpetualOrderKey
	bz, err := os.ReadFile(filename)
	if err != nil {
		return []types.PerpetualOrderKey{}, err
	}
	err = json.Unmarshal(bz, &perpetualOrderKeys)
	if err != nil {
		return []types.PerpetualOrderKey{}, err
	}

	return perpetualOrderKeys, nil
}
