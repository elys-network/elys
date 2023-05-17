package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

var _ = strconv.Itoa(0)

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [weights] [initial-deposit] [swap-fee] [exit-fee] [future-governor] [scalling-factors]",
		Short: "create a new pool and provide the liquidity to it",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argWeights, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}
			argInitialDeposit, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			argSwapFee, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argExitFee, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argFutureGovernor := args[4]
			argCastScallingFactors := strings.Split(args[5], listSeparator)
			argScallingFactors := make([]uint64, len(argCastScallingFactors))
			for i, arg := range argCastScallingFactors {
				value, err := cast.ToUint64E(arg)
				if err != nil {
					return err
				}
				argScallingFactors[i] = value
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePool(
				clientCtx.GetFromAddress().String(),
				argWeights,
				argInitialDeposit,
				argSwapFee,
				argExitFee,
				argFutureGovernor,
				argScallingFactors,
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
