package cli

import (
	"errors"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"github.com/spf13/cobra"
)

func CmdClaimRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claim-rewards [position-ids] [flags]",
		Short:   "Claim rewards from leveragelp position",
		Example: `elysd tx leveragelp claim-rewards 1,2 --from=bob --yes --gas=1000000`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			positionStrs := strings.Split(args[0], ",")
			positionIds := []uint64{}
			for _, positionStr := range positionStrs {
				id, err := strconv.Atoi(positionStr)
				if err != nil {
					return err
				}
				positionIds = append(positionIds, uint64(id))
			}

			msg := &types.MsgClaimRewards{
				Sender: signer.String(),
				Ids:    positionIds,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
