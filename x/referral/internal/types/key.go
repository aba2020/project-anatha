package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "referral"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName

	Separator = ":"

	ReferralModuleName = "referral"
	CharityFundModuleName = "charity-fund"
)

// Keys for referral store
var (
	BalancePrefix         = []byte{0x10}
	AddressParentPrefix   = []byte{0x11}
	AddressChildrenPrefix = []byte{0x12}
)

func GetBalanceKey(address sdk.AccAddress) []byte {
	return append(BalancePrefix, address...)
}

func GetAddressParentKey(address sdk.AccAddress) []byte {
	return append(AddressParentPrefix, address...)
}

func GetAddressChildrenKey(address sdk.AccAddress) []byte {
	return append(AddressChildrenPrefix, address...)
}
