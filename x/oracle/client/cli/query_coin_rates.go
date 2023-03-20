package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/spf13/cobra"
)

// CmdCoinRatesResult queries request result by reqID
func CmdCoinRatesResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coin-rates-result [request-id]",
		Short: "Query the CoinRates result data by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			r, err := queryClient.CoinRatesResult(context.Background(), &types.QueryCoinRatesRequest{RequestId: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(r)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdLastCoinRatesID queries latest request
func CmdLastCoinRatesID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-coin-rates-id",
		Short: "Query the last request id returned by CoinRates ack packet",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			r, err := queryClient.LastCoinRatesId(context.Background(), &types.QueryLastCoinRatesIdRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(r)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
