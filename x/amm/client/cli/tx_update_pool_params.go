package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdUpdatePoolParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-pool-params [pool-id] [flags]",
		Short:   "Update pool params",
		Example: "elysd tx amm update-pool-params 1 --swap-fee=0.00 --exit-fee=0.00 --use-oracle=false --from=jack --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
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

			externalLiquidityRatioStr, err := cmd.Flags().GetString(FlagExternalLiquidityRatio)
			if err != nil {
				return err
			}

			lpFeePortionStr, err := cmd.Flags().GetString(FlagLpFeePortion)
			if err != nil {
				return err
			}

			stakingFeePortionStr, err := cmd.Flags().GetString(FlagStakingFeePortion)
			if err != nil {
				return err
			}

			weightRecoveryFeePortionStr, err := cmd.Flags().GetString(FlagWeightRecoveryFeePortion)
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
				ExternalLiquidityRatio:      sdk.MustNewDecFromStr(externalLiquidityRatioStr),
				LpFeePortion:                sdk.MustNewDecFromStr(lpFeePortionStr),
				StakingFeePortion:           sdk.MustNewDecFromStr(stakingFeePortionStr),
				WeightRecoveryFeePortion:    sdk.MustNewDecFromStr(weightRecoveryFeePortionStr),
				ThresholdWeightDifference:   sdk.MustNewDecFromStr(thresholdWeightDifferenceStr),
				FeeDenom:                    feeDenom,
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePoolParams(
				clientCtx.GetFromAddress().String(),
				argPoolId,
				poolParams,
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
	cmd.Flags().String(FlagExternalLiquidityRatio, "0.00", "external liquidity ratio - valid for oracle pools")
	cmd.Flags().String(FlagLpFeePortion, "0.00", "lp fee portion")
	cmd.Flags().String(FlagStakingFeePortion, "0.00", "staking fee portion")
	cmd.Flags().String(FlagWeightRecoveryFeePortion, "0.00", "weight recovery fee portion")
	cmd.Flags().String(FlagThresholdWeightDifference, "0.00", "threshold weight difference - valid for oracle pool")
	cmd.Flags().String(FlagFeeDenom, "", "fee denom")

	return cmd
}
