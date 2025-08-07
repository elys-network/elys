package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/elys-network/elys/v7/x/masterchef/types"
	"github.com/spf13/cobra"
)

const FlagExpedited = "expedited"

func CmdAddExternalIncentive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-external-incentive [reward-denom] [pool-id] [from-block] [to-block] [amount-per-block]",
		Short: "Broadcast message add-external-incentive [reward-denom] [pool-id] [from-block] [to-block] [amount-per-block]",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			rewardDenom := args[0]
			poolId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return errors.New("invalid pool id")
			}
			fromBlock, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return errors.New("invalid from block")
			}
			toBlock, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return errors.New("invalid to block")
			}
			amountPerBlock, ok := math.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid amount per block")
			}

			msg := &types.MsgAddExternalIncentive{
				Sender:         signer.String(),
				RewardDenom:    rewardDenom,
				PoolId:         poolId,
				FromBlock:      fromBlock,
				ToBlock:        toBlock,
				AmountPerBlock: amountPerBlock,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdAddExternalRewardDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-external-reward-denom [reward-denom] [min-amount] [supported]",
		Short: "Broadcast message add-external-reward-denom [reward-denom] [min-amount] [supported]",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

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

			expedited, err := cmd.Flags().GetBool(FlagExpedited)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			rewardDenom := args[0]
			minAmount, ok := math.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid min amount")
			}
			supported, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}

			govAddress := sdk.AccAddress(address.Module("gov"))
			msg := &types.MsgAddExternalRewardDenom{
				Authority:   govAddress.String(),
				RewardDenom: rewardDenom,
				MinAmount:   minAmount,
				Supported:   supported,
			}
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
	cmd.Flags().String(FlagExpedited, "", "expedited")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdatePoolMultipliers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pool-multipliers [pool-id] [multiplier]",
		Short: "Broadcast message update-pool-multipliers [pool-id] [multiplier]",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

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

			expedited, err := cmd.Flags().GetBool(FlagExpedited)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			poolId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			multiplier := math.LegacyMustNewDecFromStr(args[1])

			poolMultipliers := []types.PoolMultiplier{}
			poolMultipliers = append(poolMultipliers, types.PoolMultiplier{
				PoolId:     poolId,
				Multiplier: multiplier,
			})
			govAddress := sdk.AccAddress(address.Module("gov"))
			msg := &types.MsgUpdatePoolMultipliers{
				Authority:       govAddress.String(),
				PoolMultipliers: poolMultipliers,
			}
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
	cmd.Flags().String(FlagExpedited, "", "expedited")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdClaimRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-rewards",
		Short: "claim rewards including external incentives",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Claim rewards from a given delegation address,
Example:
$ %s tx masterchef claim-rewards --from mykey --pool-ids [pool-ids] 
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			poolIds := []uint64{}

			poolIdsString, err := cmd.Flags().GetString(FlagPoolIds)
			if err == nil && poolIdsString != "" {
				poolIdsArray := strings.Split(poolIdsString, ",")
				for _, poolIdStr := range poolIdsArray {
					poolId, err := strconv.ParseUint(poolIdStr, 10, 64)
					if err != nil {
						return err
					}
					poolIds = append(poolIds, poolId)
				}
			}

			msgs := []sdk.Msg{&types.MsgClaimRewards{
				Sender:  delAddr.String(),
				PoolIds: poolIds,
			}}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	cmd.Flags().String(FlagPoolIds, "", "Validator's operator address to withdraw commission from")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
