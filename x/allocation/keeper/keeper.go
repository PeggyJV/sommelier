package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper

	oracleHandler types.OracleHandler
	handlerSet    bool
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper,
) Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		stakingKeeper: stakingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
