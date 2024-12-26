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

func CmdChainTVL() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain_tvl",
		Short: "show chain tvl",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ChainTVL(context.Background(), &types.QueryChainTVLRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdAllLiquidityPoolTVL() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all_liquidity_pool_tvl",
		Short: "show all pools tvl",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AllLiquidityPoolTVL(context.Background(), &types.QueryAllLiquidityPoolTVLRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryStableStakeApr() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stable-stake-apr",
		Short:   "calculate Stable Stake APR",
		Example: "elysd q masterchef stable-stake-apr [denom]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]

			params := &types.QueryStableStakeAprRequest{
				Denom: denom,
			}

			res, err := queryClient.StableStakeApr(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPoolAprs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool-aprs",
		Short:   "calculate pool APRs",
		Example: "elysd q masterchef pool-aprs [ids]",
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

func CmdQueryExternalIncentive() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "external-incentive",
		Short:   "shows external incentive",
		Example: "elysd q masterchef external-incentive [id]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			id, _ := strconv.Atoi(args[0])
			params := &types.QueryExternalIncentiveRequest{
				Id: uint64(id),
			}

			res, err := queryClient.ExternalIncentive(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPoolInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool-info",
		Short:   "shows pool info",
		Example: "elysd q masterchef pool-info [id]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			id, _ := strconv.Atoi(args[0])
			params := &types.QueryPoolInfoRequest{
				PoolId: uint64(id),
			}

			res, err := queryClient.PoolInfo(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPoolRewardInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool-reward-info",
		Short:   "shows pool reward info",
		Example: "elysd q masterchef pool-reward-info [id] [reward-denom]",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			id, _ := strconv.Atoi(args[0])
			params := &types.QueryPoolRewardInfoRequest{
				PoolId:      uint64(id),
				RewardDenom: args[1],
			}

			res, err := queryClient.PoolRewardInfo(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryUserRewardInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user-reward-info",
		Short:   "shows user reward info",
		Example: "elysd q masterchef user-reward-info [user] [pool-id] [reward-denom]",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			id, _ := strconv.Atoi(args[1])
			params := &types.QueryUserRewardInfoRequest{
				User:        args[0],
				PoolId:      uint64(id),
				RewardDenom: args[2],
			}

			res, err := queryClient.UserRewardInfo(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryUserPendingReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user-pending-reward",
		Short:   "shows user pending reward",
		Example: "elysd q masterchef user-pending-reward [user]",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryUserPendingRewardRequest{
				User: args[0],
			}

			res, err := queryClient.UserPendingReward(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
