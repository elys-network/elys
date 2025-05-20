package cli

import (
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v4/x/stablestake/types"
	"github.com/spf13/cobra"
)

func CmdAddPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-pool [deposit-denom] [interest-rate] [interest-rate-max] [interest-rate-min] [interest-rate-increase] [interest-rate-decrease] [health-factor] [max-leverage-ratio] [max-withdraw-ratio]",
		Short: "Broadcast message add-pool",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			depositDenom := args[0]

			interestRate, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			interestRateMax, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}

			interestRateMin, err := math.LegacyNewDecFromStr(args[3])
			if err != nil {
				return err
			}

			interestRateIncrease, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return err
			}

			interestRateDecrease, err := math.LegacyNewDecFromStr(args[5])
			if err != nil {
				return err
			}

			healthFactor, err := math.LegacyNewDecFromStr(args[6])
			if err != nil {
				return err
			}

			maxLeverageRatio, err := math.LegacyNewDecFromStr(args[7])
			if err != nil {
				return err
			}

			maxWithdrawRatio, err := math.LegacyNewDecFromStr(args[8])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddPool(
				clientCtx.GetFromAddress().String(),
				depositDenom,
				interestRate,
				interestRateMax,
				interestRateMin,
				interestRateIncrease,
				interestRateDecrease,
				healthFactor,
				maxLeverageRatio,
				maxWithdrawRatio,
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
