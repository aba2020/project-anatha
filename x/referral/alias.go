package referral

import (
	"github.com/anathatech/project-anatha/x/referral/internal/keeper"
	"github.com/anathatech/project-anatha/x/referral/internal/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
	ReferralModuleName = types.ReferralModuleName
	CharityFundModuleName = types.CharityFundModuleName
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewMsgWithdrawReferralRewards  = types.NewMsgWithdrawReferralRewards

	ModuleCdc     = types.ModuleCdc

	KeyNameReferralPercentage = types.KeyNameReferralPercentage

	ErrNoBalance = types.ErrNoBalance
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgWithdrawReferralRewards = types.MsgWithdrawReferralRewards
	CharityFundDistributionProposal	= types.CharityFundDistributionProposal
)
