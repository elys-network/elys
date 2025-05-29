package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	commitmenttypes "github.com/elys-network/elys/v5/x/commitment/types"
	"github.com/elys-network/elys/v5/x/masterchef/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdApr() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apr",
		Short:   "calculate APR",
		Example: "elysd q masterchef apr [withdraw-type] [denom] [days]",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			withdrawType, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			denom := args[1]
			days, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			params := &types.QueryAprRequest{
				WithdrawType: commitmenttypes.EarnType(withdrawType),
				Denom:        denom,
				Days:         uint64(days),
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
