package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/anathatech/project-anatha/x/referral/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	referralTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%S transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	referralTxCmd.AddCommand(flags.PostCommands(
		GetCmdWithdrawReferralRewards(cdc),
	)...)

	return referralTxCmd
}

func GetCmdWithdrawReferralRewards(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw-referral-rewards [address]",
		Short: "Withdraw pending referral rewards for the given address",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			address, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawReferralRewards(cliCtx.GetFromAddress(), address)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
