package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdQueryGetInterest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-get-interest",
		Short: "Query query-get-interest",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryQueryGetInterestRequest{}

			res, err := queryClient.QueryGetInterest(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
