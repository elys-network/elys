package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAprs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aprs",
		Short:   "Query aprs",
		Example: "elysd q masterchef aprs",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAprsRequest{}

			res, err := queryClient.Aprs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
