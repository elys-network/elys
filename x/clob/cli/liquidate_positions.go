package cli

import (
	"encoding/json"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/clob/types"
	"github.com/spf13/cobra"
)

func CmdLiquidatePositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "liquidate-positions [positions]",
		Short:   "liquidate-positions force-liquidate the perpetual positions",
		Example: `elysd tx clob liquidate-positions json-file-path (eg. - '[{"market_id": 1, "perpetual_id": 4}...]') --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidatePositions, err := readLiquidatePositionsJSON(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgLiquidatePositions{
				Liquidator: clientCtx.GetFromAddress().String(),
				Positions:  liquidatePositions,
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func readLiquidatePositionsJSON(filename string) ([]types.LiquidatePosition, error) {
	var positions []types.LiquidatePosition
	bz, err := os.ReadFile(filename)
	if err != nil {
		return []types.LiquidatePosition{}, err
	}
	err = json.Unmarshal(bz, &positions)
	if err != nil {
		return []types.LiquidatePosition{}, err
	}

	return positions, nil
}
