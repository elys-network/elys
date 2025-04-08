package cli

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/clob/types"
	"github.com/spf13/cobra"
	"strconv"
)

func CmdCreateMarket() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-market [base-denom] [quote-denom] [initial-margin] [maintenance-margin] [maker-fee] [taker-fee] [liquidation-fee] [min-price-tick] [min-quantity-tick] [min-notional] [max-funding-rate] [max-funding-rate-change] [twap-time]",
		Short:   "opens new perpetual market",
		Example: `elysd tx clob exit-pool create-market uatom uusdc 0.02 0.02 0 0 0.01 0.001 1 0.05 0.001 100 --from=bob --gas=1000000`,
		Args:    cobra.ExactArgs(13),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			initialMargin, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}
			maintenanceMargin, err := math.LegacyNewDecFromStr(args[3])
			if err != nil {
				return err
			}
			makerFee, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return err
			}
			takerFee, err := math.LegacyNewDecFromStr(args[5])
			if err != nil {
				return err
			}
			liquidationFee, err := math.LegacyNewDecFromStr(args[6])
			if err != nil {
				return err
			}
			minPriceTick, err := math.LegacyNewDecFromStr(args[7])
			if err != nil {
				return err
			}
			minQuantityTick, err := math.LegacyNewDecFromStr(args[8])
			if err != nil {
				return err
			}
			minNotional, err := math.LegacyNewDecFromStr(args[9])
			if err != nil {
				return err
			}
			maxFunding, err := math.LegacyNewDecFromStr(args[10])
			if err != nil {
				return err
			}
			maxFundingRateChange, err := math.LegacyNewDecFromStr(args[11])
			if err != nil {
				return err
			}
			twapTime, err := strconv.ParseUint(args[12], 10, 64)
			if err != nil {
				return err
			}

			msg := types.MsgCreatPerpetualMarket{
				Creator:                 clientCtx.GetFromAddress().String(),
				BaseDenom:               args[0],
				QuoteDenom:              args[1],
				InitialMarginRatio:      initialMargin,
				MaintenanceMarginRatio:  maintenanceMargin,
				MakerFeeRate:            makerFee,
				TakerFeeRate:            takerFee,
				LiquidationFeeShareRate: liquidationFee,
				MinPriceTickSize:        minPriceTick,
				MinQuantityTickSize:     minQuantityTick,
				MinNotional:             minNotional,
				MaxFundingRate:          maxFunding,
				MaxFundingRateChange:    maxFundingRateChange,
				MaxTwapPricesTime:       twapTime,
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
