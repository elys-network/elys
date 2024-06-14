package cli

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdCloseEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-estimation [owner] [pool_id] [lp_amount]",
		Short: "Query close estimation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			positionId, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			lpAmount, ok := math.NewIntFromString(args[2])
			if !ok {
				return err
			}

			params := &types.QueryCloseEstRequest{
				Owner:    args[0],
				Id:       uint64(positionId),
				LpAmount: lpAmount,
			}

			res, err := queryClient.CloseEst(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
