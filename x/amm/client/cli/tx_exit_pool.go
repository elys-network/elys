package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdExitPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exit-pool [pool-id] [min-amounts-out] [share-amount-in]",
		Short:   "exit a new pool and withdraw the liquidity from it",
		Example: `elysd tx amm exit-pool 0 1000uatom,1000uusdc 200000000000000000 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argMinAmountsOut, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			argShareAmountIn, ok := sdkmath.NewIntFromString(args[2])
			if !ok {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgExitPool(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				argMinAmountsOut,
				argShareAmountIn,
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
