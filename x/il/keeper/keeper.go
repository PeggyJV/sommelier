package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/peggyjv/sommelier/x/il/types"
)

// Keeper of the impermanent store
type Keeper struct {
	cdc        codec.BinaryMarshaler
	storeKey   sdk.StoreKey
	paramSpace paramtypes.Subspace

	oracleKeeper    types.OracleKeeper
	ethBridgeKeeper types.EthBridgeKeeper
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey,
	paramspace paramtypes.Subspace, oracleKeeper types.OracleKeeper,
	ethBridgeKeeper types.EthBridgeKeeper) Keeper {

	// set KeyTable if it has not already been set
	if !paramspace.HasKeyTable() {
		paramspace = paramspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		paramSpace:      paramspace,
		oracleKeeper:    oracleKeeper,
		ethBridgeKeeper: ethBridgeKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
