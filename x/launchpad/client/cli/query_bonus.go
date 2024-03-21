package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/launchpad/types"
	"github.com/spf13/cobra"
)

func CmdQueryBonus() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bonus [user]",
		Short:   "Query Bonus",
		Example: "elysd q launchpad bonus [user]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBonusRequest{
				User: args[0],
			}

			res, err := queryClient.Bonus(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
