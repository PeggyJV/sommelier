package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// Keeper of the PoA store.
type Keeper struct {
	cdc            codec.BinaryCodec
	storeKey       storetypes.StoreKey
	paramSpace     paramtypes.Subspace
	sk             types.StakingKeeper
	slashingKeeper types.SlashingKeeper
	// authority is the bech32 address allowed to submit MsgUpdateAuthoritySet
	// and MsgUpdateParams. Typically the gov module account.
	authority string
}

// NewKeeper builds a new PoA Keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ps paramtypes.Subspace,
	sk types.StakingKeeper,
	slashingKeeper types.SlashingKeeper,
	authority string,
) Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:            cdc,
		storeKey:       key,
		paramSpace:     ps,
		sk:             sk,
		slashingKeeper: slashingKeeper,
		authority:      authority,
	}
}

// Logger returns a module-scoped logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Authority returns the configured gov authority bech32 address.
func (k Keeper) Authority() string { return k.authority }

// StakingKeeper returns the underlying (non-wrapped) staking keeper. Intended
// for use by the EndBlocker, which needs to drive staking's own state machine.
func (k Keeper) StakingKeeper() types.StakingKeeper { return k.sk }

// GetParams returns the current PoA params.
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	k.paramSpace.GetParamSet(ctx, &p)
	return
}

// SetParams stores the PoA params.
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) {
	k.paramSpace.SetParamSet(ctx, &p)
}

// SetSlashingKeeper wires the slashing-keeper dependency post-construction.
// PoA's keeper is created BEFORE slashing in app.go (so slashing's NewKeeper
// can take the wrapped staking keeper); this setter fills in the back-reference
// once slashing has been built. The slashing keeper is only consulted in the
// EndBlocker for snapshot retention.
func (k *Keeper) SetSlashingKeeper(sk types.SlashingKeeper) {
	k.slashingKeeper = sk
}
