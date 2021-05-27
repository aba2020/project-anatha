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
	var children []sdk.AccAddress

	iterator := k.GetAddressChildrenIterator(ctx, address)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		child := sdk.AccAddress(iterator.Key()[21:])

		children = append(children, child)
	}

	return children, true
}

func (k Keeper) SetAddressChild(ctx sdk.Context, address sdk.AccAddress, child sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetAddressChildKey(address, child), types.StatusPresent)
}

func (k Keeper) HasAddressChild(ctx sdk.Context, address sdk.AccAddress, child sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.GetAddressChildKey(address, child))
}

func (k Keeper) GetAddressChildIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AddressChildPrefix)
}

func (k Keeper) GetAddressChildrenIterator(ctx sdk.Context, address sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := append(types.AddressChildPrefix, address...)

	return sdk.KVStorePrefixIterator(store, key)
}

func (k Keeper) IterateAddressChild(ctx sdk.Context, cb func(address sdk.AccAddress, child sdk.AccAddress) (stop bool)) {
	iterator := k.GetAddressChildIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := sdk.AccAddress(iterator.Key()[1:21])
		child := sdk.AccAddress(iterator.Key()[21:])

		if cb(address, child) {
			break
		}
	}
}
