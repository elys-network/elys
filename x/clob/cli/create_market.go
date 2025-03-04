package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/clob/types"
	"github.com/spf13/cobra"
)

func CmdCreateMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-market [base-denom] [quote-denom]",
		Short:   "exit a new pool and withdraw the liquidity from it",
		Example: `elysd tx amm exit-pool 0 1000uatom,1000uusdc 200000000000000000 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreatPerpetualMarket{
				Creator:                clientCtx.GetFromAddress().String(),
				BaseDenom:              args[0],
				QuoteDenom:             args[1],
				InitialMarginRatio:     sdkmath.LegacyMustNewDecFromStr("0.2"),
				MaintenanceMarginRatio: sdkmath.LegacyMustNewDecFromStr("0.2"),
				MakerFeeRate:           sdkmath.LegacyZeroDec(),
				TakerFeeRate:           sdkmath.LegacyZeroDec(),
				RelayerFeeShareRate:    sdkmath.LegacyZeroDec(),
				MinPriceTickSize:       sdkmath.LegacyZeroDec(),
				MinQuantityTickSize:    sdkmath.OneInt(),
				MinNotional:            sdkmath.LegacyZeroDec(),
				AllowedCollateral:      []string{"uusdc"},
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
