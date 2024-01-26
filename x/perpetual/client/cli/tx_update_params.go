package cli

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/spf13/cobra"
)

const (
	FlagLeverageMax                                    = "leverage-max"
	FlagBorrowInterestRateMax                          = "borrow-interest-rate-max"
	FlagBorrowInterestRateMin                          = "borrow-interest-rate-min"
	FlagBorrowInterestRateIncrease                     = "borrow-interest-rate-increase"
	FlagBorrowInterestRateDecrease                     = "borrow-interest-rate-decrease"
	FlagHealthGainFactor                               = "health-gain-factor"
	FlagMaxOpenPositions                               = "max-open-positions"
	FlagPoolOpenThreshold                              = "pool-open-threshold"
	FlagForceCloseFundPercentage                       = "force-close-fund-percentage"
	FlagForceCloseFundAddress                          = "force-close-fund-address"
	FlagIncrementalBorrowInterestPaymentEnabled        = "incremental-borrow-interest-payment-enabled"
	FlagIncrementalBorrowInterestPaymentFundPercentage = "incremental-borrow-interest-payment-fund-percentage"
	FlagIncrementalBorrowInterestPaymentFundAddress    = "incremental-borrow-interest-payment-fund-address"
	FlagSafetyFactor                                   = "safety-factor"
	FlagWhitelistingEnabled                            = "whitelisting-enabled"
	FlagEpochLength                                    = "epoch-length"
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

			healthGainFactor, err := cmd.Flags().GetString(FlagHealthGainFactor)
			if err != nil {
				return err
			}

			epochLength, err := cmd.Flags().GetInt64(FlagEpochLength)
			if err != nil {
				return err
			}

			maxOpenPositions, err := cmd.Flags().GetInt64(FlagMaxOpenPositions)
			if err != nil {
				return err
			}

			poolOpenThreshold, err := cmd.Flags().GetString(FlagPoolOpenThreshold)
			if err != nil {
				return err
			}

			forceCloseFundPercentage, err := cmd.Flags().GetString(FlagForceCloseFundPercentage)
			if err != nil {
				return err
			}

			forceCloseFundAddress, err := cmd.Flags().GetString(FlagForceCloseFundAddress)
			if err != nil {
				return err
			}

			incrementalBorrowInterestPaymentFundPercentage, err := cmd.Flags().GetString(FlagIncrementalBorrowInterestPaymentFundPercentage)
			if err != nil {
				return err
			}

			incrementalBorrowInterestPaymentFundAddress, err := cmd.Flags().GetString(FlagIncrementalBorrowInterestPaymentFundAddress)
			if err != nil {
				return err
			}

			safetyFactor, err := cmd.Flags().GetString(FlagSafetyFactor)
			if err != nil {
				return err
			}

			incrementalBorrowInterestPaymentEnabled, err := cmd.Flags().GetBool(FlagIncrementalBorrowInterestPaymentEnabled)
			if err != nil {
				return err
			}

			whitelistingEnabled, err := cmd.Flags().GetBool(FlagWhitelistingEnabled)
			if err != nil {
				return err
			}

			params := &types.Params{
				LeverageMax:                sdk.MustNewDecFromStr(leverageMax),
				BorrowInterestRateMax:      sdk.MustNewDecFromStr(borrowInterestRateMax),
				BorrowInterestRateMin:      sdk.MustNewDecFromStr(borrowInterestRateMin),
				BorrowInterestRateIncrease: sdk.MustNewDecFromStr(borrowInterestRateIncrease),
				BorrowInterestRateDecrease: sdk.MustNewDecFromStr(borrowInterestRateDecrease),
				HealthGainFactor:           sdk.MustNewDecFromStr(healthGainFactor),
				EpochLength:                epochLength,
				MaxOpenPositions:           maxOpenPositions,
				PoolOpenThreshold:          sdk.MustNewDecFromStr(poolOpenThreshold),
				ForceCloseFundPercentage:   sdk.MustNewDecFromStr(forceCloseFundPercentage),
				ForceCloseFundAddress:      forceCloseFundAddress,
				IncrementalBorrowInterestPaymentFundPercentage: sdk.MustNewDecFromStr(incrementalBorrowInterestPaymentFundPercentage),
				IncrementalBorrowInterestPaymentFundAddress:    incrementalBorrowInterestPaymentFundAddress,
				SafetyFactor:                            sdk.MustNewDecFromStr(safetyFactor),
				IncrementalBorrowInterestPaymentEnabled: incrementalBorrowInterestPaymentEnabled,
				WhitelistingEnabled:                     whitelistingEnabled,
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

	cmd.Flags().String(FlagLeverageMax, "", "max leverage (integer)")
	cmd.Flags().String(FlagBorrowInterestRateMax, "", "max borrow interest rate (decimal)")
	cmd.Flags().String(FlagBorrowInterestRateMin, "", "min borrow interest rate (decimal)")
	cmd.Flags().String(FlagBorrowInterestRateIncrease, "", "borrow interest rate increase (decimal)")
	cmd.Flags().String(FlagBorrowInterestRateDecrease, "", "borrow interest rate decrease (decimal)")
	cmd.Flags().String(FlagHealthGainFactor, "", "health gain factor (decimal)")
	cmd.Flags().Int64(FlagEpochLength, 1, "epoch length in blocks (integer)")
	cmd.Flags().Int64(FlagMaxOpenPositions, 10000, "max open positions")
	cmd.Flags().String(FlagPoolOpenThreshold, "", "threshold to prevent new positions (decimal range 0-1)")
	cmd.Flags().String(FlagForceCloseFundPercentage, "", "percentage of force close proceeds for fund (decimal range 0-1)")
	cmd.Flags().String(FlagForceCloseFundAddress, "", "address of fund wallet for force close")
	cmd.Flags().Bool(FlagIncrementalBorrowInterestPaymentEnabled, true, "enable incremental borrow interest payment")
	cmd.Flags().String(FlagIncrementalBorrowInterestPaymentFundPercentage, "", "percentage of incremental borrow interest payment proceeds for fund (decimal range 0-1)")
	cmd.Flags().String(FlagIncrementalBorrowInterestPaymentFundAddress, "", "address of fund wallet for incremental borrow interest payment")
	cmd.Flags().String(FlagSafetyFactor, "", "the safety factor used in liquidation ratio")
	cmd.Flags().Bool(FlagWhitelistingEnabled, false, "Enable whitelisting")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(FlagLeverageMax)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateMax)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateMin)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateIncrease)
	_ = cmd.MarkFlagRequired(FlagBorrowInterestRateDecrease)
	_ = cmd.MarkFlagRequired(FlagHealthGainFactor)
	_ = cmd.MarkFlagRequired(FlagMaxOpenPositions)
	_ = cmd.MarkFlagRequired(FlagPoolOpenThreshold)
	_ = cmd.MarkFlagRequired(FlagForceCloseFundPercentage)
	_ = cmd.MarkFlagRequired(FlagForceCloseFundAddress)
	_ = cmd.MarkFlagRequired(FlagIncrementalBorrowInterestPaymentEnabled)
	_ = cmd.MarkFlagRequired(FlagIncrementalBorrowInterestPaymentFundPercentage)
	_ = cmd.MarkFlagRequired(FlagIncrementalBorrowInterestPaymentFundAddress)
	_ = cmd.MarkFlagRequired(FlagSafetyFactor)
	_ = cmd.MarkFlagRequired(FlagWhitelistingEnabled)
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
