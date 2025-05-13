package keepers

import (
	"bytes"
	"context"
	"errors"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/cmd/cometbft/commands/debug"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ccvconsumerkeeper "github.com/cosmos/interchain-security/v6/x/ccv/consumer/keeper"
)

type ICSValidatorKeeper struct {
	consumerKeeper ccvconsumerkeeper.Keeper
}

// NewICSValidatorKeeper creates a new ICS Staking Keeper instance
func NewICSValidatorKeeper(
	consumerKeeper ccvconsumerkeeper.Keeper,
) ICSValidatorKeeper {
	_ = debug.DebugCmd

	return ICSValidatorKeeper{
		consumerKeeper: consumerKeeper,
	}
}

func (k ICSValidatorKeeper) Validator(_ context.Context, _ sdk.ValAddress) (stakingtypes.ValidatorI, error) {
	return stakingtypes.Validator{}, errors.New("not implemented for an ics chain")
}

// TotalBondedTokens We have k.consumerKeeper.GetAllCCValidator but not sure if it contains unbonding and unbonded validators.
// So we return error for now.
func (k ICSValidatorKeeper) TotalBondedTokens(_ context.Context) (sdkmath.Int, error) {
	return sdkmath.Int{}, errors.New("not implemented for an ics chain")
}

func (k ICSValidatorKeeper) Slash(ctx context.Context, address sdk.ConsAddress, i int64, i2 int64, dec sdkmath.LegacyDec) (sdkmath.Int, error) {
	return k.consumerKeeper.Slash(ctx, address, i, i2, dec)
}

// Jail internally its nil
func (k ICSValidatorKeeper) Jail(ctx context.Context, address sdk.ConsAddress) error {
	return k.consumerKeeper.Jail(ctx, address)
}

func (k ICSValidatorKeeper) ValidatorsPowerStoreIterator(_ context.Context) (storetypes.Iterator, error) {
	return nil, errors.New("not implemented for an ics chain")
}

func (k ICSValidatorKeeper) MaxValidators(_ context.Context) (uint32, error) {
	return 0, errors.New("not implemented for an ics chain")
}

// PowerReduction This value has to be the same as in the provider-chain. ccv consumer directly uses sdk.DefaultPowerReduction
func (k ICSValidatorKeeper) PowerReduction(_ context.Context) (res sdkmath.Int) {
	return sdk.DefaultPowerReduction
}

// GetPubKeyByConsAddr returns the consensus public key by consensus address.
func (k ICSValidatorKeeper) GetPubKeyByConsAddr(ctx context.Context, addr sdk.ConsAddress) (cmtprotocrypto.PublicKey, error) {
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

func (k ICSValidatorKeeper) GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	vals := k.consumerKeeper.GetAllCCValidator(sdkCtx)

	stakingVals := make([]stakingtypes.Validator, 0)
	powerReduction := k.PowerReduction(ctx)

	for _, v := range vals {
		if v.Power > 0 {
			stakingVals = append(stakingVals, stakingtypes.Validator{
				OperatorAddress: sdk.ConsAddress(v.Address).String(), // This is an incorrect operator address,
				// on a consumer chain there is no way to figure that out atm.
				// We only have a consensus address.
				ConsensusPubkey: v.GetPubkey(),
				Tokens:          sdkmath.NewInt(v.Power).Mul(powerReduction),
				Status:          stakingtypes.Bonded,
			})
		}

	}

	return stakingVals, nil
}
