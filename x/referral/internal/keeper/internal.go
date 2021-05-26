package keeper

import (
	"github.com/anathatech/project-anatha/x/referral/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAddressBalance(ctx sdk.Context, address sdk.AccAddress) (types.Balance, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetBalanceKey(address))
	if bz == nil {
		return types.NewBalance(), true
	}

	var balance types.Balance
	k.cdc.MustUnmarshalBinaryBare(bz, &balance)

	return balance, true
}

func (k Keeper) SetAddressBalance(ctx sdk.Context, address sdk.AccAddress, balance types.Balance) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetBalanceKey(address), k.cdc.MustMarshalBinaryBare(balance))
}

func (k Keeper) GetAddressBalanceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.BalancePrefix)
}

func (k Keeper) IterateAddressBalances(ctx sdk.Context, cb func(address sdk.AccAddress, balance types.Balance) (stop bool)) {
	iterator := k.GetAddressBalanceIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := sdk.AccAddress(iterator.Key()[1:])
		var balance types.Balance

		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &balance)

		if cb(address, balance) {
			break
		}
	}
}

func (k Keeper) GetAddressParent(ctx sdk.Context, address sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetAddressParentKey(address))
	if bz == nil {
		return k.SupplyKeeper.GetModuleAddress(types.CharityFundModuleName), true
	}

	var parent sdk.AccAddress
	k.cdc.MustUnmarshalBinaryBare(bz, &parent)

	return parent, true
}

func (k Keeper) SetAddressParent(ctx sdk.Context, address sdk.AccAddress, parent sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetAddressParentKey(address), k.cdc.MustMarshalBinaryBare(parent))
}

func (k Keeper) GetAddressParentIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AddressParentPrefix)
}

func (k Keeper) IterateAddressParents(ctx sdk.Context, cb func(address sdk.AccAddress, parent sdk.AccAddress) (stop bool)) {
	iterator := k.GetAddressParentIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := sdk.AccAddress(iterator.Key()[1:])
		var parent sdk.AccAddress

		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &parent)

		if cb(address, parent) {
			break
		}
	}
}

func (k Keeper) GetAddressChildren(ctx sdk.Context, address sdk.AccAddress) ([]sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetAddressChildrenKey(address))
	if bz == nil {
		return []sdk.AccAddress{}, true
	}

	var children []sdk.AccAddress
	k.cdc.MustUnmarshalBinaryBare(bz, &children)

	return children, true
}

func (k Keeper) SetAddressChildren(ctx sdk.Context, address sdk.AccAddress, children []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetAddressChildrenKey(address), k.cdc.MustMarshalBinaryBare(children))
}

func (k Keeper) GetAddressChildrenIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AddressChildrenPrefix)
}

func (k Keeper) IterateAddressChildren(ctx sdk.Context, cb func(address sdk.AccAddress, children []sdk.AccAddress) (stop bool)) {
	iterator := k.GetAddressChildrenIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := sdk.AccAddress(iterator.Key()[1:])
		var children []sdk.AccAddress

		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &children)

		if cb(address, children) {
			break
		}
	}
}