package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/elys-network/elys/x/accountedpool/types"
)

func CmdListAccountedPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-accounted-pool",
		Short: "list all accounted-pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllAccountedPoolRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AccountedPoolAll(cmd.Context(), params)
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

func CmdShowAccountedPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-accounted-pool [index]",
		Short: "shows a accounted-pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]
			poolId, err := strconv.ParseUint(argIndex, 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetAccountedPoolRequest{
				PoolId: poolId,
			}

			res, err := queryClient.AccountedPool(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
