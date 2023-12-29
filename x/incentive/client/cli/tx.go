package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
)

var (
	FlagCommission       = "commission"
	FlagValidatorAddress = "validator-address"
	FlagEarnType         = "earn-type"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		CmdWithdrawRewardsCmd(),
	)

	cmd.AddCommand(CmdUpdateIncentiveParams())
	// this line is used by starport scaffolding # 1

	return cmd
}

type newGenerateOrBroadcastFunc func(client.Context, *pflag.FlagSet, ...sdk.Msg) error

func newSplitAndApply(
	genOrBroadcastFn newGenerateOrBroadcastFunc, clientCtx client.Context,
	fs *pflag.FlagSet, msgs []sdk.Msg, chunkSize int,
) error {
	if chunkSize == 0 {
		return genOrBroadcastFn(clientCtx, fs, msgs...)
	}

	// split messages into slices of length chunkSize
	totalMessages := len(msgs)
	for i := 0; i < len(msgs); i += chunkSize {

		sliceEnd := i + chunkSize
		if sliceEnd > totalMessages {
			sliceEnd = totalMessages
		}

		msgChunk := msgs[i:sliceEnd]
		if err := genOrBroadcastFn(clientCtx, fs, msgChunk...); err != nil {
			return err
		}
	}

	return nil
}

// CmdWithdrawRewardsCmd returns a CLI command handler for creating a MsgWithdrawDelegatorReward transaction.
func CmdWithdrawRewardsCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "withdraw-rewards",
		Short: "Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw rewards from a given delegation address,
and optionally withdraw validator commission if the delegation address given is a validator operator.

Example:
$ %s tx incentive withdraw-rewards --from mykey --withdraw-type [0: withdraw all, 1: withdraw usdc program, 2: withdraw elys program, 3: withdraw eden program, 4: withdraw eden boost program.]
$ %s tx incentive withdraw-rewards --from mykey --commission --validator-address %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, bech32PrefixValAddr, bech32PrefixValAddr,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			earnType, err := cmd.Flags().GetInt64(FlagEarnType)
			if err != nil {
				earnType = int64(commitmenttypes.EarnType_ALL_PROGRAM)
			}

			msgs := []sdk.Msg{types.NewMsgWithdrawRewards(delAddr, commitmenttypes.EarnType(earnType))}

			if commission, _ := cmd.Flags().GetBool(FlagCommission); commission {
				if validatorAddr, _ := cmd.Flags().GetString(FlagValidatorAddress); len(validatorAddr) > 0 {
					valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
					if err != nil {
						return err
					}
					msgs = append(msgs, types.NewMsgWithdrawValidatorCommission(delAddr, valAddr))
				}
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	cmd.Flags().Bool(FlagCommission, false, "Withdraw the validator's commission in addition to the rewards")
	cmd.Flags().String(FlagValidatorAddress, "", "Validator's operator address to withdraw commission from")
	cmd.Flags().Int64(FlagEarnType, 0, "Earn type - 0: all earn, 1: usdc program, 2: elys program, 3: eden program, 4: eden boost program.")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
