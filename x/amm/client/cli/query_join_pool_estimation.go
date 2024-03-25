package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

func CmdJoinPoolEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "join-pool-estimation [pool_id] [tokens-in]",
		Short:   "Query JoinPoolEstimation",
		Example: "elysd q amm join-pool-estimation 1 100token,100token2",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			reqTokensIn, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryJoinPoolEstimationRequest{
				PoolId:    uint64(poolId),
				AmountsIn: reqTokensIn,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.JoinPoolEstimation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
