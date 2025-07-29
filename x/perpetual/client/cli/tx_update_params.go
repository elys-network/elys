package cli

import (
	"errors"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/spf13/cobra"
)

const (
	FlagLeverageMax                         = "leverage-max"
	FlagBorrowInterestRateMax               = "borrow-interest-rate-max"
	FlagBorrowInterestRateMin               = "borrow-interest-rate-min"
	FlagBorrowInterestRateIncrease          = "borrow-interest-rate-increase"
	FlagBorrowInterestRateDecrease          = "borrow-interest-rate-decrease"
	FlagHealthGainFactor                    = "health-gain-factor"
	FlagMaxOpenPositions                    = "max-open-positions"
	FlagPoolOpenThreshold                   = "pool-open-threshold"
	FlagBorrowInterestPaymentEnabled        = "borrow-interest-payment-enabled"
	FlagBorrowInterestPaymentFundPercentage = "borrow-interest-payment-fund-percentage"
	FlagBorrowInterestPaymentFundAddress    = "borrow-interest-payment-fund-address"
	FlagSafetyFactor                        = "safety-factor"
	FlagWhitelistingEnabled                 = "whitelisting-enabled"
	FlagEpochLength                         = "epoch-length"
)

// Governance command
func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params",
		Short: "Update perpetual params",
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

			leverageMax, err := cmd.Flags().GetString(FlagLeverageMax)
			if err != nil {
				return err
			}

			borrowInterestRateMax, err := cmd.Flags().GetString(FlagBorrowInterestRateMax)
			if err != nil {
				return err
			}

			borrowInterestRateMin, err := cmd.Flags().GetString(FlagBorrowInterestRateMin)
			if err != nil {
				return err
			}

			borrowInterestRateIncrease, err := cmd.Flags().GetString(FlagBorrowInterestRateIncrease)
			if err != nil {
				return err
			}

			borrowInterestRateDecrease, err := cmd.Flags().GetString(FlagBorrowInterestRateDecrease)
			if err != nil {
				return err
			}

			expedited, err := cmd.Flags().GetBool(FlagExpedited)
			if err != nil {
				return err
			}
			healthGainFactor, err := cmd.Flags().GetString(FlagHealthGainFactor)
			if err != nil {
				return err
			}

			maxOpenPositions, err := cmd.Flags().GetInt64(FlagMaxOpenPositions)
			if err != nil {
				return err
			}

			poolMaxLiabilitiesThreshold, err := cmd.Flags().GetString(FlagPoolOpenThreshold)
			if err != nil {
				return err
			}

			borrowInterestPaymentFundPercentage, err := cmd.Flags().GetString(FlagBorrowInterestPaymentFundPercentage)
			if err != nil {
				return err
			}

			safetyFactor, err := cmd.Flags().GetString(FlagSafetyFactor)
			if err != nil {
				return err
			}

			borrowInterestPaymentEnabled, err := cmd.Flags().GetBool(FlagBorrowInterestPaymentEnabled)
			if err != nil {
				return err
			}

			whitelistingEnabled, err := cmd.Flags().GetBool(FlagWhitelistingEnabled)
			if err != nil {
				return err
			}

			params := &types.Params{
				LeverageMax:                         sdkmath.LegacyMustNewDecFromStr(leverageMax),
				BorrowInterestRateMax:               sdkmath.LegacyMustNewDecFromStr(borrowInterestRateMax),
				BorrowInterestRateMin:               sdkmath.LegacyMustNewDecFromStr(borrowInterestRateMin),
				BorrowInterestRateIncrease:          sdkmath.LegacyMustNewDecFromStr(borrowInterestRateIncrease),
				BorrowInterestRateDecrease:          sdkmath.LegacyMustNewDecFromStr(borrowInterestRateDecrease),
				HealthGainFactor:                    sdkmath.LegacyMustNewDecFromStr(healthGainFactor),
				MaxOpenPositions:                    maxOpenPositions,
				PoolMaxLiabilitiesThreshold:         sdkmath.LegacyMustNewDecFromStr(poolMaxLiabilitiesThreshold),
				BorrowInterestPaymentFundPercentage: sdkmath.LegacyMustNewDecFromStr(borrowInterestPaymentFundPercentage),
				SafetyFactor:                        sdkmath.LegacyMustNewDecFromStr(safetyFactor),
				BorrowInterestPaymentEnabled:        borrowInterestPaymentEnabled,
				WhitelistingEnabled:                 whitelistingEnabled,
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

			if err = msg.ValidateBasic(); err != nil {
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

			govMsg, err := v1.NewMsgSubmitProposal([]sdk.Msg{msg}, deposit, signer.String(), metadata, title, summary, expedited)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), govMsg)
		},
	}

	cmd.Flags().String(FlagLeverageMax, "", "max leverage (integer)")
	cmd.Flags().String(FlagBorrowInterestRateMax, "", "max borrow interest rate (decimal)")
	cmd.Flags().String(FlagBorrowInterestRateMin, "", "min borrow interest rate (decimal)")
	cmd.Flags().String(FlagBorrowInterestRateIncrease, "", "borrow interest rate increase (decimal)")
	cmd.Flags().String(FlagBorrowInterestRateDecrease, "", "borrow interest rate decrease (decimal)")
	cmd.Flags().String(FlagHealthGainFactor, "", "health gain factor (decimal)")
	cmd.Flags().Int64(FlagEpochLength, 1, "epoch length in blocks (integer)")
	cmd.Flags().Int64(FlagMaxOpenPositions, 10000, "max open positions")
	cmd.Flags().String(FlagPoolOpenThreshold, "", "threshold to prevent new positions (decimal range 0-1)")
	cmd.Flags().Bool(FlagBorrowInterestPaymentEnabled, true, "enable incremental borrow interest payment")
	cmd.Flags().String(FlagBorrowInterestPaymentFundPercentage, "", "percentage of incremental borrow interest payment proceeds for fund (decimal range 0-1)")
	cmd.Flags().String(FlagBorrowInterestPaymentFundAddress, "", "address of fund wallet for incremental borrow interest payment")
	cmd.Flags().String(FlagSafetyFactor, "", "the safety factor used in liquidation ratio")
	cmd.Flags().Bool(FlagWhitelistingEnabled, false, "Enable whitelisting")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().Bool(FlagExpedited, false, "expedited")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(FlagLeverageMax)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateMax)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateMin)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateIncrease)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateDecrease)
	_ = cmd.MarkFlagRequired(FlagHealthGainFactor)
	_ = cmd.MarkFlagRequired(FlagMaxOpenPositions)
	_ = cmd.MarkFlagRequired(FlagPoolOpenThreshold)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestPaymentEnabled)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestPaymentFundPercentage)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestPaymentFundAddress)
	_ = cmd.MarkFlagRequired(FlagSafetyFactor)
	_ = cmd.MarkFlagRequired(FlagWhitelistingEnabled)
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
