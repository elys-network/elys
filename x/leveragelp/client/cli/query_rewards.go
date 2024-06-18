package cli

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [address] [position-ids]",
		Short: "Query position",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]
			_, err = sdk.AccAddressFromBech32(reqAddress)
			if err != nil {
				return err
			}

			positionStrs := strings.Split(args[1], ",")
			positionIds := []uint64{}
			for _, positionStr := range positionStrs {
				id, err := strconv.Atoi(positionStr)
				if err != nil {
					return err
				}
				positionIds = append(positionIds, uint64(id))
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRewardsRequest{
				Address: reqAddress,
				Ids:     positionIds,
			}

			res, err := queryClient.Rewards(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
