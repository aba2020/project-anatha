package types

import (
	"fmt"
	"github.com/anathatech/project-anatha/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	gov "github.com/anathatech/project-anatha/x/governance"
)

const (
	ProposalTypeCharityFundDistribution   = "CharityFundDistribution"
)

type CharityFundDistributionProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Amount 		sdk.Coins `json:"amount" yaml:"amount"`
	Recipient 	sdk.AccAddress `json:"recipient" yaml:"recipient"`
}

func NewCharityFundDistributionProposal(title, description string, amount sdk.Coins, recipient sdk.AccAddress) gov.Content {
	return CharityFundDistributionProposal{title, description, amount,recipient}
}

var _ gov.Content = CharityFundDistributionProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeCharityFundDistribution)
	gov.RegisterProposalTypeCodec(CharityFundDistributionProposal{}, "referral/CharityFundDistributionProposal")
}

func (p CharityFundDistributionProposal) GetTitle() string       { return p.Title }
func (p CharityFundDistributionProposal) GetDescription() string { return p.Description }
func (p CharityFundDistributionProposal) ProposalRoute() string  { return RouterKey }
func (p CharityFundDistributionProposal) ProposalType() string   { return ProposalTypeCharityFundDistribution }
func (p CharityFundDistributionProposal) ValidateBasic() error {
	if ! p.Amount.IsValid() || p.Amount.AmountOf(config.DefaultDenom).IsZero() {
		return sdkerrors.ErrInvalidCoins
	}

	return gov.ValidateAbstract(p)
}

func (sup CharityFundDistributionProposal) String() string {
	return fmt.Sprintf(`Charity Fund Distribution Proposal:
  Title: 		%s
  Description: 	%s
  Amount: 		%s
  Recipient: 	%s
`, sup.Title, sup.Description, sup.Amount, sup.Recipient)
}
