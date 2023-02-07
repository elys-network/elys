package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/x/assetprofile/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

func CmdCreateEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-entry [base-denom] [decimals] [denom] [path] [ibc-channel-id] [ibc-counterparty-channel-id] [display-name] [display-symbol] [network] [address] [external-symbol] [transfer-limit] [permissions] [unit-denom] [ibc-counterparty-denom] [ibc-counterparty-chain-id]",
		Short: "Create a new entry",
		Args:  cobra.ExactArgs(16),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexBaseDenom := args[0]

			// Get value arguments
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
			argPermissions := strings.Split(args[12], listSeparator)
			argUnitDenom := args[13]
			argIbcCounterpartyDenom := args[14]
			argIbcCounterpartyChainId := args[15]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateEntry(
				clientCtx.GetFromAddress().String(),
				indexBaseDenom,
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

func CmdUpdateEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-entry [base-denom] [decimals] [denom] [path] [ibc-channel-id] [ibc-counterparty-channel-id] [display-name] [display-symbol] [network] [address] [external-symbol] [transfer-limit] [permissions] [unit-denom] [ibc-counterparty-denom] [ibc-counterparty-chain-id]",
		Short: "Update a entry",
		Args:  cobra.ExactArgs(16),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexBaseDenom := args[0]

			// Get value arguments
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
			argPermissions := strings.Split(args[12], listSeparator)
			argUnitDenom := args[13]
			argIbcCounterpartyDenom := args[14]
			argIbcCounterpartyChainId := args[15]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateEntry(
				clientCtx.GetFromAddress().String(),
				indexBaseDenom,
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

func CmdDeleteEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-entry [base-denom]",
		Short: "Delete a entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexBaseDenom := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteEntry(
				clientCtx.GetFromAddress().String(),
				indexBaseDenom,
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
