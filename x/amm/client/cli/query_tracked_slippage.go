package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

func CmdTrackedSlippage() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tracked-slippage [pool_id]",
		Short:   "Query Tracked Slippage",
		Example: "elysd q amm tracked-slippage 1",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySlippageTrackRequest{
				PoolId: uint64(poolId),
			}

			res, err := queryClient.SlippageTrack(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdTrackedSlippageAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tracked-slippage-all",
		Short:   "Query All Tracked Slippage",
		Example: "elysd q amm tracked-slippage-all",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySlippageTrackAllRequest{}

			res, err := queryClient.SlippageTrackAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
