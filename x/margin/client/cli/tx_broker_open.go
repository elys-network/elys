package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdBrokerOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broker-open [position] [leverage] [collateral-asset] [collateral-amount] [borrow-asset] [flags]",
		Short: "Open margin position as a broker",
		Example: `Infinte profitability:
elysd tx margin broker-open long 5 uusdc 100000000 uatom sif123 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
Finite profitability:
elysd tx margin broker-open short 5 uusdc 100000000 uatom sif123 --take-profit 100 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args: cobra.ExactArgs(6),
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

			argLeverage, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			argCollateralAsset := args[2]

			argCollateralAmount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("invalid collateral amount")
			}

			argBorrowAsset := args[4]

			argMtpOwner := args[5]
			if argMtpOwner == "" {
				return errors.New("invalid mtp owner address")
			}

			takeProfitPriceStr, err := cmd.Flags().GetString(flagTakeProfitPrice)
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

			msg := types.NewMsgBrokerOpen(
				signer.String(),
				argCollateralAsset,
				argCollateralAmount,
				argBorrowAsset,
				argPosition,
				argLeverage,
				takeProfitPrice,
				argMtpOwner,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagTakeProfitPrice, types.InfinitePriceString, "Optional take profit price")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
