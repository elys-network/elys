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
		Use:     "open [leverage] [collateral-asset] [collateral-amount] [amm-pool-id] [flags]",
		Short:   "Open leveragelp position",
		Example: `elysd tx leveragelp open 5 uusdc 100000000 1 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
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

			leverage, err := sdk.NewDecFromStr(args[0])
			if err != nil {
				return err
			}

			collateralAsset := args[1]
			collateralAmount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return errors.New("invalid collateral amount")
			}

			ammPoolId, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgOpen(
				signer.String(),
				collateralAsset,
				collateralAmount,
				uint64(ammPoolId),
				leverage,
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
