package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


var (
	ErrNoBalance 				= sdkerrors.Register(ModuleName, 101, "No pending balance available.")
)
