package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdOutRouteByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "out-route-by-denom [denom-out] [denom-in]",
		Example: "elysd q amm out-route-by-denom uusdc uelys",
		Short:   "Query out-route-by-denom",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDenomOut := args[0]
			reqDenomIn := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryOutRouteByDenomRequest{
				DenomOut: reqDenomOut,
				DenomIn:  reqDenomIn,
			}

			res, err := queryClient.OutRouteByDenom(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
