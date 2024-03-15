package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrInvalidEvmAddress    = errorsmod.Register(ModuleName, 2, "invalid evm address")
	ErrInvalidCosmosAddress = errorsmod.Register(ModuleName, 3, "invalid cosmos address")
)
