package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrInvalidEvmAddress    = errorsmod.Register(ModuleName, 2, "invalid evm address")
	ErrInvalidCosmosAddress = errorsmod.Register(ModuleName, 3, "invalid cosmos address")
	ErrNilCosmosAddress     = errorsmod.Register(ModuleName, 4, "cosmos address cannot be nil")
	ErrNilEvmAddress        = errorsmod.Register(ModuleName, 5, "evm address cannot be nil")
)
