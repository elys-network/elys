package cli

import (
	"errors"
	"strconv"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdOpenEstimationByFinal() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open-estimation-by-final [position] [leverage] [final-amount] [pool-id] [collateral-denom]",
		Short:   "Query open-estimation-by-final amount",
		Example: "elysd q perpetual open-estimation-by-final long 5 100000000uusdc 1 uatom",
		Args:    cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPosition := types.GetPositionFromString(args[0])

			reqLeverage, err := sdkmath.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			finalAmount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			reqPoolId, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			collateralDenom := args[4]
			if err = sdk.ValidateDenom(collateralDenom); err != nil {
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

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryOpenEstimationByFinalRequest{
				Position:        reqPosition,
				Leverage:        reqLeverage,
				FinalAmount:     finalAmount,
				TakeProfitPrice: takeProfitPrice,
				PoolId:          reqPoolId,
				CollateralDenom: collateralDenom,
				Address:         address,
			}

			res, err := queryClient.OpenEstimationByFinal(cmd.Context(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagTakeProfitPrice, types.InfinitePriceString, "Optional take profit price")
	cmd.Flags().String(FlagAddress, "", "address of the account which will open the position")

	return cmd
}
