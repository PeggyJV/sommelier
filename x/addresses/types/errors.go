package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidEvmAddress    = sdkerrors.Register(ModuleName, 2, "invalid evm address")
	ErrInvalidCosmosAddress = sdkerrors.Register(ModuleName, 3, "invalid cosmos address")
)
