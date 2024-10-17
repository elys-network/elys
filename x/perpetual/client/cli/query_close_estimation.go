package cli

import (
	"cosmossdk.io/math"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCloseEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-estimation [address] [position-id] [closing-amount] [closing-price]",
		Short: "Query close-estimation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]

			reqPositionId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			reqClosingAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid closing amount")
			}

			reqClosingPrice, err := math.LegacyNewDecFromStr(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCloseEstimationRequest{
				Address:      reqAddress,
				PositionId:   reqPositionId,
				CloseAmount:  reqClosingAmount,
				ClosingPrice: reqClosingPrice,
			}

			res, err := queryClient.CloseEstimation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
