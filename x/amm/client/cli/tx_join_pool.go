package cli

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdJoinPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "join-pool [pool-id] [max-amounts-in] [share-amount-out]",
		Short:   "join a new pool and provide the liquidity to it",
		Example: `elysd tx amm join-pool 0 2000uatom,2000uusdc 200000000000000000 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			poolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			maxAmountsIn, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			shareAmountOut, ok := math.NewIntFromString(args[2])
			if !ok {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgJoinPool(
				clientCtx.GetFromAddress().String(),
				poolId,
				maxAmountsIn,
				shareAmountOut,
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
