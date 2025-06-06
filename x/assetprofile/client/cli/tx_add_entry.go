package cli

import (
	"strconv"

	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/assetprofile/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAddEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-entry [base-denom] [decimals] [denom] [path] [ibc-channel-id] [ibc-counterparty-channel-id] [display-name] [display-symbol] [network] [address] [external-symbol] [transfer-limit] [permissions] [unit-denom] [ibc-counterparty-denom] [ibc-counterparty-chain-id] [commit-enabled] [withdraw-enabled]",
		Short: "Create a new asset",
		Args:  cobra.ExactArgs(18),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argPermissions := []string{}
			if args[12] != "" {
				argPermissions = strings.Split(args[12], listSeparator)
			}

			argBaseDenom := args[0]
			argDecimals, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argDenom := args[2]
			argPath := args[3]
			argIbcChannelId := args[4]
			argIbcCounterpartyChannelId := args[5]
			argDisplayName := args[6]
			argDisplaySymbol := args[7]
			argNetwork := args[8]
			argAddress := args[9]
			argExternalSymbol := args[10]
			argTransferLimit := args[11]

			argUnitDenom := args[13]
			argIbcCounterpartyDenom := args[14]
			argIbcCounterpartyChainId := args[15]
			argCommitEnabled, err := cast.ToBoolE(args[16])
			if err != nil {
				return err
			}
			argWithdrawEnabled, err := cast.ToBoolE(args[17])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddEntry(
				clientCtx.GetFromAddress().String(),
				argBaseDenom,
				argDecimals,
				argDenom,
				argPath,
				argIbcChannelId,
				argIbcCounterpartyChannelId,
				argDisplayName,
				argDisplaySymbol,
				argNetwork,
				argAddress,
				argExternalSymbol,
				argTransferLimit,
				argPermissions,
				argUnitDenom,
				argIbcCounterpartyDenom,
				argIbcCounterpartyChainId,
				argCommitEnabled,
				argWithdrawEnabled,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
