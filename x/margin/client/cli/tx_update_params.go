package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ = strconv.Itoa(0)

// Governance command
func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params",
		Short: "Update margin params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			summary, err := cmd.Flags().GetString(cli.FlagSummary)
			if err != nil {
				return err
			}

			metadata, err := cmd.Flags().GetString(cli.FlagMetadata)
			if err != nil {
				return err
			}

			leverage_max := sdk.MustNewDecFromStr(viper.GetString("leverage-max"))

			if leverage_max.GT(sdk.NewDec(10)) || leverage_max.LT(sdk.NewDec(1)) {
				return errors.New("invalid leverage max, it has to be between 1-10.")
			}

			params := &types.Params{
				LeverageMax:                              leverage_max,
				InterestRateMax:                          sdk.MustNewDecFromStr(viper.GetString("interest-rate-max")),
				InterestRateMin:                          sdk.MustNewDecFromStr(viper.GetString("interest-rate-min")),
				InterestRateIncrease:                     sdk.MustNewDecFromStr(viper.GetString("interest-rate-increase")),
				InterestRateDecrease:                     sdk.MustNewDecFromStr(viper.GetString("interest-rate-decrease")),
				HealthGainFactor:                         sdk.MustNewDecFromStr(viper.GetString("health-gain-factor")),
				EpochLength:                              viper.GetInt64("epoch-length"),
				RemovalQueueThreshold:                    sdk.MustNewDecFromStr(viper.GetString("removal-queue-threshold")),
				MaxOpenPositions:                         viper.GetInt64("max-open-positions"),
				PoolOpenThreshold:                        sdk.MustNewDecFromStr(viper.GetString("pool-open-threshold")),
				ForceCloseFundPercentage:                 sdk.MustNewDecFromStr(viper.GetString("force-close-fund-percentage")),
				ForceCloseFundAddress:                    viper.GetString("force-close-fund-address"),
				IncrementalInterestPaymentFundPercentage: sdk.MustNewDecFromStr(viper.GetString("incremental-interest-payment-fund-percentage")),
				IncrementalInterestPaymentFundAddress:    viper.GetString("incremental-interest-payment-fund-address"),
				SqModifier:                               sdk.MustNewDecFromStr(viper.GetString("sq-modifier")),
				SafetyFactor:                             sdk.MustNewDecFromStr(viper.GetString("safety-factor")),
				IncrementalInterestPaymentEnabled:        viper.GetBool("incremental-interest-payment-enabled"),
				WhitelistingEnabled:                      viper.GetBool("whitelisting-enabled"),
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			govAddress := sdk.AccAddress(address.Module("gov"))
			msg := types.NewMsgUpdateParams(
				govAddress.String(),
				params,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			govMsg, err := v1.NewMsgSubmitProposal([]sdk.Msg{msg}, deposit, signer.String(), metadata, title, summary)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), govMsg)
		},
	}

	cmd.Flags().String("leverage-max", "", "max leverage (integer)")
	cmd.Flags().String("interest-rate-max", "", "max interest rate (decimal)")
	cmd.Flags().String("interest-rate-min", "", "min interest rate (decimal)")
	cmd.Flags().String("interest-rate-increase", "", "interest rate increase (decimal)")
	cmd.Flags().String("interest-rate-decrease", "", "interest rate decrease (decimal)")
	cmd.Flags().String("health-gain-factor", "", "health gain factor (decimal)")
	cmd.Flags().Int64("epoch-length", 1, "epoch length in blocks (integer)")
	cmd.Flags().Uint64("max-open-positions", 10000, "max open positions")
	cmd.Flags().String("removal-queue-threshold", "", "removal queue threshold (decimal range 0-1)")
	cmd.Flags().String("pool-open-threshold", "", "threshold to prevent new positions (decimal range 0-1)")
	cmd.Flags().String("force-close-fund-percentage", "", "percentage of force close proceeds for fund (decimal range 0-1)")
	cmd.Flags().String("force-close-fund-address", "", "address of fund wallet for force close")
	cmd.Flags().Bool("incremental-interest-payment-enabled", true, "enable incremental interest payment")
	cmd.Flags().String("incremental-interest-payment-fund-percentage", "", "percentage of incremental interest payment proceeds for fund (decimal range 0-1)")
	cmd.Flags().String("incremental-interest-payment-fund-address", "", "address of fund wallet for incremental interest payment")
	cmd.Flags().String("sq-modifier", "", "the modifier value for the removal queue's sq formula")
	cmd.Flags().String("safety-factor", "", "the safety factor used in liquidation ratio")
	cmd.Flags().Bool("whitelisting-enabled", false, "Enable whitelisting")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired("leverage-max")
	_ = cmd.MarkFlagRequired("interest-rate-max")
	_ = cmd.MarkFlagRequired("interest-rate-min")
	_ = cmd.MarkFlagRequired("interest-rate-increase")
	_ = cmd.MarkFlagRequired("interest-rate-decrease")
	_ = cmd.MarkFlagRequired("health-gain-factor")
	_ = cmd.MarkFlagRequired("removal-queue-threshold")
	_ = cmd.MarkFlagRequired("max-open-positions")
	_ = cmd.MarkFlagRequired("pool-open-threshold")
	_ = cmd.MarkFlagRequired("force-close-fund-percentage")
	_ = cmd.MarkFlagRequired("force-close-fund-address")
	_ = cmd.MarkFlagRequired("incremental-interest-payment-enabled")
	_ = cmd.MarkFlagRequired("incremental-interest-payment-fund-percentage")
	_ = cmd.MarkFlagRequired("incremental-interest-payment-fund-address")
	_ = cmd.MarkFlagRequired("sq-modifier")
	_ = cmd.MarkFlagRequired("safety-factor")
	_ = cmd.MarkFlagRequired("whitelisting-enabled")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
