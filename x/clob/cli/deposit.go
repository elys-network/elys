package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
	"github.com/spf13/cobra"
	"strconv"
)

func CmdDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deposit [sub-account-id] [amount]",
		Short:   "exit a new pool and withdraw the liquidity from it",
		Example: `elysd tx amm exit-pool 0 1000uatom,1000uusdc 200000000000000000 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			subAccountId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgDeposit{
				Sender:       clientCtx.GetFromAddress().String(),
				SubAccountId: subAccountId,
				Coin:         amount,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
