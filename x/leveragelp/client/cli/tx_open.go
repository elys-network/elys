package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open [position] [leverage] [collateral-asset] [collateral-amount] [borrow-asset] [flags]",
		Short:   "Open leveragelp position",
		Example: `elysd tx leveragelp open long 5 uusdc 100000000 uatom --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(5),
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

			msg := types.NewMsgOpen(
				signer.String(),
				argCollateralAsset,
				argCollateralAmount,
				argBorrowAsset,
				argPosition,
				argLeverage,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
