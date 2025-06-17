package cli

import (
	"errors"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/spf13/cobra"
)

const (
	FlagSwapFee   = "swap-fee"
	FlagUseOracle = "use-oracle"
)

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-pool [weights] [initial-deposit] [fee-denom]",
		Short:   "create a new pool and provide the liquidity to it",
		Example: `elysd tx amm create-pool 100uatom,100uusdc 100000000000uatom,100000000000uusdc --swap-fee=0.00 --use-oracle=false  --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argWeights, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}
			argInitialDeposit, err := sdk.ParseCoinsNormalized(args[1])
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

				// External liquidity ratio initially is set to 1, this value can be changed by feeder only with relevant orderbook liquidity
				// Setting this ratio to 1 is equivalent to considering only amm liquidity when making swaps
				poolAssets = append(poolAssets, types.PoolAsset{
					Weight:                 argWeights[i].Amount,
					Token:                  argInitialDeposit[i],
					ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
				})
			}

			swapFeeStr, err := cmd.Flags().GetString(FlagSwapFee)
			if err != nil {
				return err
			}

			useOracle, err := cmd.Flags().GetBool(FlagUseOracle)
			if err != nil {
				return err
			}

			feeDenom := args[2]

			poolParams := types.PoolParams{
				SwapFee:   sdkmath.LegacyMustNewDecFromStr(swapFeeStr),
				UseOracle: useOracle,
				FeeDenom:  feeDenom,
			}

			msg := types.NewMsgCreatePool(
				clientCtx.GetFromAddress().String(),
				poolParams,
				poolAssets,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagSwapFee, "0.00", "swap fee")
	cmd.Flags().Bool(FlagUseOracle, false, "flag to be an oracle pool or non-oracle pool")

	return cmd
}
