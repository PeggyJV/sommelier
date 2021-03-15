package types

import (
	bridgetypes "github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
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
