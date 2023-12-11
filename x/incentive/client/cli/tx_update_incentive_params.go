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
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateIncentiveParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-incentive-params [community-tax] [withdraw-addr-enabled] [reward-portion-for-lps] [reward-portion-for-stakers] [elys-stake-tracking-rate] [max-eden-reward-apr-stakers] [max-eden-reward-apr-lps] [distribution-epoch-for-stakers] [distribution-epoch-for-lps]",
		Short: "Broadcast message update-incentive-params update-incentive-params [community-tax] [withdraw-addr-enabled] [reward-portion-for-lps] [reward-portion-for-stakers] [elys-stake-tracking-rate] [max-eden-reward-apr-stakers] [max-eden-reward-apr-lps] [distribution-epoch-for-stakers] [distribution-epoch-for-lps]",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCommunityTax := args[0]
			argWithdrawAddrEnabled := args[1]
			argRewardPortionForLps := args[2]
			argRewardPortionForStakers := args[3]
			argElysStakeTrackingRate := args[4]
			argMaxEdenRewardAprStakers := args[5]
			argMaxEdenRewardAprLps := args[6]
			argDistributionEpochForStakers := args[7]
			argDistributionEpochForLps := args[8]

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

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			communityTax := sdk.MustNewDecFromStr(argCommunityTax)
			withdarwAddrEnabled, err := strconv.ParseBool(argWithdrawAddrEnabled)
			if err != nil {
				return err
			}
			rewardPortionForLps := sdk.MustNewDecFromStr(argRewardPortionForLps)
			rewardPortionForStakers := sdk.MustNewDecFromStr(argRewardPortionForStakers)
			elysStakeTrackingRate, err := strconv.ParseInt(argElysStakeTrackingRate, 10, 64)
			if err != nil {
				return err
			}
			maxEdenRewardAprStakers := sdk.MustNewDecFromStr(argMaxEdenRewardAprStakers)
			maxEdenRewardLps := sdk.MustNewDecFromStr(argMaxEdenRewardAprLps)
			distributionEpochForStaker, err := strconv.ParseInt(argDistributionEpochForStakers, 10, 64)
			if err != nil {
				return err
			}
			distributionEpochForLps, err := strconv.ParseInt(argDistributionEpochForLps, 10, 64)
			if err != nil {
				return err
			}

			govAddress := sdk.AccAddress(address.Module("gov"))
			msg := types.NewMsgUpdateIncentiveParams(
				govAddress.String(),
				communityTax,
				withdarwAddrEnabled,
				rewardPortionForLps,
				rewardPortionForStakers,
				elysStakeTrackingRate,
				maxEdenRewardAprStakers,
				maxEdenRewardLps,
				distributionEpochForStaker,
				distributionEpochForLps,
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

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), govMsg)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
