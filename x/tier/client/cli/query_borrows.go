package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tier/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdTotalBorrows() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-borrows [user]",
		Short: "Query total-borrows",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqUser := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTotalBorrowsRequest{
				User: reqUser,
			}

			res, err := queryClient.QueryTotalBorrows(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
