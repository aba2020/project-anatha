package cli

import (
	"bufio"
	"github.com/anathatech/project-anatha/x/governance"
	referralutils "github.com/anathatech/project-anatha/x/referral/client/utils"
	"github.com/anathatech/project-anatha/x/referral/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

func GetCmdSubmitCharityFundDistributionProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "charity-fund-distribution [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a Charity Fund Distribution proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			proposal, err := referralutils.ParseCharityFundDistributionProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewCharityFundDistributionProposal(
				proposal.Title,
				proposal.Description,
				proposal.Amount,
				proposal.Recipient,
			)

			msg := governance.NewMsgSubmitProposal(content, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
