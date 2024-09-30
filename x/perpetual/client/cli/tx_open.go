package cli

import (
	sdkmath "cosmossdk.io/math"
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open [position] [leverage] [trading-asset] [collateral] [stop-loss-price] [flags]",
		Short: "Open perpetual position",
		Example: `Infinte profitability:
elysd tx perpetual open long 5 uatom 100000000uusdc 100.0 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
Finite profitability:
elysd tx perpetual open short 5 uatom 100000000uusdc 100.0 --take-profit 100 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
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

			argTradingAsset := args[2]

			argCollateral, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			takeProfitPriceStr, err := cmd.Flags().GetString(FlagTakeProfitPrice)
			if err != nil {
				return err
			}

			stopLossPrice, err := sdkmath.LegacyNewDecFromStr(args[4])
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
				takeProfitPrice, err = sdkmath.LegacyNewDecFromStr(types.TakeProfitPriceDefault)
				if err != nil {
					return errors.New("failed to set default take profit price")
				}
			}

			msg := types.NewMsgOpen(
				signer.String(),
				argPosition,
				argLeverage,
				argTradingAsset,
				argCollateral,
				takeProfitPrice,
				stopLossPrice,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTakeProfitPrice, types.InfinitePriceString, "Optional take profit price")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
