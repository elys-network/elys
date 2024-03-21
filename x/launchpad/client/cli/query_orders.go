package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/launchpad/types"
	"github.com/spf13/cobra"
)

func CmdQueryOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "orders [user]",
		Short:   "Query Orders",
		Example: "elysd q launchpad orders [user]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryOrdersRequest{
				User: args[0],
			}

			res, err := queryClient.Orders(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAllOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "all-orders",
		Short:   "Query Orders",
		Example: "elysd q launchpad all-orders",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllOrdersRequest{}

			res, err := queryClient.AllOrders(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
