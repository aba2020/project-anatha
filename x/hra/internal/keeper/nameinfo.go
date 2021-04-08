package keeper

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/hra/internal/types"
	"github.com/anathatech/project-anatha/x/referral"
)


func (k Keeper) HandleRegisterName(ctx sdk.Context, name string, owner sdk.AccAddress) error {
	if k.IsNameRegistered(ctx, name) {
		return types.ErrNameRegistered
	}

	err := k.SupplyKeeper.SendCoinsFromAccountToModule(
		ctx,
		owner,
		k.feeCollectorName,
		k.NameInfoRegistrationFee(ctx),
	)
	if err != nil {
		return err
	}

	if ! k.OwnsAnyName(ctx, owner) {
		k.SetCredits(ctx, owner, k.AddressCredits(ctx))
		k.AfterFirstNameCreated(ctx, owner)
	}

	k.RegisterName(ctx, name, owner)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegister,
			sdk.NewAttribute(types.AttributeKeyName, name),
			sdk.NewAttribute(types.AttributeKeyExpires, ctx.BlockTime().Add(k.NameInfoDuration(ctx)).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, owner.String()),
		),
	})

	return nil
}

func (k Keeper) HandleRegisterNameV2(ctx sdk.Context, name string, owner sdk.AccAddress, referrer sdk.AccAddress) error {
	if k.IsNameRegistered(ctx, name) {
		return types.ErrNameRegistered
	}

	fee := k.NameInfoRegistrationFee(ctx)

	if owner.Equals(referrer) {
		return types.ErrNoSelfReferral
	}

	feeAmount := fee.AmountOf(config.DefaultDenom)

	referralAmount := feeAmount.ToDec().Mul(k.ReferralKeeper.ReferralPercentage(ctx)).TruncateInt()
	referralCoin := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, referralAmount))

	fee = fee.Sub(referralCoin)

	err := k.SupplyKeeper.SendCoinsFromAccountToModule(
		ctx,
		owner,
		k.feeCollectorName,
		fee,
	)
	if err != nil {
		return err
	}

	err = k.SupplyKeeper.SendCoinsFromAccountToModule(
		ctx,
		owner,
		referral.ReferralModuleName,
		referralCoin,
	)
	if err != nil {
		return err
	}

	existingParent, _ := k.ReferralKeeper.GetAddressParent(ctx, owner)
	ownerChildren, _ := k.ReferralKeeper.GetAddressChildren(ctx, owner)

	eventReferral := referrer
	if existingParent.Equals(k.SupplyKeeper.GetModuleAddress(referral.CharityFundModuleName)) && len(ownerChildren) == 0 {
		k.ReferralKeeper.SetAddressParent(ctx, owner, referrer)
		currentChildren, _ := k.ReferralKeeper.GetAddressChildren(ctx, referrer)

		found := false

		for _, el := range currentChildren {
			if el.Equals(owner) {
				found = true

				break
			}
		}

		if !found {
			currentChildren = append(currentChildren, owner)
			k.ReferralKeeper.SetAddressChildren(ctx, referrer, currentChildren)
		}

		referralBalance, _ := k.ReferralKeeper.GetAddressBalance(ctx, referrer)
		referralBalance.PendingReward = referralBalance.PendingReward.Add(referralCoin...)
		k.ReferralKeeper.SetAddressBalance(ctx, referrer, referralBalance)
	} else {
		referralBalance, _ := k.ReferralKeeper.GetAddressBalance(ctx, existingParent)
		referralBalance.PendingReward = referralBalance.PendingReward.Add(referralCoin...)
		k.ReferralKeeper.SetAddressBalance(ctx, existingParent, referralBalance)

		eventReferral = existingParent
	}

	if ! k.OwnsAnyName(ctx, owner) {
		k.SetCredits(ctx, owner, k.AddressCredits(ctx))
		k.AfterFirstNameCreated(ctx, owner)
	}

	k.RegisterName(ctx, name, owner)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterV2,
			sdk.NewAttribute(types.AttributeKeyName, name),
			sdk.NewAttribute(types.AttributeKeyExpires, ctx.BlockTime().Add(k.NameInfoDuration(ctx)).String()),
			sdk.NewAttribute(types.AttributeKeyReferral, eventReferral.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, owner.String()),
		),
	})

	return nil
}

func (k Keeper) RegisterName(ctx sdk.Context, name string, owner sdk.AccAddress) {
	nameInfo := types.NewNameInfo(name)

	nameInfo.Owner = owner
	nameInfo.CreationTime = ctx.BlockTime()
	nameInfo.ExpiryTime = nameInfo.CreationTime.Add(k.NameInfoDuration(ctx))

	k.SetNameInfo(ctx, name, nameInfo)
	k.SetNameInfoStatusMap(ctx, owner, name)

	k.InsertExpiredNameInfoQueue(ctx, name, nameInfo.ExpiryTime)
}

