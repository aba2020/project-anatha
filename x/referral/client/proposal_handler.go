package client

import (
	govclient "github.com/anathatech/project-anatha/x/governance/client"
	"github.com/anathatech/project-anatha/x/referral/client/cli"
)

var (
	CharityFundDistributionProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitCharityFundDistributionProposal)
)