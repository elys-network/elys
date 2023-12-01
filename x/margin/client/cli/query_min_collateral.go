package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMinCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "min-collateral [position] [leverage] [trading-asset] [collateral-asset]",
		Short:   "Query min-collateral",
		Example: "elysd q margin min-collateral long 5 uatom uusdc",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPosition := types.GetPositionFromString(args[0])

			reqLeverage, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			reqTradingAsset := args[2]
			reqCollateralAsset := args[3]

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

			params := &types.QueryMinCollateralRequest{
				Position:        reqPosition,
				Leverage:        reqLeverage,
				TradingAsset:    reqTradingAsset,
				CollateralAsset: reqCollateralAsset,
				Discount:        discount,
			}

			res, err := queryClient.MinCollateral(cmd.Context(), params)
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