func (k Keeper) HandleRenewName(ctx sdk.Context, name string, owner sdk.AccAddress) error {
	nameInfo, found := k.GetNameInfo(ctx, name)
	if ! found {
		return types.ErrNameNotRegistered
	}

	if ! owner.Equals(nameInfo.Owner) {
		return types.ErrNotOwner
	}

	err := k.SupplyKeeper.SendCoinsFromAccountToModule(
		ctx,
		owner,
		k.feeCollectorName,
		k.NameInfoRenewalFee(ctx),
	)
	if err != nil {
		return err
	}

	oldExpiryTime := nameInfo.ExpiryTime

	if ctx.BlockTime().Before(nameInfo.ExpiryTime) {
		// If renewal is done before the expiry time, renew for 1 year after the expiry time
		nameInfo.ExpiryTime = nameInfo.ExpiryTime.Add(k.NameInfoDuration(ctx))
	} else {
		return types.ErrExpiredNameRenewal
	}

	if nameInfo.ExpiryTime.After(ctx.BlockTime().Add(k.NameInfoMaxDuration(ctx))) {
		return types.ErrMaximumDurationExceeded
	}

	k.RemoveFromExpiredNameInfoQueue(ctx, name, oldExpiryTime)
	k.InsertExpiredNameInfoQueue(ctx, name, nameInfo.ExpiryTime)

	k.SetNameInfo(ctx, name, nameInfo)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRenew,
			sdk.NewAttribute(types.AttributeKeyName, name),
			sdk.NewAttribute(types.AttributeKeyExpires, nameInfo.ExpiryTime.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, owner.String()),
		),
	})

	return nil
}

func (k Keeper) IsNameRegistered(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.GetNameInfoByNameKey(name))
}

func (k Keeper) OwnsAnyName(ctx sdk.Context, owner sdk.AccAddress) bool {
	iterator := k.GetNamesByAddressIterator(ctx, owner)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if bytes.Compare(iterator.Value(), types.StatusPresent) == 0 {
			return true
		}
	}

	return false
}

func (k Keeper) HandleDeleteName(ctx sdk.Context, name string, owner sdk.AccAddress) error {
	nameInfo, found := k.GetNameInfo(ctx, name)
	if ! found {
		return types.ErrNameNotRegistered
	}

	if ! owner.Equals(nameInfo.Owner) {
		return types.ErrNotOwner
	}

	k.RemoveFromExpiredNameInfoQueue(ctx, nameInfo.Name, nameInfo.ExpiryTime)

	k.DeleteNameInfo(ctx, name)
	k.DeleteNameInfoStatusMap(ctx, owner, name)

	// if last HRA remove all associated addresses
	if ! k.OwnsAnyName(ctx, owner) {
		k.RemoveAllAddresses(ctx, owner)

		k.SetCredits(ctx, owner, sdk.ZeroInt())
		err := k.AfterLastNameRemoved(ctx, owner)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) HandleTransferName(ctx sdk.Context, name string, owner sdk.AccAddress, newOwner sdk.AccAddress) error {
	nameInfo, found := k.GetNameInfo(ctx, name)
	if ! found {
		return types.ErrNameNotRegistered
	}

	if ! owner.Equals(nameInfo.Owner) {
		return types.ErrNotOwner
	}

	if nameInfo.Owner.Equals(newOwner) {
		return types.ErrAlreadyOwned
	}

	account := k.AccountKeeper.GetAccount(ctx, newOwner)
	if account == nil {
		account = k.AccountKeeper.NewAccountWithAddress(ctx, newOwner)
		k.AccountKeeper.SetAccount(ctx, account)
	}

	if ! k.OwnsAnyName(ctx, newOwner) {
		k.SetCredits(ctx, newOwner, k.AddressCredits(ctx))
		k.AfterFirstNameCreated(ctx, newOwner)
	}

	// update the status mapping
	k.DeleteNameInfoStatusMap(ctx, nameInfo.Owner, name)
	k.SetNameInfoStatusMap(ctx, newOwner, name)

	// update the owner and reset the price
	nameInfo.Owner = newOwner
	nameInfo.Price = sdk.NewCoins(sdk.NewInt64Coin(config.DefaultDenom, 0))

	k.SetNameInfo(ctx, name, nameInfo)

	if ! k.OwnsAnyName(ctx, owner) {
		k.RemoveAllAddresses(ctx, owner)
		k.SetCredits(ctx, owner, sdk.ZeroInt())
		err := k.AfterLastNameRemoved(ctx, owner)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.NameInfoByNameKeyPrefix)
}
