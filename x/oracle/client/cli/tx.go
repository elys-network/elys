package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/oracle/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdRequestCoinRatesData())
	cmd.AddCommand(CmdCreateAssetInfo())
	cmd.AddCommand(CmdUpdateAssetInfo())
	cmd.AddCommand(CmdDeleteAssetInfo())
	cmd.AddCommand(CmdSubmitAddAssetInfoProposal())
	cmd.AddCommand(CmdSubmitRemoveAssetInfoProposal())
	cmd.AddCommand(CmdCreatePrice())
	cmd.AddCommand(CmdUpdatePrice())
	cmd.AddCommand(CmdDeletePrice())
cmd.AddCommand(CmdCreatePriceFeeder())
	cmd.AddCommand(CmdUpdatePriceFeeder())
	cmd.AddCommand(CmdDeletePriceFeeder())
// this line is used by starport scaffolding # 1

	return cmd
}

func CmdSubmitAddAssetInfoProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset-info-proposal [denom] [display] [bandTicker] [binanceTicker] [osmosisTicker]",
		Args:  cobra.ExactArgs(5),
		Short: "Submit an add asset info proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			content := types.NewProposalAddAssetInfo(
				title,
				description,
				args[0],
				args[1],
				args[2],
				args[3],
				args[4],
			)

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := v1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSubmitRemoveAssetInfoProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-asset-info-proposal [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an add asset info proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			content := types.NewProposalRemoveAssetInfo(
				title,
				description,
				args[0],
			)

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := v1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
