package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

const FlagFeeDenom = "fee-denom"

func CmdUpdatePoolParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-pool-params [pool-id] [flags]",
		Short:   "Update pool params",
		Example: "elysd tx amm update-pool-params 1 --swap-fee=0.00 --exit-fee=0.00 --use-oracle=false --from=bob --yes --gas=1000000",
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

			useOracle, err := cmd.Flags().GetBool(FlagUseOracle)
			if err != nil {
				return err
			}

			feeDenom, err := cmd.Flags().GetString(FlagFeeDenom)
			if err != nil {
				return err
			}

			poolParams := types.PoolParams{
				SwapFee:   sdkmath.LegacyMustNewDecFromStr(swapFeeStr),
				UseOracle: useOracle,
				FeeDenom:  feeDenom,
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
	cmd.Flags().Bool(FlagUseOracle, false, "flag to be an oracle pool or non-oracle pool")
	cmd.Flags().String(FlagFeeDenom, "", "fee denom")

	return cmd
}
