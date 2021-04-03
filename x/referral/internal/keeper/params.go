package keeper

import (
	"github.com/anathatech/project-anatha/x/referral/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


// NameInfoDuration
func (k Keeper) ReferralPercentage(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, types.KeyNameReferralPercentage, &res)
	return
}

// GetParams returns the total set of referral parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the referral parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}