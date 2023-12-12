package cli

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
)

func CmdOpenEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open-estimation [position] [leverage] [trading-asset] [collateral]",
		Short:   "Query open-estimation",
		Example: "elysd q margin open-estimation long 5 uatom 100000000uusdc",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPosition := types.GetPositionFromString(args[0])

			reqLeverage, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			reqTradingAsset := args[2]

			reqCollateral, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			discountStr, err := cmd.Flags().GetString(FlagDiscount)
			if err != nil {
				return err
			}
			discount, err := sdk.NewDecFromStr(discountStr)
			if err != nil {
				return err
			}

			takeProfitPriceStr, err := cmd.Flags().GetString(FlagTakeProfitPrice)
			if err != nil {
				return err
			}

			var takeProfitPrice sdk.Dec
			if takeProfitPriceStr != types.InfinitePriceString {
				takeProfitPrice, err = sdk.NewDecFromStr(takeProfitPriceStr)
				if err != nil {
					return errors.New("invalid take profit price")
				}
			} else {
				takeProfitPrice, err = sdk.NewDecFromStr(types.TakeProfitPriceDefault)
				if err != nil {
					return errors.New("failed to set default take profit price")
				}
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryOpenEstimationRequest{
				Position:        reqPosition,
				Leverage:        reqLeverage,
				TradingAsset:    reqTradingAsset,
				Collateral:      reqCollateral,
				Discount:        discount,
				TakeProfitPrice: takeProfitPrice,
			}

			res, err := queryClient.OpenEstimation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagDiscount, "0.0", "discount to apply to the swap fee")
	cmd.Flags().String(FlagTakeProfitPrice, types.InfinitePriceString, "Optional take profit price")

	return cmd
}
