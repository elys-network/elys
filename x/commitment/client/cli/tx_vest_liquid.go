package cli

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/x/commitment/types"
	"github.com/spf13/cobra"
)

func CmdVestLiquid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vest-liquid [amount] [denom]",
		Short: "Broadcast message vest-liquid",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount, found := math.NewIntFromString(args[0])
			if !found {
				return errorsmod.Wrap(sdkerrors.ErrInvalidType, "cannot convert string to int")
			}
			argDenom := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgVestLiquid(
				clientCtx.GetFromAddress().String(),
				argAmount,
				argDenom,
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
