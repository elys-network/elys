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
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"github.com/spf13/cobra"
)

const (
	FlagLeverageMax         = "leverage-max"
	FlagMaxOpenPositions    = "max-open-positions"
	FlagPoolOpenThreshold   = "pool-open-threshold"
	FlagSafetyFactor        = "safety-factor"
	FlagWhitelistingEnabled = "whitelisting-enabled"
	FlagEpochLength         = "epoch-length"
)

// Governance command
func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params",
		Short: "Update leveragelp params",
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

			epoch_length, err := cmd.Flags().GetInt64(FlagEpochLength)
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

			safetyFactor, err := cmd.Flags().GetString(FlagSafetyFactor)
			if err != nil {
				return err
			}

			whitelistingEnabled, err := cmd.Flags().GetBool(FlagWhitelistingEnabled)
			if err != nil {
				return err
			}

			expedited, err := cmd.Flags().GetBool(FlagExpedited)
			if err != nil {
				return err
			}

			params := &types.Params{
				LeverageMax:         sdkmath.LegacyMustNewDecFromStr(leverageMax),
				EpochLength:         epoch_length,
				MaxOpenPositions:    maxOpenPositions,
				PoolOpenThreshold:   sdkmath.LegacyMustNewDecFromStr(poolOpenThreshold),
				SafetyFactor:        sdkmath.LegacyMustNewDecFromStr(safetyFactor),
				WhitelistingEnabled: whitelistingEnabled,
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

			govMsg, err := v1.NewMsgSubmitProposal([]sdk.Msg{msg}, deposit, signer.String(), metadata, title, summary, expedited)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), govMsg)
		},
	}

	cmd.Flags().String(FlagLeverageMax, "", "max leverage (integer)")
	cmd.Flags().Int64(FlagEpochLength, 1, "epoch length in blocks (integer)")
	cmd.Flags().Int64(FlagMaxOpenPositions, 10000, "max open positions")
	cmd.Flags().String(FlagPoolOpenThreshold, "", "threshold to prevent new positions (decimal range 0-1)")
	cmd.Flags().String(FlagSafetyFactor, "", "the safety factor used in liquidation ratio")
	cmd.Flags().Bool(FlagWhitelistingEnabled, false, "Enable whitelisting")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	cmd.Flags().String(FlagExpedited, "", "expedited")
	_ = cmd.MarkFlagRequired(FlagLeverageMax)
	_ = cmd.MarkFlagRequired(FlagMaxOpenPositions)
	_ = cmd.MarkFlagRequired(FlagPoolOpenThreshold)
	_ = cmd.MarkFlagRequired(FlagSafetyFactor)
	_ = cmd.MarkFlagRequired(FlagWhitelistingEnabled)
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
