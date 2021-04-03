package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgWithdrawReferralRewards{}, "referral/WithdrawReferralRewards", nil)

	cdc.RegisterConcrete(CharityFundDistributionProposal{}, "referral/CharityFundDistributionProposal", nil)
}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)

	codec.RegisterCrypto(ModuleCdc)

	ModuleCdc.Seal()
}

