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

func CmdAddCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-collateral [id] [collateral]",
		Short: "Broadcast message add-collateral",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPositionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errors.New("invalid position id")
			}

			argAmount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddCollateral(
				clientCtx.GetFromAddress().String(),
				argPositionId,
				argAmount,
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
