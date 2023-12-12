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
)

func CmdJoinPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "join-pool [pool-id] [max-amounts-in] [share-amount-out]",
		Short:   "join a new pool and provide the liquidity to it",
		Example: `elysd tx amm join-pool 0 2000uatom,2000uusdc 200000000000000000 true --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			poolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			maxAmountsIn, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			shareAmountOut, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			noRemaining, err := strconv.ParseBool(args[3])

			msg := types.NewMsgJoinPool(
				clientCtx.GetFromAddress().String(),
				poolId,
				maxAmountsIn,
				shareAmountOut,
				noRemaining,
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
