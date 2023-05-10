package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

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

func CmdDelegationRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation-rewards [delegator-address] [validator-address]",
		Short: "Query delegation-rewards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDelegatorAddress := args[0]
			reqValidatorAddress := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegationRewardsRequest{

				DelegatorAddress: reqDelegatorAddress,
				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.DelegationRewards(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdDelegationTotalRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation-total-rewards [delegator-address]",
		Short: "Query delegation-total-rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDelegatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegationTotalRewardsRequest{

				DelegatorAddress: reqDelegatorAddress,
			}

			res, err := queryClient.DelegationTotalRewards(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdDelegatorValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-validators [delegator-address]",
		Short: "Query delegator-validators",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDelegatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegatorValidatorsRequest{

				DelegatorAddress: reqDelegatorAddress,
			}

			res, err := queryClient.DelegatorValidators(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdDelegatorWithdrawAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-withdraw-address [delegator-address]",
		Short: "Query delegator-withdraw-address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDelegatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDelegatorWithdrawAddressRequest{

				DelegatorAddress: reqDelegatorAddress,
			}

			res, err := queryClient.DelegatorWithdrawAddress(cmd.Context(), params)
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

func CmdValidatorCommission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-commission [validator-address]",
		Short: "Query validator-commission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryValidatorCommissionRequest{

				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.ValidatorCommission(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdValidatorOutstandingRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-outstanding-rewards [validator-address]",
		Short: "Query validator-outstanding-rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryValidatorOutstandingRewardsRequest{

				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.ValidatorOutstandingRewards(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdValidatorSlashes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-slashes [validator-address]",
		Short: "Query validator-slashes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqValidatorAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryValidatorSlashesRequest{

				ValidatorAddress: reqValidatorAddress,
			}

			res, err := queryClient.ValidatorSlashes(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
