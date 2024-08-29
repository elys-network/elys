package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgOpen         = "open"
	TypeMsgClose        = "close"
	TypeMsgUpdateParams = "update_params"
	TypeMsgWhitelist    = "whitelist"
	TypeMsgAddPool      = "add_pool"
	TypeMsgRemovePool   = "remove_pool"
	TypeMsgUpdatePool   = "update_pool"
	TypeMsgDewhitelist  = "dewhitelist"
	TypeMsgClaimRewards = "claim_rewards"
)

var (
	_ sdk.Msg = &MsgClose{}
	_ sdk.Msg = &MsgOpen{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgWhitelist{}
	_ sdk.Msg = &MsgAddPool{}
	_ sdk.Msg = &MsgUpdatePool{}
	_ sdk.Msg = &MsgRemovePool{}
	_ sdk.Msg = &MsgDewhitelist{}
	_ sdk.Msg = &MsgClaimRewards{}
)

func NewMsgClose(creator string, id uint64, amount math.Int) *MsgClose {
	return &MsgClose{
		Creator:  creator,
		Id:       id,
		LpAmount: amount,
	}
}

func (msg *MsgClose) Route() string {
	return RouterKey
}

func (msg *MsgClose) Type() string {
	return TypeMsgClose
}

func (msg *MsgClose) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClose) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.LpAmount.IsNegative() {
		return fmt.Errorf("invalid lp amount: cannot be negative")
	}
	return nil
}

func NewMsgOpen(creator string, collateralAsset string, collateralAmount math.Int, ammPoolId uint64, leverage sdk.Dec, stopLossPrice sdk.Dec) *MsgOpen {
	return &MsgOpen{
		Creator:          creator,
		CollateralAsset:  collateralAsset,
		CollateralAmount: collateralAmount,
		AmmPoolId:        ammPoolId,
		Leverage:         leverage,
		StopLossPrice:    stopLossPrice,
	}
}

func (msg *MsgOpen) Route() string {
	return RouterKey
}

func (msg *MsgOpen) Type() string {
	return TypeMsgOpen
}

func (msg *MsgOpen) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgOpen) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// leverage should be greater than 1
	if msg.Leverage.LTE(sdk.OneDec()) {
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

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
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

func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

func NewMsgWhitelist(signer string, whitelistedAddress string) *MsgWhitelist {
	return &MsgWhitelist{
		Authority:          signer,
		WhitelistedAddress: whitelistedAddress,
	}
}

func (msg *MsgWhitelist) Route() string { return RouterKey }

func (msg *MsgWhitelist) Type() string {
	return TypeMsgWhitelist
}

func (msg *MsgWhitelist) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgWhitelist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
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

func (msg *MsgAddPool) Route() string {
	return RouterKey
}

func (msg *MsgAddPool) Type() string { return TypeMsgAddPool }

func (msg *MsgAddPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Pool.LeverageMax.LTE(sdk.OneDec()) {
		return ErrLeverageTooSmall
	}
	return nil
}

func NewMsgAddPool(signer string, pool AddPool) *MsgAddPool {

	return &MsgAddPool{
		Authority: signer,
		Pool:      pool,
	}
}

func (msg *MsgRemovePool) Route() string {
	return RouterKey
}

func (msg *MsgRemovePool) Type() string {
	return TypeMsgRemovePool
}

func (msg *MsgRemovePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemovePool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
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

func NewMsgUpdatePool(signer string, pool UpdatePool) *MsgUpdatePool {
	return &MsgUpdatePool{
		Authority:  signer,
		UpdatePool: &pool,
	}
}

func (msg *MsgUpdatePool) Route() string { return RouterKey }

func (msg *MsgUpdatePool) Type() string { return TypeMsgUpdatePool }

func (msg *MsgUpdatePool) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgDewhitelist(signer string, whitelistedAddress string) *MsgDewhitelist {
	return &MsgDewhitelist{
		Authority:          signer,
		WhitelistedAddress: whitelistedAddress,
	}
}

func (msg *MsgDewhitelist) Route() string {
	return RouterKey
}

func (msg *MsgDewhitelist) Type() string {
	return TypeMsgDewhitelist
}

func (msg *MsgDewhitelist) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgDewhitelist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
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

func (msg *MsgClaimRewards) Route() string {
	return RouterKey
}

func (msg *MsgClaimRewards) Type() string {
	return TypeMsgClaimRewards
}

func (msg *MsgClaimRewards) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgClaimRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
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
