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
	AddressChildPrefix = []byte{0x12}

	StatusPresent = []byte{0x1}
)

func GetBalanceKey(address sdk.AccAddress) []byte {
	return append(BalancePrefix, address...)
}

func GetAddressParentKey(address sdk.AccAddress) []byte {
	return append(AddressParentPrefix, address...)
}

func GetAddressChildKey(address sdk.AccAddress, child sdk.AccAddress) []byte {
	key := append(AddressChildPrefix, address...)
	key = append(key, child...)

	return key
}
