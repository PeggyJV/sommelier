package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/pubsub module sentinel errors
var (
	// TODO(bolten): fill out errors here
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
)
