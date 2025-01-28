package keepers

import (
	"context"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/cmd/cometbft/commands/debug"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cmtjsonclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type ICSStakingKeeper struct {
	stakingkeeper.Keeper
}

// NewICSStakingKeeper creates a new ICS Staking Keeper instance
func NewICSStakingKeeper(
	keeper stakingkeeper.Keeper,
) ICSStakingKeeper {
	_ = debug.DebugCmd

	return ICSStakingKeeper{
		Keeper: keeper,
	}
}

// GetPubKeyByConsAddr returns the consensus public key by consensus address.
func (k ICSStakingKeeper) GetPubKeyByConsAddr(ctx context.Context, addr sdk.ConsAddress) (cmtprotocrypto.PublicKey, error) {
	var height *int64
	page := 1
	limit := 100

	cmtRPCEndpoint, err := debug.DebugCmd.PersistentFlags().GetString("rpc-laddr")
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	cmtHTTPClient, err := cmtjsonclient.DefaultHTTPClient(cmtRPCEndpoint)
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	cmtRPCClient, err := rpchttp.NewWithClient(cmtRPCEndpoint, "/websocket", cmtHTTPClient)
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	clientCtx := client.Context{
		Client: cmtRPCClient,
	}

	response, err := cmtservice.ValidatorsOutput(ctx, clientCtx, height, page, limit)
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	vals := response.Validators
	for _, v := range vals {
		if v.Address == addr.String() {
			pubkey := v.PubKey.GetCachedValue().(cryptotypes.PubKey)
			tmPk, err := cryptocodec.ToCmtProtoPublicKey(pubkey)
			if err != nil {
				return cmtprotocrypto.PublicKey{}, err
			}
			return tmPk, nil
		}
	}

	return cmtprotocrypto.PublicKey{}, stakingtypes.ErrNoValidatorFound
}

func (k ICSStakingKeeper) GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error) {
	var height *int64
	page := 1
	limit := 100

	cmtRPCEndpoint, err := debug.DebugCmd.PersistentFlags().GetString("rpc-laddr")
	if err != nil {
		return nil, err
	}

	cmtHTTPClient, err := cmtjsonclient.DefaultHTTPClient(cmtRPCEndpoint)
	if err != nil {
		return nil, err
	}

	cmtRPCClient, err := rpchttp.NewWithClient(cmtRPCEndpoint, "/websocket", cmtHTTPClient)
	if err != nil {
		return nil, err
	}

	clientCtx := client.Context{
		Client: cmtRPCClient,
	}

	response, err := cmtservice.ValidatorsOutput(ctx, clientCtx, height, page, limit)
	if err != nil {
		return nil, err
	}

	vals := response.Validators
	stakingVals := make([]stakingtypes.Validator, 0)
	powerReduction := k.PowerReduction(ctx)

	for _, v := range vals {
		if v.VotingPower > 0 {
			stakingVals = append(stakingVals, stakingtypes.Validator{
				OperatorAddress: v.Address,
				Tokens:          sdkmath.NewInt(v.VotingPower).Mul(powerReduction),
				Status:          stakingtypes.Bonded,
			})
		}

	}

	return stakingVals, nil
}
