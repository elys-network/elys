package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/spf13/cobra"
)

// CmdBandPriceResult queries request result by reqID
func CmdBandPriceResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "band-price-result [request-id]",
		Short: "Query the BandPrice result data by id",
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
			r, err := queryClient.BandPriceResult(context.Background(), &types.QueryBandPriceRequest{RequestId: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(r)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdLastBandRequestId queries latest request
func CmdLastBandRequestId() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-band-request-id",
		Short: "Query the last request id returned by BandPrice ack packet",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			r, err := queryClient.LastBandRequestId(context.Background(), &types.QueryLastBandRequestIdRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(r)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
