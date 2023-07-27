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

func CmdOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open margin position",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			collateralAsset, err := cmd.Flags().GetString("collateral_asset")
			if err != nil {
				return err
			}

			collateralAmount, err := cmd.Flags().GetString("collateral_amount")
			if err != nil {
				return err
			}

			borrowAsset, err := cmd.Flags().GetString("borrow_asset")
			if err != nil {
				return err
			}

			position, err := cmd.Flags().GetString("position")
			if err != nil {
				return err
			}
			positionEnum := types.GetPositionFromString(position)

			leverage, err := cmd.Flags().GetString("leverage")
			if err != nil {
				return err
			}

			leverageDec, err := sdk.NewDecFromStr(leverage)
			if err != nil {
				return err
			}

			collateralAmt, ok := sdk.NewIntFromString(collateralAmount)
			if !ok {
				return errors.New("invalid collateral amount")
			}

			msg := types.NewMsgOpen(
				clientCtx.GetFromAddress().String(),
				collateralAsset,
				collateralAmt,
				borrowAsset,
				positionEnum,
				leverageDec,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("collateral_amount", "0", "amount of collateral asset")
	cmd.Flags().String("collateral_asset", "", "symbol of asset")
	cmd.Flags().String("borrow_asset", "", "symbol of asset")
	cmd.Flags().String("position", "", "type of position")
	cmd.Flags().String("leverage", "", "leverage of position")
	_ = cmd.MarkFlagRequired("collateral_amount")
	_ = cmd.MarkFlagRequired("collateral_asset")
	_ = cmd.MarkFlagRequired("borrow_asset")
	_ = cmd.MarkFlagRequired("position")
	_ = cmd.MarkFlagRequired("leverage")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
