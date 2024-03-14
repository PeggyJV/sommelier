package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidEvmAddress = sdkerrors.Register(ModuleName, 2, "invalid evm address")
)
