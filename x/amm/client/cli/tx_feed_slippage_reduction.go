package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdFeedSlippageReduction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-slippage-reduction [pool-id] [slippage-reduction]",
		Short: "Broadcast message feed-slippage-reduction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			poolId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			reduction := sdk.MustNewDecFromStr(args[1])
			msg := types.NewMsgFeedSlippageReduction(
				clientCtx.GetFromAddress().String(),
				uint64(poolId),
				reduction,
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
