package keepers

import (
	"bytes"
	"context"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/cmd/cometbft/commands/debug"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ccvconsumerkeeper "github.com/cosmos/interchain-security/v6/x/ccv/consumer/keeper"
)

type ICSStakingKeeper struct {
	*stakingkeeper.Keeper
	consumerKeeper ccvconsumerkeeper.Keeper
}

// NewICSStakingKeeper creates a new ICS Staking Keeper instance
func NewICSStakingKeeper(
	stakingKeeper *stakingkeeper.Keeper,
	consumerKeeper ccvconsumerkeeper.Keeper,
) ICSStakingKeeper {
	_ = debug.DebugCmd

	return ICSStakingKeeper{
		Keeper:         stakingKeeper,
		consumerKeeper: consumerKeeper,
	}
}

// GetPubKeyByConsAddr returns the consensus public key by consensus address.
func (k ICSStakingKeeper) GetPubKeyByConsAddr(ctx context.Context, addr sdk.ConsAddress) (cmtprotocrypto.PublicKey, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	vals := k.consumerKeeper.GetAllCCValidator(sdkCtx)

	for _, v := range vals {
		if bytes.Equal(v.Address, addr.Bytes()) {
			pubkey := v.Pubkey.GetCachedValue().(cryptotypes.PubKey)
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
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	vals := k.consumerKeeper.GetAllCCValidator(sdkCtx)

	stakingVals := make([]stakingtypes.Validator, 0)
	powerReduction := k.PowerReduction(ctx)

	for _, v := range vals {
		if v.Power > 0 {
			stakingVals = append(stakingVals, stakingtypes.Validator{
				OperatorAddress: string(v.Address),
				Tokens:          sdkmath.NewInt(v.Power).Mul(powerReduction),
				Status:          stakingtypes.Bonded,
			})
		}

	}

	return stakingVals, nil
}
