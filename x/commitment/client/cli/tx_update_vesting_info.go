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
	"github.com/elys-network/elys/x/commitment/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateVestingInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-vesting-info",
		Short: "Broadcast message update-vesting-info",
		Args:  cobra.ExactArgs(0),
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

			argBaseDenom, err := cmd.Flags().GetString("base-denom")
			if err != nil {
				return err
			}

			argVestingDenom, err := cmd.Flags().GetString("vesting-denom")
			if err != nil {
				return err
			}

			argEpochIdentifier, err := cmd.Flags().GetString("epoch-identifier")
			if err != nil {
				return err
			}

			argNumEpochs, err := cmd.Flags().GetString("num-epochs")
			if err != nil {
				return err
			}

			argVestNowFactor, err := cmd.Flags().GetString("vest-now-factor")
			if err != nil {
				return err
			}

			argNumMaxVestings, err := cmd.Flags().GetString("num-max-vestings")
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			govAddress := sdk.AccAddress(address.Module("gov"))
			msg := types.NewMsgUpdateVestingInfo(
				govAddress.String(),
				argBaseDenom,
				argVestingDenom,
				argEpochIdentifier,
				argNumEpochs,
				argVestNowFactor,
				argNumMaxVestings,
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

	cmd.Flags().String("base-denom", "", "base denom (string)")
	cmd.Flags().String("vesting-denom", "", "vesting-denom (string)")
	cmd.Flags().String("epoch-identifier", "", "epoch-identifier (string)")
	cmd.Flags().String("num-epochs", "", "num-epochs (int64)")
	cmd.Flags().String("vest-now-factor", "", "vest-now-factor (decimal)")
	cmd.Flags().String("num-max-vestings", "", "num-max-vestings (int64)")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired("base-denom")
	_ = cmd.MarkFlagRequired("vesting-denom")
	_ = cmd.MarkFlagRequired("epoch-identifier")
	_ = cmd.MarkFlagRequired("num-epochs")
	_ = cmd.MarkFlagRequired("vest-now-factor")
	_ = cmd.MarkFlagRequired("num-max-vestings")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
