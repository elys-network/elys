package types

import (
	"errors"
	"fmt"
	"github.com/elys-network/elys/v7/utils"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgClose{}
	_ sdk.Msg = &MsgOpen{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgWhitelist{}
	_ sdk.Msg = &MsgAddPool{}
	_ sdk.Msg = &MsgRemovePool{}
	_ sdk.Msg = &MsgDewhitelist{}
	_ sdk.Msg = &MsgClaimRewards{}
	_ sdk.Msg = &MsgClaimAllUserRewards{}
)

func NewMsgClose(creator string, id uint64, amount math.Int) *MsgClose {
	return &MsgClose{
		Creator:  creator,
		Id:       id,
		LpAmount: amount,
	}
}

func (msg *MsgClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.LpAmount.IsNil() {
		return errors.New("invalid lp amount: cannot be nil")
	}
	if msg.LpAmount.IsNegative() || msg.LpAmount.IsZero() {
		return errors.New("invalid lp amount: cannot be zero or negative")
	}
	return nil
}

func NewMsgOpen(creator string, collateralAsset string, collateralAmount math.Int, ammPoolId uint64, leverage math.LegacyDec, stopLossPrice math.LegacyDec) *MsgOpen {
	return &MsgOpen{
		Creator:          creator,
		CollateralAsset:  collateralAsset,
		CollateralAmount: collateralAmount,
		AmmPoolId:        ammPoolId,
		Leverage:         leverage,
		StopLossPrice:    stopLossPrice,
	}
}

func (msg *MsgOpen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// leverage should be greater than or equal to 1
	if msg.Leverage.LT(math.LegacyOneDec()) {
		return ErrLeverageTooSmall
	}
	collateralCoin := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	// sdk.NewCoin already coin.Validate(), but it does not check if amount is 0
	if collateralCoin.IsZero() {
		return ErrInvalidCollateralAsset.Wrapf("(amount cannot be equal to 0)")
	}

	// 0 StopLoss price is allowed. It means not set
	if msg.StopLossPrice.IsNegative() {
		return fmt.Errorf("stop loss price cannot be negative")
	}
	return nil
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if err = msg.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}
	return nil
}

func NewMsgUpdateParams(signer string, params *Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: signer,
		Params:    params,
	}
}

func NewMsgUpdateEnabledPools(signer string, enabledPools []uint64) *MsgUpdateEnabledPools {
	return &MsgUpdateEnabledPools{
		Authority:    signer,
		EnabledPools: enabledPools,
	}
}

func NewMsgWhitelist(signer string, whitelistedAddress string) *MsgWhitelist {
	return &MsgWhitelist{
		Authority:          signer,
		WhitelistedAddress: whitelistedAddress,
	}
}

func (msg *MsgWhitelist) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.WhitelistedAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid whitelist address (%s)", err)
	}
	return nil
}

func (msg *MsgAddPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err = utils.CheckLegacyDecNilAndNegative(msg.Pool.LeverageMax, "LeverageMax"); err != nil {
		return err
	}
	if msg.Pool.LeverageMax.LTE(math.LegacyOneDec()) {
		return ErrLeverageTooSmall
	}

	if err = utils.CheckLegacyDecNilAndNegative(msg.Pool.PoolMaxLeverageRatio, "PoolMaxLeverageRatio"); err != nil {
		return err
	}
	if err = utils.CheckLegacyDecNilAndNegative(msg.Pool.AdlTriggerRatio, "AdlTriggerRatio"); err != nil {
		return err
	}
	if !msg.Pool.PoolMaxLeverageRatio.GT(math.LegacyZeroDec()) || !msg.Pool.PoolMaxLeverageRatio.LT(math.LegacyOneDec()) {
		return errors.New("invalid pool max leverage ratio")
	}
	if msg.Pool.AdlTriggerRatio.LT(msg.Pool.PoolMaxLeverageRatio) {
		return errors.New("adl trigger ratio must be greater than max leverage ratio")
	}
	return nil
}

func NewMsgAddPool(signer string, pool AddPool) *MsgAddPool {

	return &MsgAddPool{
		Authority: signer,
		Pool:      pool,
	}
}

func (msg *MsgRemovePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgRemovePool(signer string, poolId uint64) *MsgRemovePool {
	return &MsgRemovePool{
		Authority: signer,
		Id:        poolId,
	}
}

func NewMsgDewhitelist(signer string, whitelistedAddress string) *MsgDewhitelist {
	return &MsgDewhitelist{
		Authority:          signer,
		WhitelistedAddress: whitelistedAddress,
	}
}

func (msg *MsgDewhitelist) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.WhitelistedAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid whitelisted address (%s)", err)
	}
	return nil
}

func NewMsgClaimRewards(signer string, ids []uint64) *MsgClaimRewards {
	return &MsgClaimRewards{
		Sender: signer,
		Ids:    ids,
	}
}

func (msg *MsgClaimRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if len(msg.Ids) == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "empty ids")
	}

	poolIdsMap := make(map[uint64]bool)
	for _, id := range msg.Ids {
		if poolIdsMap[id] {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate pool id %d", id)
		} else {
			poolIdsMap[id] = true
		}
	}
	return nil
}

func NewMsgClaimAllUserRewards(signer string) *MsgClaimAllUserRewards {
	return &MsgClaimAllUserRewards{
		Sender: signer,
	}
}

func (msg *MsgClaimAllUserRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func (msg *MsgUpdatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.PoolId == 0 {
		return errors.New("invalid pool id")
	}

	if msg.LeverageMax.IsNil() || msg.LeverageMax.LTE(math.LegacyOneDec()) {
		return errors.New("invalid leverage max")
	}

	if msg.MaxLeveragelpRatio.IsNil() || msg.MaxLeveragelpRatio.IsNegative() {
		return errors.New("invalid max leverage ratio")
	}
	return nil
}

func (msg *MsgUpdateEnabledPools) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	poolIdsMap := make(map[uint64]bool)
	for _, id := range msg.EnabledPools {
		if poolIdsMap[id] {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate pool id %d", id)
		} else {
			poolIdsMap[id] = true
		}
	}
	return nil
}
