package cli

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClosePositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-positions [liquidate.json] [stoploss.json]",
		Short: "Broadcast message close-positions",
		Args:  cobra.ExactArgs(2),
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
			var liqudiatePtrs []*types.PositionRequest
			for i := range liquidate {
				liqudiatePtrs = append(liqudiatePtrs, &liquidate[i])
			}

			stopLoss, err := readPositionRequestJSON(args[1])
			if err != nil {
				return err
			}
			// Convert to slice of pointers
			var stoplossPtrs []*types.PositionRequest
			for i := range stopLoss {
				stoplossPtrs = append(stoplossPtrs, &stopLoss[i])
			}

			msg := types.NewMsgClosePositions(
				clientCtx.GetFromAddress(),
				liqudiatePtrs,
				stoplossPtrs,
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
	bz, err := ioutil.ReadFile(filename)
	if err != nil {
		return []types.PositionRequest{}, err
	}
	err = json.Unmarshal(bz, &positions)
	if err != nil {
		return []types.PositionRequest{}, err
	}

	return positions, nil
}
