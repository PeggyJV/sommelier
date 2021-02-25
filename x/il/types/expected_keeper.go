package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
)

// OracleKeeper is expected keeper for the oracle module
type OracleKeeper interface {
	IterateAggregatedOracleDataByHeight(ctx sdk.Context, cb func(oracleData oracletypes.OracleData) (stop bool))
}

// EthBridgeKeeper is expected keeper for the peggy bridge module
type EthBridgeKeeper interface {
	AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount, fee sdk.Coin) (uint64, error)
}

// // AccountKeeper is the expected account keeper
// type AccountKeeper interface {
// 	GetAccount(ctx sdk.Context, address sdk.AccAddress) authtypes.AccountI
// }
