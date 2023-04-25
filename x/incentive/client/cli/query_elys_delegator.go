package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/spf13/cobra"
)

func CmdListElysDelegator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-elys-delegator",
		Short: "list all elys-delegator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllElysDelegatorRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ElysDelegatorAll(context.Background(), params)
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

func CmdShowElysDelegator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-elys-delegator [index]",
		Short: "shows a elys-delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]

			params := &types.QueryGetElysDelegatorRequest{
				Index: argIndex,
			}

			res, err := queryClient.ElysDelegator(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
