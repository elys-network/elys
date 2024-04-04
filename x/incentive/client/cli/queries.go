package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/spf13/cobra"
)

func CmdCommunityPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-pool",
		Short: "Query community-pool",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCommunityPoolRequest{}

			res, err := queryClient.CommunityPool(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdApr() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apr",
		Short:   "calculate APR",
		Example: "elysd q incentive apr [withdraw-type] [denom]",
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

func CmdAprs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aprs",
		Short:   "calculate APRs",
		Example: "elysd q incentive aprs",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAprsRequest{}

			res, err := queryClient.Aprs(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdPoolAprs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool-aprs",
		Short:   "calculate pool APRs",
		Example: "elysd q incentive pool-aprs [ids]",
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
			params := &types.QueryPoolAprsRequest{
				PoolIds: ids,
			}

			res, err := queryClient.PoolAprs(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdPoolRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool-rewards",
		Short:   "calculate pool rewards",
		Example: "elysd q incentive pool-rewards [ids]",
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
