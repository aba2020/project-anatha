package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/spf13/cobra"

	"github.com/anathatech/project-anatha/x/referral/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	referralQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	referralQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdGetAddressChildren(queryRoute, cdc),
			GetCmdGetAddressParent(queryRoute, cdc),
			GetCmdGetAddressBalance(queryRoute, cdc),
			GetCmdQueryParams(queryRoute, cdc),
		)...,
	)

	referralQueryCmd.AddCommand(GetModuleAccountCmd(cdc))

	return referralQueryCmd

}

func GetCmdGetAddressChildren(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address-children [address]",
		Short: "Query all referral children of the given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/address-children/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("Could not query address children - %s \n", address)
				return nil
			}

			var out []sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetAddressParent(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address-parent [address]",
		Short: "Query referral parent of the given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/address-parent/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("Could not query address parent - %s \n", address)
				return nil
			}

			var out sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetAddressBalance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address-balance [address]",
		Short: "Query referral balances of the given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/address-balance/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("Could not query address balance - %s \n", address)
				return nil
			}

			var out types.Balance
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetModuleAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module [name]",
		Short: "Query module account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			accGetter := auth.NewAccountRetriever(cliCtx)

			key := supply.NewModuleAddress(args[0])

			acc, err := accGetter.GetAccount(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(acc)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/parameters", queryRoute)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}