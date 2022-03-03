package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/pubsub module sentinel errors
var (
	ErrAlreadyExists = sdkerrors.Register(ModuleName, 2, "entity already exists")
	ErrInvalid       = sdkerrors.Register(ModuleName, 3, "entity is invalid")
)
