package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
	"github.com/spf13/cobra"
)

func CmdQueryReturnElysEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "return-elys-estimation [orderId] [amount]",
		Short:   "Query ReturnElys estimation",
		Example: "elysd q launchpad return-elys-estimation [orderId] [amount]",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			orderId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid spending token amount")
			}

			params := &types.QueryReturnElysEstRequest{
				OrderId:    uint64(orderId),
				ElysAmount: amount,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ReturnElysEst(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
