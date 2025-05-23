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

func CmdStake() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake [amount] [asset] [validator-address]",
		Short: "Stake Elys tokens",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount, found := math.NewIntFromString(args[0])
			if !found {
				return errorsmod.Wrap(sdkerrors.ErrInvalidType, "cannot convert string to int")
			}
			argAsset := args[1]
			argValidatorAddress := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgStake(
				clientCtx.GetFromAddress().String(),
				argAmount,
				argAsset,
				argValidatorAddress,
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
