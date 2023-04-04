package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/spf13/cobra"
)

func CmdListPriceFeeder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-price-feeder",
		Short: "list all priceFeeder",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPriceFeederRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PriceFeederAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowPriceFeeder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-price-feeder [index]",
		Short: "shows a priceFeeder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetPriceFeederRequest{
				Feeder: args[0],
			}

			res, err := queryClient.PriceFeeder(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
