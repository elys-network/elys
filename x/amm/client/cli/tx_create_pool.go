package cli

import (
	sdkmath "cosmossdk.io/math"
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

const (
	FlagSwapFee                     = "swap-fee"
	FlagExitFee                     = "exit-fee"
	FlagUseOracle                   = "use-oracle"
	FlagWeightBreakingFeeMultiplier = "weight-breaking-fee-multiplier"
	FlagWeightBreakingFeeExponent   = "weight-breaking-fee-exponent"
	FlagWeightRecoveryFeePortion    = "weight-recovery-fee-portion"
	FlagWeightBreakingFeePortion    = "weight-breaking-fee-portion"
	FlagThresholdWeightDifference   = "threshold-weight-diff"
	FlagFeeDenom                    = "fee-denom"
)

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-pool [weights] [initial-deposit]",
		Short:   "create a new pool and provide the liquidity to it",
		Example: `elysd tx amm create-pool 100uatom,100uusdc 100000000000uatom,100000000000uusdc --swap-fee=0.00 --exit-fee=0.00 --use-oracle=false  --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000`,
		Args:    cobra.ExactArgs(2),
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
					ExternalLiquidityRatio: sdk.OneDec(),
				})
			}

			swapFeeStr, err := cmd.Flags().GetString(FlagSwapFee)
			if err != nil {
				return err
			}

			exitFeeStr, err := cmd.Flags().GetString(FlagExitFee)
			if err != nil {
				return err
			}

			useOracle, err := cmd.Flags().GetBool(FlagUseOracle)
			if err != nil {
				return err
			}

			weightBreakingFeeMultiplierStr, err := cmd.Flags().GetString(FlagWeightBreakingFeeMultiplier)
			if err != nil {
				return err
			}

			weightBreakingFeeExponentStr, err := cmd.Flags().GetString(FlagWeightBreakingFeeExponent)
			if err != nil {
				return err
			}

			weightRecoveryFeePortionStr, err := cmd.Flags().GetString(FlagWeightRecoveryFeePortion)
			if err != nil {
				return err
			}

			weightBreakingFeePortionStr, err := cmd.Flags().GetString(FlagWeightBreakingFeePortion)
			if err != nil {
				return err
			}

			thresholdWeightDifferenceStr, err := cmd.Flags().GetString(FlagThresholdWeightDifference)
			if err != nil {
				return err
			}

			feeDenom, err := cmd.Flags().GetString(FlagFeeDenom)
			if err != nil {
				return err
			}

			poolParams := &types.PoolParams{
				SwapFee:                     sdk.MustNewDecFromStr(swapFeeStr),
				ExitFee:                     sdk.MustNewDecFromStr(exitFeeStr),
				UseOracle:                   useOracle,
				WeightBreakingFeeMultiplier: sdk.MustNewDecFromStr(weightBreakingFeeMultiplierStr),
				WeightBreakingFeeExponent:   sdk.MustNewDecFromStr(weightBreakingFeeExponentStr),
				WeightRecoveryFeePortion:    sdk.MustNewDecFromStr(weightRecoveryFeePortionStr),
				WeightBreakingFeePortion:    sdkmath.LegacyMustNewDecFromStr(weightBreakingFeePortionStr),
				ThresholdWeightDifference:   sdk.MustNewDecFromStr(thresholdWeightDifferenceStr),
				FeeDenom:                    feeDenom,
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

	cmd.Flags().String(FlagSwapFee, "0.00", "swap fee")
	cmd.Flags().String(FlagExitFee, "0.00", "exit fee")
	cmd.Flags().Bool(FlagUseOracle, false, "flag to be an oracle pool or non-oracle pool")
	cmd.Flags().String(FlagWeightBreakingFeeMultiplier, "0.00", "weight breaking fee multiplier")
	cmd.Flags().String(FlagWeightBreakingFeeExponent, "2.50", "weight breaking fee exponent")
	cmd.Flags().String(FlagWeightRecoveryFeePortion, "0.10", "weight recovery fee portion")
	cmd.Flags().String(FlagWeightBreakingFeePortion, "0.10", "weight breaking fee portion")
	cmd.Flags().String(FlagThresholdWeightDifference, "0.00", "threshold weight difference - valid for oracle pool")
	cmd.Flags().String(FlagFeeDenom, "", "fee denom")

	return cmd
}
