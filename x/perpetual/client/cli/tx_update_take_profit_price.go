package cli

import (
	"errors"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/spf13/cobra"
)

func CmdUpdateTakeProfitPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-take-profit-price [amount] [id] [poolId]",
		Short: "Broadcast message update-take-profit-price",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPrice, err := math.LegacyNewDecFromStr(args[0])
			if err != nil {
				return errors.New("invalid take profit amount")
			}
			positionId, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateTakeProfitPrice(
				clientCtx.GetFromAddress().String(),
				uint64(positionId),
				argPrice,
				poolId,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
