package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdGetPositionsByPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-positions-by-pool [amm_pool_id]",
		Short: "Query get-positions-by-pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqId := args[0]
			id, err := strconv.ParseUint(reqId, 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.PositionsByPoolRequest{
				AmmPoolId: id,
			}

			res, err := queryClient.QueryPositionsByPool(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
