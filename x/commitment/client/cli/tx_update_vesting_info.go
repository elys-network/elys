package cli

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

			argBaseDenom, err := cmd.Flags().GetString(BaseDenom)
			if err != nil {
				return err
			}

			argVestingDenom, err := cmd.Flags().GetString(VestingDenom)
			if err != nil {
				return err
			}

			argEpochIdentifier, err := cmd.Flags().GetString(EpochIdentifier)
			if err != nil {
				return err
			}

			argNumEpochs, err := cmd.Flags().GetString(NumEpochs)
			if err != nil {
				return err
			}

			argVestNowFactor, err := cmd.Flags().GetString(VestNowFactor)
			if err != nil {
				return err
			}

			argNumMaxVestings, err := cmd.Flags().GetString(NumMaxVestings)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			numEpochs, err := strconv.ParseInt(argNumEpochs, 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(govtypes.ErrInvalidProposalMsg, "invalid proposal; %d", argNumEpochs)
			}

			vestNowFactor, err := strconv.ParseInt(argVestNowFactor, 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(govtypes.ErrInvalidProposalMsg, "invalid proposal; %d", argVestNowFactor)
			}

			maxVestings, err := strconv.ParseInt(argNumMaxVestings, 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(govtypes.ErrInvalidProposalMsg, "invalid proposal; %d", argNumMaxVestings)
			}

			govAddress := sdk.AccAddress(address.Module(govtypes.ModuleName))
			msg := types.NewMsgUpdateVestingInfo(
				govAddress.String(),
				argBaseDenom,
				argVestingDenom,
				argEpochIdentifier,
				numEpochs,
				vestNowFactor,
				maxVestings,
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

	cmd.Flags().String(BaseDenom, "", "base denom")
	cmd.Flags().String(VestingDenom, "", "vesting-denom")
	cmd.Flags().String(EpochIdentifier, "", "epoch-identifier")
	cmd.Flags().String(NumEpochs, "", "num-epochs")
	cmd.Flags().String(VestNowFactor, "", "vest-now-factor")
	cmd.Flags().String(NumMaxVestings, "", "num-max-vestings")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagSummary, "", "summary of proposal")
	cmd.Flags().String(cli.FlagMetadata, "", "metadata of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(BaseDenom)
	_ = cmd.MarkFlagRequired(VestingDenom)
	_ = cmd.MarkFlagRequired(EpochIdentifier)
	_ = cmd.MarkFlagRequired(NumEpochs)
	_ = cmd.MarkFlagRequired(VestNowFactor)
	_ = cmd.MarkFlagRequired(NumMaxVestings)
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagSummary)
	_ = cmd.MarkFlagRequired(cli.FlagMetadata)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
