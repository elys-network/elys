package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/launchpad/types"
	"github.com/spf13/cobra"
)

func CmdQueryModuleBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "module-balances",
		Short:   "Query module balances",
		Example: "elysd q launchpad module-balances",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryModuleBalancesRequest{}

			res, err := queryClient.ModuleBalances(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
