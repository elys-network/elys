package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/assetprofile/types"
	"github.com/spf13/cobra"
)

func CmdListEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-entry",
		Short: "list all entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllEntryRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.EntryAll(context.Background(), params)
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

func CmdShowEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-entry [base-denom]",
		Short: "shows a entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argBaseDenom := args[0]

			params := &types.QueryEntryRequest{
				BaseDenom: argBaseDenom,
			}

			res, err := queryClient.Entry(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowEntryByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-entry-by-denom [denom]",
		Short: "shows a entry by denom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argDenom := args[0]

			params := &types.QueryEntryByDenomRequest{
				Denom: argDenom,
			}

			res, err := queryClient.EntryByDenom(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
