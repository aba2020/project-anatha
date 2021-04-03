package referral

import (
	"github.com/anathatech/project-anatha/x/referral/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)

	for _, record := range data.AddressBalancePairs {
		keeper.SetAddressBalance(ctx, record.Address, record.Balance)
	}

	for _, record := range data.AddressChildrenPairs {
		keeper.SetAddressChildren(ctx, record.Address, record.Children)
	}

	for _, record := range data.AddressParentPairs {
		keeper.SetAddressParent(ctx, record.Address, record.Parent)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)

	var addressBalancePairs []types.AddressBalancePair
	k.IterateAddressBalances(ctx, func (address sdk.AccAddress, balance types.Balance) (stop bool) {
		addressBalancePairs = append(addressBalancePairs, types.NewAddressBalancePair(address, balance))

		return false
	})

	var addressParentPairs []types.AddressParentPair
	k.IterateAddressParents(ctx, func (address sdk.AccAddress, parent sdk.AccAddress) (stop bool) {
		addressParentPairs = append(addressParentPairs, types.NewAddressParentPair(address, parent))

		return false
	})

	var addressChildrenPairs []types.AddressChildrenPair
	k.IterateAddressChildren(ctx, func (address sdk.AccAddress, children []sdk.AccAddress) (stop bool) {
		addressChildrenPairs = append(addressChildrenPairs, types.NewAddressChildrenPair(address, children))

		return false
	})

	return GenesisState{
		Params:      params,
		AddressBalancePairs: addressBalancePairs,
		AddressParentPairs: addressParentPairs,
		AddressChildrenPairs: addressChildrenPairs,
	}
}