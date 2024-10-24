package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdApr() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apr",
		Short:   "calculate APR",
		Example: "elysd q masterchef apr [withdraw-type] [denom]",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			withdrawType, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			denom := args[1]

			params := &types.QueryAprRequest{
				WithdrawType: commitmenttypes.EarnType(withdrawType),
				Denom:        denom,
			}

			res, err := queryClient.Apr(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
