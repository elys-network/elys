package cli

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/elys-network/elys/v7/x/vaults/types"
)

func CmdPerformActionJoinPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join-pool [vault-id] [pool-id] [share-amount-out] [max-amounts-in]",
		Short: "Join a pool from a vault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get vault ID
			vaultId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid vault ID: %w", err)
			}

			// Get pool ID
			poolId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool ID: %w", err)
			}

			// Get share amount out
			shareAmountOut, ok := sdkmath.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid share amount out: %s", args[2])
			}

			// Get max amounts in
			maxAmountsIn, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return fmt.Errorf("invalid max amounts in: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgPerformActionJoinPool{
				Creator:        clientCtx.GetFromAddress().String(),
				VaultId:        vaultId,
				PoolId:         poolId,
				ShareAmountOut: shareAmountOut,
				MaxAmountsIn:   maxAmountsIn,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdPerformActionExitPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exit-pool [vault-id] [pool-id] [share-amount-in] [min-amounts-out] [token-out-denom]",
		Short: "Exit a pool from a vault",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get vault ID
			vaultId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid vault ID: %w", err)
			}

			// Get pool ID
			poolId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool ID: %w", err)
			}

			// Get share amount in
			shareAmountIn, ok := sdkmath.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid share amount in: %s", args[2])
			}

			// Get min amounts out
			minAmountsOut, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return fmt.Errorf("invalid min amounts out: %w", err)
			}

			// Get token out denom
			tokenOutDenom := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgPerformActionExitPool{
				Creator:       clientCtx.GetFromAddress().String(),
				VaultId:       vaultId,
				PoolId:        poolId,
				ShareAmountIn: shareAmountIn,
				MinAmountsOut: minAmountsOut,
				TokenOutDenom: tokenOutDenom,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdPerformActionSwapByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-by-denom [vault-id] [amount] [min-amount] [max-amount] [denom-in] [denom-out]",
		Short: "Swap tokens by denomination from a vault",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get vault ID
			vaultId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid vault ID: %w", err)
			}

			// Get amount
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			// Get min amount
			minAmount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return fmt.Errorf("invalid min amount: %w", err)
			}

			// Get max amount
			maxAmount, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return fmt.Errorf("invalid max amount: %w", err)
			}

			// Get denom in
			denomIn := args[4]

			// Get denom out
			denomOut := args[5]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgPerformActionSwapByDenom{
				Creator:   clientCtx.GetFromAddress().String(),
				VaultId:   vaultId,
				Amount:    amount,
				MinAmount: minAmount,
				MaxAmount: maxAmount,
				DenomIn:   denomIn,
				DenomOut:  denomOut,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
