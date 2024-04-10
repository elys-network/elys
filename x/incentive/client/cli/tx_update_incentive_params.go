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

func CmdUpdateIncentiveParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-incentive-params [reward-portion-for-lps] [reward-portion-for-stakers] [elys-stake-tracking-rate] [max-eden-reward-apr-stakers] [max-eden-reward-apr-lps] [distribution-interval]",
		Short: "Broadcast message update-incentive-params update-incentive-params [reward-portion-for-lps] [reward-portion-for-stakers] [elys-stake-tracking-rate] [max-eden-reward-apr-stakers] [max-eden-reward-apr-lps] [distribution-interval]",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRewardPortionForLps := args[0]
			argRewardPortionForStakers := args[1]
			argElysStakeSnapInterval := args[2]
			argMaxEdenRewardAprStakers := args[3]
			argMaxEdenRewardAprLps := args[4]

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

			rewardPortionForLps := sdk.MustNewDecFromStr(argRewardPortionForLps)
			rewardPortionForStakers := sdk.MustNewDecFromStr(argRewardPortionForStakers)
			elysStakeSnapInterval, err := strconv.ParseInt(argElysStakeSnapInterval, 10, 64)
			if err != nil {
				return err
			}
			maxEdenRewardAprStakers := sdk.MustNewDecFromStr(argMaxEdenRewardAprStakers)
			maxEdenRewardLps := sdk.MustNewDecFromStr(argMaxEdenRewardAprLps)

			govAddress := sdk.AccAddress(address.Module("gov"))
			msg := types.NewMsgUpdateIncentiveParams(
				govAddress.String(),
				rewardPortionForLps,
				rewardPortionForStakers,
				elysStakeSnapInterval,
				maxEdenRewardAprStakers,
				maxEdenRewardLps,
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
