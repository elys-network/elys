package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/burner/types"
	"github.com/spf13/cobra"
)

func CmdListHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-history",
		Short: "list all history",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllHistoryRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.HistoryAll(context.Background(), params)
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

func CmdShowHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-history [timestamp] [denom]",
		Short: "shows a history",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argTimestamp := args[0]
			argDenom := args[1]

			params := &types.QueryGetHistoryRequest{
				Timestamp: argTimestamp,
				Denom:     argDenom,
			}

			res, err := queryClient.History(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
