package cli

import (
	"errors"
	"strconv"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdBrokerOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broker-open [position] [leverage] [pool-id] [trading-asset] [collateral] [owner] [flags]",
		Short: "Broker Open perpetual position",
		Example: `Infinte profitability:
elysd tx perpetual broker-open long 5 1 uatom 100000000uusdc elys1w9uac4zrf9z7qd604qxk2y4n74568lfl8vutz4 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
Finite profitability:
elysd tx perpetual broker-open short 5 1 uatom 100000000uusdc elys1w9uac4zrf9z7qd604qxk2y4n74568lfl8vutz4 --take-profit 100 --stop-loss 10 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			_, err = sdk.AccAddressFromBech32(args[5])
			if err != nil {
				return errors.New("invalid owner address")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			argPosition := types.GetPositionFromString(args[0])

			argLeverage, err := sdkmath.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			argPoolId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			argTradingAsset := args[3]

			argCollateral, err := sdk.ParseCoinNormalized(args[4])
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

			stopLossPriceStr, err := cmd.Flags().GetString(FlagStopLossPrice)
			if err != nil {
				return err
			}

			var stopLossPrice sdkmath.LegacyDec
			if stopLossPriceStr != types.ZeroPriceString {
				stopLossPrice, err = sdkmath.LegacyNewDecFromStr(stopLossPriceStr)
				if err != nil {
					return errors.New("invalid stop loss price")
				}
			} else {
				stopLossPrice = types.StopLossPriceDefault
			}

			msg := types.NewMsgBrokerOpen(
				signer.String(),
				argPosition,
				argLeverage,
				argTradingAsset,
				argCollateral,
				takeProfitPrice,
				args[5],
				stopLossPrice,
				argPoolId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTakeProfitPrice, types.InfinitePriceString, "Optional take profit price")
	cmd.Flags().String(FlagStopLossPrice, types.ZeroPriceString, "Optional stop loss price")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
