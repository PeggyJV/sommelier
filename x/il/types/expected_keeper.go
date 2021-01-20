package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// OracleKeeper is expected keeper for the oracle module
type OracleKeeper interface {
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingtypes.ValidatorI // get validator by operator address; nil when validator not found
	TotalBondedTokens(sdk.Context) sdk.Int                                     // total bonded tokens within the validator set
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec)                 // slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Jail(sdk.Context, sdk.ConsAddress)                                         // jail a validator
	IterateValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
}

// EthBridgeKeeper is expected keeper for the peggy bridge module
type EthBridgeKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins)
}

// // AccountKeeper is the expected account keeper
// type AccountKeeper interface {
// 	GetAccount(ctx sdk.Context, address sdk.AccAddress) authtypes.AccountI
// 	GetModuleAddress(name string) sdk.AccAddress
// 	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
// 	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
// }

// // BankKeeper is expected bank keeper
// type BankKeeper interface {
// 	GetSupply(ctx sdk.Context) (supply bankexported.SupplyI)
// 	SetSupply(ctx sdk.Context, supply bankexported.SupplyI)
// 	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
// 	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
// 	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
// }
