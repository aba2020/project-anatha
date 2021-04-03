package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type GenesisState struct {
	Params 						Params 		`json:"params" yaml:"params"`
	AddressParentPairs 			[]AddressParentPair 	`json:"address_parent_pairs" yaml:"address_parent_pairs"`
	AddressChildrenPairs		[]AddressChildrenPair `json:"address_children_pairs" yaml:"address_children_pairs"`
	AddressBalancePairs         []AddressBalancePair    `json:"address_balance_pairs" yaml:"address_balance_pairs"`
}


func NewGenesisState(params Params, addressParentPairs []AddressParentPair, addressChildrenPairs []AddressChildrenPair, addressBalancePairs []AddressBalancePair) GenesisState {
	return GenesisState{
		Params: params,
		AddressParentPairs: addressParentPairs,
		AddressChildrenPairs: addressChildrenPairs,
		AddressBalancePairs: addressBalancePairs,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		AddressParentPairs: []AddressParentPair{},
		AddressChildrenPairs: []AddressChildrenPair{},
		AddressBalancePairs: []AddressBalancePair{},
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, record := range data.AddressParentPairs {
		if record.Address.Empty() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid AddressParentPair: Value: %s. Error: Missing address", record.Address)
		}
	}
	for _, record := range data.AddressChildrenPairs {
		if record.Address.Empty() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid AddressChildrenPair: Value: %s. Error: Missing address", record.Address)
		}
	}
	for _, record := range data.AddressBalancePairs {
		if record.Address.Empty() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid AddressBalancePair: Value: %s. Error: Missing address", record.Address)
		}
	}
	return nil
}