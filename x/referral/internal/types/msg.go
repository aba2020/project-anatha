package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgWithdrawReferralRewards
type MsgWithdrawReferralRewards struct {
	Caller sdk.AccAddress `json:"caller" yaml:"caller"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

func NewMsgWithdrawReferralRewards(caller sdk.AccAddress, address sdk.AccAddress) MsgWithdrawReferralRewards {
	return MsgWithdrawReferralRewards{
		Caller: caller,
		Address: address,
	}
}

func (msg MsgWithdrawReferralRewards) Route() string { return RouterKey }

func (msg MsgWithdrawReferralRewards) Type() string { return "withdraw_referral_rewards" }

func (msg MsgWithdrawReferralRewards) ValidateBasic() error {
	if msg.Caller.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Caller.String())
	}

	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Address.String())
	}

	return nil
}

func (msg MsgWithdrawReferralRewards) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgWithdrawReferralRewards) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Caller}
}
