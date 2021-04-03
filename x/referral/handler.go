package referral

import (
	"github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/referral/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgWithdrawReferralRewards:
			return handleWithdrawReferralRewards(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleWithdrawReferralRewards(ctx sdk.Context, msg MsgWithdrawReferralRewards, k Keeper) (*sdk.Result, error) {
	balance, ok := k.GetAddressBalance(ctx, msg.Address)

	if !ok {
		return nil, types.ErrNoBalance
	}

	amount := balance.PendingReward.AmountOf(config.DefaultDenom)

	if amount == sdk.ZeroInt() {
		return nil, types.ErrNoBalance
	}

	parentAmount := amount.ToDec().Mul(k.ReferralPercentage(ctx)).TruncateInt()
	withdrawAmount := amount.Sub(parentAmount)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(config.DefaultDenom, withdrawAmount))

	err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ReferralModuleName, msg.Address, withdrawCoins)

	if err != nil {
		return nil, err
	}

	balance.WithdrawnReward = balance.WithdrawnReward.Add(withdrawCoins...)
	balance.PendingReward = sdk.NewCoins()

	k.SetAddressBalance(ctx, msg.Address, balance)

	parent, _ := k.GetAddressParent(ctx, msg.Address)
	parentBalance, _ := k.GetAddressBalance(ctx, parent)

	parentBalance.PendingReward = parentBalance.PendingReward.Add(sdk.NewCoin(config.DefaultDenom, parentAmount))

	k.SetAddressBalance(ctx, parent, parentBalance)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawReferralRewards,
			sdk.NewAttribute(types.AttributeKeyCaller, msg.Caller.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
			sdk.NewAttribute(types.AttributeKeySender, msg.Caller.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func NewGovernanceProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
			case CharityFundDistributionProposal:
				return handleCharityFundDistribution(ctx, k, c)
			default:
					return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized referral proposal content type: %T", c)
		}
	}
}

func handleCharityFundDistribution(ctx sdk.Context, k Keeper, p types.CharityFundDistributionProposal) error {
	err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.CharityFundModuleName, p.Recipient, p.Amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCharityFundDistributionProposal,
			sdk.NewAttribute(types.AttributeKeyTitle, p.Title),
			sdk.NewAttribute(types.AttributeKeyDescription, p.Description),
			sdk.NewAttribute(types.AttributeKeyAmount, p.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, p.Recipient.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueModule),
		),
	)

	return nil
}