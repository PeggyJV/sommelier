package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/peggyjv/sommelier/v8/x/addresses/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ps paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   key,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

////////////
// Params //
////////////

// GetParamSet returns the vote period from the parameters
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramstore.GetParamSet(ctx, &p)
	return p
}

// setParams sets the parameters in the store
func (k Keeper) setParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

///////////////////////
/// AddressMappings ///
///////////////////////

// SetAddressMapping stores the mapping between the cosmos and evm addresses
func (k Keeper) SetAddressMapping(ctx sdk.Context, cosmosAddr []byte, evmAddr []byte) error {
	// sanity check, shouldn't be possible with proper validation in the message handler
	if len(cosmosAddr) == 0 {
		return errorsmod.Wrap(types.ErrNilCosmosAddress, "cosmos address cannot be empty")
	}

	if len(evmAddr) == 0 {
		return errorsmod.Wrap(types.ErrNilEvmAddress, "evm address cannot be empty")
	}

	k.setCosmosToEvmMapping(ctx, cosmosAddr, evmAddr)
	k.setEvmToCosmosMapping(ctx, evmAddr, cosmosAddr)

	return nil
}

func (k Keeper) setCosmosToEvmMapping(ctx sdk.Context, cosmosAddr []byte, evmAddr []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetCosmosToEvmMapKey(cosmosAddr), evmAddr)
}

func (k Keeper) setEvmToCosmosMapping(ctx sdk.Context, evmAddr []byte, cosmosAddr []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetEvmToCosmosMapKey(evmAddr), cosmosAddr)
}

// DeleteAddressMapping deletes the mapping between the cosmos and evm addresses
func (k Keeper) DeleteAddressMapping(ctx sdk.Context, cosmosAddr []byte) error {
	// sanity check, shouldn't be possible with proper validation in the message handler
	if len(cosmosAddr) == 0 {
		return errorsmod.Wrap(types.ErrNilCosmosAddress, "cosmos address cannot be empty")
	}

	evmAddr := k.GetEvmAddressByCosmosAddress(ctx, cosmosAddr)

	if len(evmAddr) == 0 {
		return nil
	}

	k.deleteEvmToCosmosMapping(ctx, evmAddr)
	k.deleteCosmosToEvmMapping(ctx, cosmosAddr)

	return nil
}

func (k Keeper) deleteCosmosToEvmMapping(ctx sdk.Context, cosmosAddr []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetCosmosToEvmMapKey(cosmosAddr))
}

func (k Keeper) deleteEvmToCosmosMapping(ctx sdk.Context, evmAddr []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetEvmToCosmosMapKey(evmAddr))
}

func (k Keeper) GetCosmosAddressByEvmAddress(ctx sdk.Context, evmAddr []byte) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GetEvmToCosmosMapKey(evmAddr))
}

func (k Keeper) GetEvmAddressByCosmosAddress(ctx sdk.Context, cosmosAddr []byte) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GetCosmosToEvmMapKey(cosmosAddr))
}

// IterateAddressMappings iterates over all Cosmos to EVM address mappings
func (k Keeper) IterateAddressMappings(ctx sdk.Context, cb func(cosmosAddr []byte, evmAddr []byte) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetCosmosToEvmMapPrefix())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		cosmosAddr := iterator.Key()[1:]
		evmAddr := iterator.Value()
		if cb(cosmosAddr, evmAddr) {
			break
		}
	}
}
