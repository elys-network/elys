package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/tier/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdLeverageLpTotal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leverage-lp-total [user]",
		Short: "Query leverage-lp-total",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqUser := args[0]
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLeverageLpTotalRequest{
				User: reqUser,
			}

			res, err := queryClient.LeverageLpTotal(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
