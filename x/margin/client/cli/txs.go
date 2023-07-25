package cli

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ = strconv.Itoa(0)

func CmdOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open margin position",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			collateralAsset, err := cmd.Flags().GetString("collateral_asset")
			if err != nil {
				return err
			}

			collateralAmount, err := cmd.Flags().GetString("collateral_amount")
			if err != nil {
				return err
			}

			borrowAsset, err := cmd.Flags().GetString("borrow_asset")
			if err != nil {
				return err
			}

			position, err := cmd.Flags().GetString("position")
			if err != nil {
				return err
			}
			positionEnum := types.GetPositionFromString(position)

			leverage, err := cmd.Flags().GetString("leverage")
			if err != nil {
				return err
			}

			leverageDec, err := sdk.NewDecFromStr(leverage)
			if err != nil {
				return err
			}

			collateralAmt, ok := sdk.NewIntFromString(collateralAmount)
			if !ok {
				return errors.New("invalid collateral amount")
			}

			msg := types.NewMsgOpen(
				clientCtx.GetFromAddress().String(),
				collateralAsset,
				collateralAmt,
				borrowAsset,
				positionEnum,
				leverageDec,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("collateral_amount", "0", "amount of collateral asset")
	cmd.Flags().String("collateral_asset", "", "symbol of asset")
	cmd.Flags().String("borrow_asset", "", "symbol of asset")
	cmd.Flags().String("position", "", "type of position")
	cmd.Flags().String("leverage", "", "leverage of position")
	_ = cmd.MarkFlagRequired("collateral_amount")
	_ = cmd.MarkFlagRequired("collateral_asset")
	_ = cmd.MarkFlagRequired("borrow_asset")
	_ = cmd.MarkFlagRequired("position")
	_ = cmd.MarkFlagRequired("leverage")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close margin position",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			id, err := cmd.Flags().GetUint64("id")
			if err != nil {
				return err
			}

			msg := types.NewMsgClose(
				signer.String(),
				id,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Uint64("id", 0, "id of the position")
	_ = cmd.MarkFlagRequired("id")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Governance command
func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params",
		Short: "Update margin params",
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

			leverage_max := sdk.MustNewDecFromStr(viper.GetString("leverage-max"))

			if leverage_max.GT(sdk.NewDec(10)) || leverage_max.LT(sdk.NewDec(1)) {
				return errors.New("invalid leverage max, it has to be between 1-10.")
			}

			content := types.NewProposalUpdateParams(
				title,
				description,
				sdk.MustNewDecFromStr(viper.GetString("leverage-max")),
				sdk.MustNewDecFromStr(viper.GetString("interest-rate-max")),
				sdk.MustNewDecFromStr(viper.GetString("interest-rate-min")),
				sdk.MustNewDecFromStr(viper.GetString("interest-rate-increase")),
				sdk.MustNewDecFromStr(viper.GetString("interest-rate-decrease")),
				sdk.MustNewDecFromStr(viper.GetString("health-gain-factor")),
				viper.GetUint64("epoch-length"),
				sdk.MustNewDecFromStr(viper.GetString("removal-queue-threshold")),
				viper.GetUint64("max-open-positions"),
				sdk.MustNewDecFromStr(viper.GetString("pool-open-threshold")),
				sdk.MustNewDecFromStr(viper.GetString("force-close-fund-percentage")),
				viper.GetString("force-close-fund-address"),
				sdk.MustNewDecFromStr(viper.GetString("incremental-interest-payment-fund-percentage")),
				viper.GetString("incremental-interest-payment-fund-address"),
				sdk.MustNewDecFromStr(viper.GetString("sq-modifier")),
				sdk.MustNewDecFromStr(viper.GetString("safety-factor")),
				viper.GetBool("incremental-interest-payment-enabled"),
				viper.GetBool("whitelisting-enabled"),
			)

			from := clientCtx.GetFromAddress()
			if from == nil {
				return errors.New("signer address is missing")
			}

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

	cmd.Flags().String("leverage-max", "", "max leverage (integer)")
	cmd.Flags().String("interest-rate-max", "", "max interest rate (decimal)")
	cmd.Flags().String("interest-rate-min", "", "min interest rate (decimal)")
	cmd.Flags().String("interest-rate-increase", "", "interest rate increase (decimal)")
	cmd.Flags().String("interest-rate-decrease", "", "interest rate decrease (decimal)")
	cmd.Flags().String("health-gain-factor", "", "health gain factor (decimal)")
	cmd.Flags().Int64("epoch-length", 1, "epoch length in blocks (integer)")
	cmd.Flags().Uint64("max-open-positions", 10000, "max open positions")
	cmd.Flags().String("removal-queue-threshold", "", "removal queue threshold (decimal range 0-1)")
	cmd.Flags().String("pool-open-threshold", "", "threshold to prevent new positions (decimal range 0-1)")
	cmd.Flags().String("force-close-fund-percentage", "", "percentage of force close proceeds for fund (decimal range 0-1)")
	cmd.Flags().String("force-close-fund-address", "", "address of fund wallet for force close")
	cmd.Flags().Bool("incremental-interest-payment-enabled", true, "enable incremental interest payment")
	cmd.Flags().String("incremental-interest-payment-fund-percentage", "", "percentage of incremental interest payment proceeds for fund (decimal range 0-1)")
	cmd.Flags().String("incremental-interest-payment-fund-address", "", "address of fund wallet for incremental interest payment")
	cmd.Flags().String("sq-modifier", "", "the modifier value for the removal queue's sq formula")
	cmd.Flags().String("safety-factor", "", "the safety factor used in liquidation ratio")
	cmd.Flags().Bool("whitelisting-enabled", false, "Enable whitelisting")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired("leverage-max")
	_ = cmd.MarkFlagRequired("interest-rate-max")
	_ = cmd.MarkFlagRequired("interest-rate-min")
	_ = cmd.MarkFlagRequired("interest-rate-increase")
	_ = cmd.MarkFlagRequired("interest-rate-decrease")
	_ = cmd.MarkFlagRequired("health-gain-factor")
	_ = cmd.MarkFlagRequired("removal-queue-threshold")
	_ = cmd.MarkFlagRequired("max-open-positions")
	_ = cmd.MarkFlagRequired("pool-open-threshold")
	_ = cmd.MarkFlagRequired("force-close-fund-percentage")
	_ = cmd.MarkFlagRequired("force-close-fund-address")
	_ = cmd.MarkFlagRequired("incremental-interest-payment-enabled")
	_ = cmd.MarkFlagRequired("incremental-interest-payment-fund-percentage")
	_ = cmd.MarkFlagRequired("incremental-interest-payment-fund-address")
	_ = cmd.MarkFlagRequired("sq-modifier")
	_ = cmd.MarkFlagRequired("safety-factor")
	_ = cmd.MarkFlagRequired("whitelisting-enabled")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// Governance command
// TODO
// Need to confirm the format of enabled pools and disabled pools list
// Pool.json
func CmdUpdatePools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pools [pool.json]",
		Short: "Update margin enabled pools, and closed pools",
		Args:  cobra.ExactArgs(1),
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

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			pools, err := readPoolsJSON(args[0])
			if err != nil {
				return err
			}

			closedPools, err := readPoolsJSON(viper.GetString("closed-pools"))
			if err != nil {
				return err
			}

			content := types.NewProposalUpdatePools(
				title,
				description,
				pools,
				closedPools)

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := v1beta1.NewMsgSubmitProposal(content, deposit, signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("closed-pools", "", "pools that new positions cannot be opened on")
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired("closed-pools")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func readPoolsJSON(filename string) ([]string, error) {
	var pools []string
	bz, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}
	err = json.Unmarshal(bz, &pools)
	if err != nil {
		return []string{}, err
	}

	return pools, nil
}

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

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			content := types.NewProposalWhitelist(
				title,
				description,
				args[0],
			)

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := v1beta1.NewMsgSubmitProposal(content, deposit, signer)
			if err != nil {
				return err
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdDewhitelist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dewhitelist [address]",
		Short: "Dewhitelist the provided address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.New("invalid whitelisted address")
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			content := types.NewProposalDewhitelist(
				title,
				description,
				args[0],
			)

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := v1beta1.NewMsgSubmitProposal(content, deposit, signer)
			if err != nil {
				return err
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
