package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdPoolRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool-rewards",
		Short:   "calculate pool rewards",
		Example: "elysd q masterchef pool-rewards [ids]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			idStrs := strings.Split(args[0], ",")
			ids := []uint64{}
			if args[0] != "" {
				for _, idStr := range idStrs {
					id, err := strconv.Atoi(idStr)
					if err != nil {
						return err
					}
					ids = append(ids, uint64(id))
				}
			}
			params := &types.QueryPoolRewardsRequest{
				PoolIds: ids,
			}

			res, err := queryClient.PoolRewards(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
