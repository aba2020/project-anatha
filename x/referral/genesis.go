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

	for _, record := range data.AddressChildPairs {
		keeper.SetAddressChild(ctx, record.Address, record.Child)
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

	var addressChildPairs []types.AddressChildPair
	k.IterateAddressChild(ctx, func (address sdk.AccAddress, child sdk.AccAddress) (stop bool) {
		addressChildPairs = append(addressChildPairs, types.NewAddressChildrenPair(address, child))

		return false
	})

	return GenesisState{
		Params:      params,
		AddressBalancePairs: addressBalancePairs,
		AddressParentPairs: addressParentPairs,
		AddressChildPairs: addressChildPairs,
	}
}