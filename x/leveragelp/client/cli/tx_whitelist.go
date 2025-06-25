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
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdWhitelist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist [address]",
		Short: "Whitelist the provided address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.New("invalid whitelisted address")
			}

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

			govAddress := sdk.AccAddress(address.Module("gov"))
			msg := types.NewMsgWhitelist(
				govAddress.String(),
				args[0],
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
