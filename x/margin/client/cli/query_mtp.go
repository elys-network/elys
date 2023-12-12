package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/spf13/cobra"
)

func CmdMtp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-mtp [address] [id]",
		Short: "Query mtp",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]
			_, err = sdk.AccAddressFromBech32(reqAddress)
			if err != nil {
				return err
			}

			reqId := args[1]
			id, err := strconv.ParseUint(reqId, 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.MTPRequest{
				Address: reqAddress,
				Id:      id,
			}

			res, err := queryClient.MTP(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
