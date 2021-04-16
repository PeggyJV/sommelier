package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bridgetypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
	oracletypes "github.com/peggyjv/sommelier/x/allocation/types"
)

// OracleKeeper is expected keeper for the oracle module
type OracleKeeper interface {
	GetLatestAggregatedOracleData(ctx sdk.Context, dataType, id string) (oracletypes.OracleData, int64)
}

// EthBridgeKeeper is expected keeper for the peggy bridge module
type EthBridgeKeeper interface {
	SetOutgoingLogicCall(ctx sdk.Context, call *bridgetypes.OutgoingLogicCall)
	GetLastObservedEthereumBlockHeight(ctx sdk.Context) bridgetypes.LastObservedEthereumBlockHeight
}
