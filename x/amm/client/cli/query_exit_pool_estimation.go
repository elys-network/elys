package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

func CmdExitPoolEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exit-pool-estimation [pool_id] [token-in] [token_out_denom]",
		Short:   "Query ExitPoolEstimation",
		Example: "elysd q amm exit-pool-estimation 1 10000 token",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			shareAmountIn, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid share in amount")
			}

			tokenOutDenom := args[2]

			params := &types.QueryExitPoolEstimationRequest{
				PoolId:        uint64(poolId),
				ShareAmountIn: shareAmountIn,
				TokenOutDenom: tokenOutDenom,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ExitPoolEstimation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
