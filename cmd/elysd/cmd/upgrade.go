package cmd

import (
	"fmt"
	"strconv"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/spf13/cobra"
)

// SoftwareUpgradeTxCmd implements submitting a proposal transaction command for chain upgrade.
func SoftwareUpgradeTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software-upgrade-tx [name] [height] [deposit] [description] [info] [expedited]",
		Short: "cmd to submit software upgrade proposal",
		Args:  cobra.ExactArgs(6),
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			height, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}
			description := args[3]
			info := args[4]
			expedited, err := strconv.ParseBool(args[5])
			if err != nil {
				return err
			}

			softwareUpgrade := &upgradetypes.MsgSoftwareUpgrade{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Plan: upgradetypes.Plan{
					Name:                name,
					Time:                time.Time{},
					Height:              int64(height),
					Info:                info,
					UpgradedClientState: nil,
				},
			}

			msg, err := v1.NewMsgSubmitProposal([]sdk.Msg{softwareUpgrade}, deposit, clientCtx.GetFromAddress().String(), "", name, description, expedited)
			if err != nil {
				return fmt.Errorf("invalid message: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
