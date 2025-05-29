package cli

import (
	"errors"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
	"github.com/spf13/cobra"
)

const (
	FlagMarketOrderEnabled   = "market-order-enabled"
	FlagStakeEnabled         = "stake-enabled"
	FlagProcessOrdersEnabled = "process-orders-enabled"
	FlagSwapEnabled          = "swap-enabled"
	FlagPerpetualEnabled     = "perpetual-enabled"
	FlagRewardEnabled        = "reward-enabled"
	FlagLeverageEnabled      = "leverage-enabled"
	FlagLimitProcessOrder    = "limit-process-order"
	FlagRewardPercentage     = "reward-percentage"
	FlagMarginError          = "margin-error"
	FlagMinimumDeposit       = "minimum-deposit"
	FlagExpedited            = "expedited"
)

// Governance command
func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params",
		Short: "Update tradeshield params",
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

			marketOrderEnabled, err := cmd.Flags().GetBool(FlagMarketOrderEnabled)
			if err != nil {
				return err
			}

			stakeEnabled, err := cmd.Flags().GetBool(FlagStakeEnabled)
			if err != nil {
				return err
			}

			processOrdersEnabled, err := cmd.Flags().GetBool(FlagProcessOrdersEnabled)
			if err != nil {
				return err
			}

			swapEnabled, err := cmd.Flags().GetBool(FlagSwapEnabled)
			if err != nil {
				return err
			}

			perpetualEnabled, err := cmd.Flags().GetBool(FlagPerpetualEnabled)
			if err != nil {
				return err
			}

			rewardEnabled, err := cmd.Flags().GetBool(FlagRewardEnabled)
			if err != nil {
				return err
			}

			leverageEnabled, err := cmd.Flags().GetBool(FlagLeverageEnabled)
			if err != nil {
				return err
			}

			limitProcessOrder, err := cmd.Flags().GetUint64(FlagLimitProcessOrder)
			if err != nil {
				return err
			}

			rewardPercentage, err := cmd.Flags().GetString(FlagRewardPercentage)
			if err != nil {
				return err
			}

			marginError, err := cmd.Flags().GetString(FlagMarginError)
			if err != nil {
				return err
			}

			minDeposit, err := cmd.Flags().GetString(FlagMinimumDeposit)
			if err != nil {
				return err
			}

			expedited, err := cmd.Flags().GetBool(FlagExpedited)
			if err != nil {
				return err
			}

			minimumDeposit, ok := math.NewIntFromString(minDeposit)
			if !ok {
				return errors.New("invalid minimum deposit amount")
			}

			params := types.Params{
				MarketOrderEnabled:   marketOrderEnabled,
				StakeEnabled:         stakeEnabled,
				ProcessOrdersEnabled: processOrdersEnabled,
				SwapEnabled:          swapEnabled,
				PerpetualEnabled:     perpetualEnabled,
				RewardEnabled:        rewardEnabled,
				LeverageEnabled:      leverageEnabled,
				LimitProcessOrder:    limitProcessOrder,
				RewardPercentage:     math.LegacyMustNewDecFromStr(rewardPercentage),
				MarginError:          math.LegacyMustNewDecFromStr(marginError),
				MinimumDeposit:       minimumDeposit,
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

	cmd.Flags().Bool(FlagMarketOrderEnabled, false, "market order enabled")
	cmd.Flags().Bool(FlagStakeEnabled, false, "stake enabled")
	cmd.Flags().Bool(FlagProcessOrdersEnabled, false, "process order enabled")
	cmd.Flags().Bool(FlagSwapEnabled, false, "swap enabled")
	cmd.Flags().Bool(FlagPerpetualEnabled, false, "perpetual enabled")
	cmd.Flags().Bool(FlagRewardEnabled, false, "reward enabled")
	cmd.Flags().Bool(FlagLeverageEnabled, false, "leverage enabled")
	cmd.Flags().Uint64(FlagLimitProcessOrder, 10000000, "max limit order processed")
	cmd.Flags().String(FlagRewardPercentage, "", "percentage of rewards given to watchers (decimal range 0-1)")
	cmd.Flags().String(FlagMarginError, "", "percentage of margin error on orders triggered by price (decimal range 0-1)")
	cmd.Flags().String(FlagMinimumDeposit, "", "minimum deposit amount for watchers")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	cmd.Flags().Bool(FlagExpedited, false, "expedited")
	_ = cmd.MarkFlagRequired(FlagMarketOrderEnabled)
	_ = cmd.MarkFlagRequired(FlagStakeEnabled)
	_ = cmd.MarkFlagRequired(FlagProcessOrdersEnabled)
	_ = cmd.MarkFlagRequired(FlagSwapEnabled)
	_ = cmd.MarkFlagRequired(FlagPerpetualEnabled)
	_ = cmd.MarkFlagRequired(FlagRewardEnabled)
	_ = cmd.MarkFlagRequired(FlagLeverageEnabled)
	_ = cmd.MarkFlagRequired(FlagLimitProcessOrder)
	_ = cmd.MarkFlagRequired(FlagRewardPercentage)
	_ = cmd.MarkFlagRequired(FlagMarginError)
	_ = cmd.MarkFlagRequired(FlagMinimumDeposit)
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
