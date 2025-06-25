package cli

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/commitment/types"
	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CmdUncommitTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uncommit-tokens [amount] [denom]",
		Short: "Broadcast message uncommit-tokens",
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

			msg := types.NewMsgUncommitTokens(
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
