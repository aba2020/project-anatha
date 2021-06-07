package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Balance struct {
	PendingReward sdk.Coins		`json:"pending_reward" yaml:"pending_reward"`
	WithdrawnReward sdk.Coins	`json:"withdrawn_reward" yaml:"withdrawn_reward"`
}

func NewBalance() Balance {
	return Balance{
		PendingReward: sdk.NewCoins(),
		WithdrawnReward: sdk.NewCoins(),
	}
}

func (h Balance) String() string {
 	return fmt.Sprintf(`Pending reward: %s
Withdrawn reward: %s`, h.PendingReward, h.WithdrawnReward)
}

type AddressParentPair struct {
	Address sdk.AccAddress
	Parent sdk.AccAddress
}

func NewAddressParentPair(address sdk.AccAddress, parent sdk.AccAddress) AddressParentPair {
	return AddressParentPair{
		Address: address,
		Parent: parent,
	}
}

func (h AddressParentPair) String() string {
	return fmt.Sprintf(`Address: %s
Parent: %s`, h.Address, h.Parent)
}

type AddressChildPair struct {
	Address sdk.AccAddress
	Child sdk.AccAddress
}

func NewAddressChildrenPair(address sdk.AccAddress, child sdk.AccAddress) AddressChildPair {
	return AddressChildPair{
		Address: address,
		Child: child,
	}
}

func (h AddressChildPair) String() string {
	return fmt.Sprintf(`Address: %s
Child: %s`, h.Address, h.Child)
}

type AddressBalancePair struct {
	Address sdk.AccAddress
	Balance Balance
}

func NewAddressBalancePair(address sdk.AccAddress, balance Balance) AddressBalancePair {
	return AddressBalancePair{
		Address: address,
		Balance: balance,
	}
}

func (h AddressBalancePair) String() string {
	return fmt.Sprintf(`Address: %s
Balance: %s`, h.Address, h.Balance)
}

type ChildBalance struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Balance Balance `json:"balance" yaml:"balance"`
}

func NewChildBalance(address sdk.AccAddress, balance Balance) ChildBalance {
	return ChildBalance{
		Address: address,
		Balance: balance,
	}
}

func (h ChildBalance) String() string {
	return fmt.Sprintf(`Address: %s
Balance: %s`, h.Address, h.Balance)
}