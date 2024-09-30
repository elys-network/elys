package cli

import (
	"cosmossdk.io/math"
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAddCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-collateral [amount] [id]",
		Short: "Broadcast message add-collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount, ok := math.NewIntFromString(args[0])
			if !ok {
				return errors.New("invalid collateral amount")
			}
			positionId, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddCollateral(
				clientCtx.GetFromAddress().String(),
				argAmount,
				uint64(positionId),
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
