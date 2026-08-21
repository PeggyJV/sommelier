package keeper

import (
	"encoding/binary"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
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

// SetActivationHeight records the height at which PoA became active, plus the
// block time at that height. Called once from InitGenesis (and the v10 upgrade
// handler); idempotent callers should guard with GetActivationHeight.
//
// The timestamp exists so pruneSnapshots can derive the chain's ACTUAL average
// block time (see estimateBlockNanos) instead of trusting a compiled-in
// constant. Under-estimating block count silently truncates snapshot retention,
// which makes authority slashes get skipped — a security-relevant failure that
// a hardcoded constant cannot detect.
func (k Keeper) SetActivationHeight(ctx sdk.Context, height int64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(height))
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ActivationHeightKey, bz)

	tbz := make([]byte, 8)
	binary.BigEndian.PutUint64(tbz, uint64(ctx.BlockTime().UnixNano()))
	store.Set(types.ActivationTimeKey, tbz)
}

// SetActivationStamp restores a previously exported activation height/time
// pair verbatim, rather than stamping the current block. Used by InitGenesis on
// an export/import round-trip so the activation height — which decides whether
// a missing slash snapshot is benign or corruption — survives the restart.
func (k Keeper) SetActivationStamp(ctx sdk.Context, height, unixNanos int64) {
	store := ctx.KVStore(k.storeKey)

	hbz := make([]byte, 8)
	binary.BigEndian.PutUint64(hbz, uint64(height))
	store.Set(types.ActivationHeightKey, hbz)

	tbz := make([]byte, 8)
	binary.BigEndian.PutUint64(tbz, uint64(unixNanos))
	store.Set(types.ActivationTimeKey, tbz)
}

// GetActivationTime returns the block time recorded at the activation height.
func (k Keeper) GetActivationTime(ctx sdk.Context) (int64, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.ActivationTimeKey)
	if bz == nil {
		return 0, false
	}
	return int64(binary.BigEndian.Uint64(bz)), true
}

// GetActivationHeight returns the PoA activation height and whether it has been
// set. A false found-bool means the module is not yet active (no height should
// be treated as post-activation).
func (k Keeper) GetActivationHeight(ctx sdk.Context) (int64, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.ActivationHeightKey)
	if bz == nil {
		return 0, false
	}
	return int64(binary.BigEndian.Uint64(bz)), true
}
