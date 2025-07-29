package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdAddCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-collateral [mtp-id] [collateral] [pool-id]",
		Short: "Add collateral perpetual position",
		Example: `
elysd tx add-collateral 11 100000000uusdc 1 --from=bob --yes --gas=1000000`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			mtpID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.MsgAddCollateral{
				Creator:       signer.String(),
				Id:            mtpID,
				AddCollateral: collateral,
				PoolId:        poolId,
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(FlagTakeProfitPrice, types.InfinitePriceString, "Optional take profit price")
	cmd.Flags().String(FlagStopLossPrice, types.ZeroPriceString, "Optional stop loss price")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
