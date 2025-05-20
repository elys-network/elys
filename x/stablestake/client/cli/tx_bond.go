package cli

import (
	"fmt"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v4/x/stablestake/types"
	"github.com/spf13/cobra"
)

func CmdBond() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bond [amount] [pool-id]",
		Short: "Broadcast message bond",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			amount, ok := math.NewIntFromString(args[0])
			if !ok {
				return fmt.Errorf("unable to parse bonding amount")
			}

			poolId, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("unable to parse pool id")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgBond(
				clientCtx.GetFromAddress().String(),
				amount,
				poolId.Uint64(),
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
