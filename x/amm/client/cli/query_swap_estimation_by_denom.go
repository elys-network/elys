package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSwapEstimationByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "swap-estimation-by-denom [amount] [denom-in] [denom-out]",
		Short:   "Query swap-estimation-by-denom",
		Example: "elysd q amm swap-estimation-by-denom 100uatom uatom uosmo",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAmount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			reqDenomIn := args[1]
			reqDenomOut := args[2]

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

			params := &types.QuerySwapEstimationByDenomRequest{
				Amount:   reqAmount,
				DenomIn:  reqDenomIn,
				DenomOut: reqDenomOut,
				Discount: discount,
			}

			res, err := queryClient.SwapEstimationByDenom(cmd.Context(), params)
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
