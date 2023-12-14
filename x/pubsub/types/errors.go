package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/pubsub module sentinel errors
var (
	ErrAlreadyExists = errorsmod.Register(ModuleName, 2, "entity already exists")
	ErrInvalid       = errorsmod.Register(ModuleName, 3, "entity is invalid")
)
