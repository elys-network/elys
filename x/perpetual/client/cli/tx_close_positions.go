package cli

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClosePositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-positions [liquidate] [stoploss] [take-profit]",
		Short: "Broadcast message close-positions",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidate, err := readPositionRequestJSON(args[0])
			if err != nil {
				return err
			}
			// Convert to slice of pointers
			var liqudiatePtrs []types.PositionRequest
			for i := range liquidate {
				liqudiatePtrs = append(liqudiatePtrs, liquidate[i])
			}

			stopLoss, err := readPositionRequestJSON(args[1])
			if err != nil {
				return err
			}
			// Convert to slice of pointers
			var stoplossPtrs []types.PositionRequest
			for i := range stopLoss {
				stoplossPtrs = append(stoplossPtrs, stopLoss[i])
			}

			takeProfit, err := readPositionRequestJSON(args[2])
			if err != nil {
				return err
			}
			// Convert to slice of pointers
			var takeProfitPtrs []types.PositionRequest
			for i := range takeProfit {
				takeProfitPtrs = append(takeProfitPtrs, takeProfit[i])
			}

			msg := types.NewMsgClosePositions(
				clientCtx.GetFromAddress().String(),
				liqudiatePtrs,
				stoplossPtrs,
				takeProfitPtrs,
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

func readPositionRequestJSON(filename string) ([]types.PositionRequest, error) {
	var positions []types.PositionRequest
	bz, err := os.ReadFile(filename)
	if err != nil {
		return []types.PositionRequest{}, err
	}
	err = json.Unmarshal(bz, &positions)
	if err != nil {
		return []types.PositionRequest{}, err
	}

	return positions, nil
}
