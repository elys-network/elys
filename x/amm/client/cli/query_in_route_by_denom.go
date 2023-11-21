package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdInRouteByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "in-route-by-denom [denom-in] [denom-out]",
		Example: "elysd q amm in-route-by-denom uelys uusdc",
		Short:   "Query in-route-by-denom",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDenomIn := args[0]
			reqDenomOut := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryInRouteByDenomRequest{
				DenomIn:  reqDenomIn,
				DenomOut: reqDenomOut,
			}

			res, err := queryClient.InRouteByDenom(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
