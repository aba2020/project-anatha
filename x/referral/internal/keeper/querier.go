package keeper

import (
	"github.com/anathatech/project-anatha/x/referral/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryAddressChildren = "address-children"
	QueryAddressChildrenWithBalance = "address-children-with-balance"
	QueryAddressParent = "address-parent"
	QueryAddressBalance = "address-balance"
	QueryParameters = "parameters"
	QueryModule = "module"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryAddressChildren:
			return queryAddressChildren(ctx, path[1:], req, k)
		case QueryAddressChildrenWithBalance:
			return queryAddressChildrenWithBalance(ctx, path[1:], req, k)
		case QueryAddressParent:
			return queryAddressParent(ctx, path[1:], req, k)
		case QueryAddressBalance:
			return queryAddressBalance(ctx, path[1:], req, k)
		case QueryParameters:
			return queryParams(ctx, k)
		case QueryModule:
			return queryModuleAccount(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unknown referral query endpoint: %s", path[0])
		}
	}
}

func queryAddressChildren(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return nil, err
	}

	children, _ := k.GetAddressChildren(ctx, address)

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, children)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryAddressChildrenWithBalance(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return nil, err
	}

	children, _ := k.GetAddressChildren(ctx, address)

	var returnData []types.ChildBalance

	for i := 0; i < len(children); i++ {
		child := children[i]

		balance, _ := k.GetAddressBalance(ctx, child)

		returnData = append(returnData, types.NewChildBalance(child, balance))
	}

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, returnData)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryAddressParent(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return nil, err
	}

	parent, _ := k.GetAddressParent(ctx, address)

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, parent)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryAddressBalance(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return nil, err
	}

	balance, _ := k.GetAddressBalance(ctx, address)

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, balance)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryModuleAccount(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	key := supply.NewModuleAddress(path[0])

	acc := k.AccountKeeper.GetAccount(ctx, key)

	res, marshalErr := codec.MarshalJSONIndent(types.ModuleCdc, acc)

	if marshalErr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, marshalErr.Error())
	}

	return res, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}