package cli

import (
	"errors"
	"strconv"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdOpenEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open-estimation [position] [leverage] [trading-asset] [collateral] [pool-id]",
		Short:   "Query open-estimation",
		Example: "elysd q perpetual open-estimation long 5 uatom 100000000uusdc 1",
		Args:    cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPosition := types.GetPositionFromString(args[0])

			reqLeverage, err := sdkmath.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			reqTradingAsset := args[2]

			reqCollateral, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			reqPoolId, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			takeProfitPriceStr, err := cmd.Flags().GetString(FlagTakeProfitPrice)
			if err != nil {
				return err
			}

			var takeProfitPrice sdkmath.LegacyDec
			if takeProfitPriceStr != types.InfinitePriceString {
				takeProfitPrice, err = sdkmath.LegacyNewDecFromStr(takeProfitPriceStr)
				if err != nil {
					return errors.New("invalid take profit price")
				}
			} else {
				takeProfitPrice = types.TakeProfitPriceDefault
			}

			address, err := cmd.Flags().GetString(FlagAddress)
			if err != nil {
				return err
			}

			limitPriceStr, err := cmd.Flags().GetString(FlagLimitPrice)
			if err != nil {
				return err
			}

			limitPrice, err := sdkmath.LegacyNewDecFromStr(limitPriceStr)
			if err != nil {
				return errors.New("invalid limit price")
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
				TakeProfitPrice: takeProfitPrice,
				PoolId:          reqPoolId,
				LimitPrice:      limitPrice,
				Address:         address,
			}

			res, err := queryClient.OpenEstimation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagTakeProfitPrice, types.InfinitePriceString, "Optional take profit price")
	cmd.Flags().String(FlagLimitPrice, "0.0", "limit price, default 0 which calculates at market price")
	cmd.Flags().String(FlagAddress, "", "address of the account which will open the position")

	return cmd
}
