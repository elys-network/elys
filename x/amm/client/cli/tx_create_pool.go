package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-pool [weights] [initial-deposit] [swap-fee] [exit-fee]",
		Short:   "create a new pool and provide the liquidity to it",
		Example: `elysd tx amm create-pool 100uatom,100uusdc 100000000000uatom,100000000000uusdc 0.00 0.00  --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argWeights, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}
			argInitialDeposit, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			argSwapFee, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}
			argExitFee, err := sdk.NewDecFromStr(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(argInitialDeposit) != len(argWeights) {
				return errors.New("deposit tokens and token weights should have same length")
			}

			var poolAssets []types.PoolAsset
			for i := 0; i < len(argWeights); i++ {
				if argWeights[i].Denom != argInitialDeposit[i].Denom {
					return errors.New("deposit tokens and token weights should have same denom order")
				}

				poolAssets = append(poolAssets, types.PoolAsset{
					Weight: argWeights[i].Amount,
					Token:  argInitialDeposit[i],
				})
			}

			poolParams := &types.PoolParams{
				SwapFee: argSwapFee,
				ExitFee: argExitFee,
			}

			msg := types.NewMsgCreatePool(
				clientCtx.GetFromAddress().String(),
				poolParams,
				poolAssets,
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
