package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/spf13/cobra"
)

func CmdListDenomLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-denom-liquidity",
		Short: "list all denom-liquidity",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDenomLiquidityRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DenomLiquidityAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowDenomLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-denom-liquidity [denom]",
		Short: "shows a denom-liquidity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argDenom := args[0]

			params := &types.QueryGetDenomLiquidityRequest{
				Denom: argDenom,
			}

			res, err := queryClient.DenomLiquidity(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
