package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/spf13/cobra"
)

func CmdListPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-price",
		Short: "list all price",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPriceRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PriceAll(context.Background(), params)
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

func CmdShowPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-price [asset] [flags]",
		Short: "shows a price",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			// retrieve the amount of coins allowed to be paid for oracle request fee from the pool account.
			source, err := cmd.Flags().GetString(flagSource)
			if err != nil {
				return err
			}

			// retrieve the amount of gas allowed for the prepare step of the oracle script.
			timestamp, err := cmd.Flags().GetUint64(flagTimestamp)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryGetPriceRequest{
				Asset:     args[0],
				Source:    source,
				Timestamp: timestamp,
			}

			res, err := queryClient.Price(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(flagSource, "", "the source of token price")
	cmd.Flags().Uint64(flagTimestamp, 0, "timestamp when token price is set (optional).")

	return cmd
}
