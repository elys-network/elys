package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetPositionsByPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-positions-by-pool [asset]",
		Short: "Query get-positions-by-pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAsset := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.PositionsByPoolRequest{

				Asset: reqAsset,
			}

			res, err := queryClient.GetPositionsByPool(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
