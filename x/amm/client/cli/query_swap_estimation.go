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

func CmdSwapEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swap-estimation [token-in] {pool_id token_out_denom}...",
		Short:   "Query SwapEstimation",
		Example: "elysd q amm swap-estimation 100token 1 token_out1 2 token_out2 ...",
		Args:    cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqTokenIn, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			if (len(args)-1)%2 != 0 {
				return errors.New("you must provide pairs of pool_id and token_out_denom for routes")
			}

			var reqRoutes []*types.SwapAmountInRoute
			for i := 1; i+1 < len(args); i += 2 {
				poolID, err := strconv.ParseUint(args[i], 10, 64)
				if err != nil {
					return err
				}

				reqRoutes = append(reqRoutes, &types.SwapAmountInRoute{
					PoolId:        poolID,
					TokenOutDenom: args[i+1],
				})
			}

			discountStr, err := cmd.Flags().GetString(FlagDiscount)
			if err != nil {
				return err
			}
			discount, err := sdk.NewDecFromStr(discountStr)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySwapEstimationRequest{
				Routes:   reqRoutes,
				TokenIn:  reqTokenIn,
				Discount: discount,
			}

			res, err := queryClient.SwapEstimation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagDiscount, "0.0", "discount to apply to the swap fee")

	return cmd
}
