package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/elys-network/elys/x/incentive/types"
)

var (
	FlagCommission                        = "commission"
	FlagValidatorAddress                  = "validator-address"
	FlagMaxMessagesPerTx                  = "max-msgs"
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	MaxMessagesPerTxDefault    = 0
	listSeparator              = ","
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
		NewWithdrawRewardsCmd(),
	)

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

// NewWithdrawRewardsCmd returns a CLI command handler for creating a MsgWithdrawDelegatorReward transaction.
func NewWithdrawRewardsCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "withdraw-rewards",
		Short: "Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw rewards from a given delegation address,
and optionally withdraw validator commission if the delegation address given is a validator operator.

Example:
$ %s tx incentive withdraw-rewards --from mykey
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
			msgs := []sdk.Msg{types.NewMsgWithdrawRewards(delAddr)}

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
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
