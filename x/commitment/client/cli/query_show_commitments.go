package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/spf13/cobra"
)

func CmdShowCommitments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-commitments [creator]",
		Short: "Query show-commitments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryShowCommitmentsRequest{
				Creator: args[0],
			}
			res, err := queryClient.ShowCommitments(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdCommittedTokensLocked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "committed-tokens-locked [address]",
		Short: "Show locked coins in commitment not unlockable",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryCommittedTokensLockedRequest{
				Address: args[0],
			}
			res, err := queryClient.CommittedTokensLocked(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdNumberOfCommitments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "number-of-commitments",
		Short: "Query number-of-commitments",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryNumberOfCommitmentsRequest{}
			res, err := queryClient.NumberOfCommitments(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
