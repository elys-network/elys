package cli

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdOpenEstimation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-estimation [amm_pool_id] [collateral] [leverage]",
		Short: "Query open estimation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			collateral, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			poolId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			leverage, err := math.LegacyNewDecFromStr(args[2])
			if err != nil {
				return err
			}

			params := &types.QueryOpenEstRequest{
				CollateralAsset:  collateral.Denom,
				CollateralAmount: collateral.Amount,
				AmmPoolId:        uint64(poolId),
				Leverage:         leverage,
			}

			res, err := queryClient.OpenEst(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
