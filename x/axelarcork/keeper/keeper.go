package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"
)

var _ porttypes.ICS4Wrapper = &Keeper{}

// Keeper of the oracle store
type Keeper struct {
	storeKey           storetypes.StoreKey
	cdc                codec.BinaryCodec
	paramSpace         paramtypes.Subspace
	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	stakingKeeper      types.StakingKeeper
	transferKeeper     types.TransferKeeper
	distributionKeeper types.DistributionKeeper
	gravityKeeper      types.GravityKeeper
	poaKeeper          types.PoaKeeper

	Ics4Wrapper types.ICS4Wrapper
}

// NewKeeper creates a new x/axelarcork Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper,
	transferKeeper types.TransferKeeper, distributionKeeper types.DistributionKeeper,
	wrapper types.ICS4Wrapper, gravityKeeper types.GravityKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		paramSpace:         paramSpace,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		stakingKeeper:      stakingKeeper,
		transferKeeper:     transferKeeper,
		distributionKeeper: distributionKeeper,
		gravityKeeper:      gravityKeeper,

		Ics4Wrapper: wrapper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetTransferKeeper sets the transferKeeper
func (k *Keeper) SetTransferKeeper(transferKeeper types.TransferKeeper) {
	k.transferKeeper = transferKeeper
}

// SetPoaKeeper wires the read-only PoA safe-mode dependency post-construction
// (PoA is constructed before axelarcork in app.go).
func (k *Keeper) SetPoaKeeper(poaKeeper types.PoaKeeper) {
	k.poaKeeper = poaKeeper
}

// inSafeMode reports whether the chain is in authority-empty safe mode, in which
// case axelarcork operations are frozen. A nil PoA keeper means no freeze.
func (k Keeper) inSafeMode(ctx sdk.Context) bool {
	return k.poaKeeper != nil && k.poaKeeper.SafeModeActive(ctx)
}

////////////
// Params //
////////////

// GetParamSet returns the vote period from the parameters
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	// using this direct method instead of k.paramSpace.SetParamSet because our
	// param validation should have happened on ValidateBasic, where it is
	// contingent on being "enabled"
	for _, pair := range params.ParamSetPairs() {
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()

		k.paramSpace.Set(ctx, pair.Key, v)
	}
}

/////////////////////
// Scheduled Corks //
/////////////////////

// SetAuthorityAxelarCork stores a cork scheduled by the cork authority for
// execution at blockHeight on chainID, and returns its ID.
func (k Keeper) SetAuthorityAxelarCork(ctx sdk.Context, chainID uint64, blockHeight uint64, cork types.AxelarCork) []byte {
	id := cork.IDHash(blockHeight)
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(
		types.GetAuthorityCorkKey(chainID, blockHeight, id, common.HexToAddress(cork.TargetContractAddress)),
		bz,
	)
	return id
}

// DeleteAuthorityAxelarCork removes a scheduled authority cork.
func (k Keeper) DeleteAuthorityAxelarCork(ctx sdk.Context, chainID uint64, blockHeight uint64, id []byte, contract common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetAuthorityCorkKey(chainID, blockHeight, id, contract))
}

// IterateAuthorityAxelarCorksByBlockHeight walks the authority corks targeting
// blockHeight on chainID.
func (k Keeper) IterateAuthorityAxelarCorksByBlockHeight(
	ctx sdk.Context,
	chainID uint64,
	blockHeight uint64,
	cb func(id []byte, contract common.Address, cork types.AxelarCork) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAuthorityCorkKeyByBlockHeightPrefix(chainID, blockHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cork types.AxelarCork
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		keyPair.Next(8) // trim chain id, filtered by the prefix
		keyPair.Next(8) // trim block height, filtered by the prefix
		// Copy: Next returns a slice into the buffer backed by iter.Key(), and
		// callers retain the id past iter.Close().
		id := make([]byte, 32)
		copy(id, keyPair.Next(32))
		contract := common.BytesToAddress(keyPair.Next(20))

		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if cb(id, contract, cork) {
			break
		}
	}
}

// IterateAllAuthorityAxelarCorks walks every queued authority cork on a chain,
// at any height. Used by ExportGenesis; the per-height iterator is the hot path.
func (k Keeper) IterateAllAuthorityAxelarCorks(ctx sdk.Context, chainID uint64, cb func(blockHeight uint64, id []byte, contract common.Address, cork types.AxelarCork) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAuthorityCorkKeyPrefix(chainID))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cork types.AxelarCork
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		keyPair.Next(8) // trim chain id, filtered by the prefix
		blockHeight := sdk.BigEndianToUint64(keyPair.Next(8))
		// Copy: Next returns a slice into the buffer backed by iter.Key().
		id := make([]byte, 32)
		copy(id, keyPair.Next(32))
		contract := common.BytesToAddress(keyPair.Next(20))

		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if cb(blockHeight, id, contract, cork) {
			break
		}
	}
}

// GetAuthorityAxelarCorksByBlockHeight returns queued authority corks for one
// chain at one height.
func (k Keeper) GetAuthorityAxelarCorksByBlockHeight(ctx sdk.Context, chainID uint64, height uint64) []*types.ScheduledAxelarCork {
	var out []*types.ScheduledAxelarCork
	k.IterateAuthorityAxelarCorksByBlockHeight(ctx, chainID, height, func(id []byte, _ common.Address, cork types.AxelarCork) (stop bool) {
		c := cork
		out = append(out, &types.ScheduledAxelarCork{Cork: &c, BlockHeight: height, Id: hex.EncodeToString(id)})
		return false
	})
	return out
}

// GetAuthorityAxelarCorksByID returns queued authority corks matching an ID.
func (k Keeper) GetAuthorityAxelarCorksByID(ctx sdk.Context, chainID uint64, queriedID []byte) []*types.ScheduledAxelarCork {
	var out []*types.ScheduledAxelarCork
	k.IterateAllAuthorityAxelarCorks(ctx, chainID, func(blockHeight uint64, id []byte, _ common.Address, cork types.AxelarCork) (stop bool) {
		if bytes.Equal(id, queriedID) {
			c := cork
			out = append(out, &types.ScheduledAxelarCork{Cork: &c, BlockHeight: blockHeight, Id: hex.EncodeToString(id)})
		}
		return false
	})
	return out
}

// GetAuthorityAxelarCorks returns every queued authority cork on a chain for
// genesis export. The Validator field is left empty: authority corks have no
// scheduling validator. It is retained on the type for wire compatibility.
func (k Keeper) GetAuthorityAxelarCorks(ctx sdk.Context, chainID uint64) []*types.ScheduledAxelarCork {
	var out []*types.ScheduledAxelarCork
	k.IterateAllAuthorityAxelarCorks(ctx, chainID, func(blockHeight uint64, id []byte, _ common.Address, cork types.AxelarCork) (stop bool) {
		c := cork
		out = append(out, &types.ScheduledAxelarCork{
			Cork:        &c,
			BlockHeight: blockHeight,
			Id:          hex.EncodeToString(id),
		})
		return false
	})
	return out
}

/////////////////
// WinningCork //
/////////////////

func (k Keeper) SetWinningAxelarCork(ctx sdk.Context, chainID uint64, blockHeight uint64, cork types.AxelarCork) {
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(types.GetWinningAxelarCorkKey(chainID, blockHeight, common.HexToAddress(cork.TargetContractAddress)), bz)
}

func (k Keeper) IterateWinningAxelarCorks(ctx sdk.Context, chainID uint64, cb func(contract common.Address, blockHeight uint64, cork types.AxelarCork) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetWinningAxelarCorkKeyPrefix(chainID))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cork types.AxelarCork
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		keyPair.Next(8) // trim chain ID
		blockHeight := binary.BigEndian.Uint64(keyPair.Next(8))
		contractAddress := common.BytesToAddress(keyPair.Next(20)) // contract

		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if cb(contractAddress, blockHeight, cork) {
			break
		}
	}
}

func (k Keeper) GetWinningAxelarCork(ctx sdk.Context, chainID uint64, contractAddr common.Address) (uint64, types.AxelarCork, bool) {
	var bh uint64
	var c types.AxelarCork
	found := false
	k.IterateWinningAxelarCorks(ctx, chainID, func(contract common.Address, blockHeight uint64, cork types.AxelarCork) (stop bool) {
		if bytes.Equal(contractAddr.Bytes(), contract.Bytes()) {
			bh = blockHeight
			c = cork
			found = true
			return true
		}

		return false
	})

	return bh, c, found
}

func (k Keeper) DeleteWinningAxelarCorkByBlockheight(ctx sdk.Context, chainID uint64, blockHeight uint64, cork types.AxelarCork) {
	ctx.KVStore(k.storeKey).Delete(types.GetWinningAxelarCorkKey(chainID, blockHeight, common.HexToAddress(cork.TargetContractAddress)))
}

// TODO (Collin): Need pruning logic. This method is unused.
func (k Keeper) DeleteWinningAxelarCork(ctx sdk.Context, chainID uint64, c types.AxelarCork) {

	k.IterateWinningAxelarCorks(ctx, chainID, func(contract common.Address, blockHeight uint64, cork types.AxelarCork) (stop bool) {
		if c.Equals(cork) {
			k.DeleteWinningAxelarCorkByBlockheight(ctx, chainID, blockHeight, cork)

			return true
		}

		return false
	})
}

///////////////////////////
// ScheduledBlockHeights //
///////////////////////////

func (k Keeper) GetScheduledBlockHeights(ctx sdk.Context, chainID uint64) []uint64 {
	var heights []uint64

	latestHeight := uint64(0)
	k.IterateAllAuthorityAxelarCorks(ctx, chainID, func(blockHeight uint64, _ []byte, _ common.Address, _ types.AxelarCork) (stop bool) {
		if blockHeight > latestHeight {
			heights = append(heights, blockHeight)
		}
		latestHeight = blockHeight
		return false
	})

	return heights
}

////////////////////////
// AxelarCork Results //
////////////////////////

func (k Keeper) SetAxelarCorkResult(ctx sdk.Context, chainID uint64, id []byte, corkResult types.AxelarCorkResult) {
	bz := k.cdc.MustMarshal(&corkResult)
	ctx.KVStore(k.storeKey).Set(types.GetAxelarCorkResultKey(chainID, id), bz)
}

func (k Keeper) GetAxelarCorkResult(ctx sdk.Context, chainID uint64, id []byte) (types.AxelarCorkResult, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAxelarCorkResultKey(chainID, id))
	if len(bz) == 0 {
		return types.AxelarCorkResult{}, false
	}

	var corkResult types.AxelarCorkResult
	k.cdc.MustUnmarshal(bz, &corkResult)
	return corkResult, true
}

func (k Keeper) DeleteAxelarCorkResult(ctx sdk.Context, chainID uint64, id []byte) {
	ctx.KVStore(k.storeKey).Delete(types.GetAxelarCorkResultKey(chainID, id))
}

// IterateAxelarCorkResults iterates over all cork results by chain ID
func (k Keeper) IterateAxelarCorkResults(ctx sdk.Context, chainID uint64, cb func(id []byte, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.AxelarCorkResult) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAxelarCorkResultPrefix(chainID))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var corkResult types.AxelarCorkResult
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		id := keyPair.Next(32)

		k.cdc.MustUnmarshal(iter.Value(), &corkResult)
		if cb(id, corkResult.BlockHeight, corkResult.Approved, corkResult.ApprovalPercentage, corkResult) {
			break
		}
	}
}

// GetAxelarCorkResults returns AxelarCorkResults
func (k Keeper) GetAxelarCorkResults(ctx sdk.Context, chainID uint64) []*types.AxelarCorkResult {
	var corkResults []*types.AxelarCorkResult
	k.IterateAxelarCorkResults(ctx, chainID, func(id []byte, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.AxelarCorkResult) (stop bool) {
		corkResults = append(corkResults, &corkResult)
		return false
	})

	return corkResults
}

///////////
// Votes //
///////////

/////////////
// Cellars //
/////////////

func (k Keeper) SetCellarIDs(ctx sdk.Context, chainID uint64, c types.CellarIDSet) {
	bz := k.cdc.MustMarshal(&c)
	// always sort before writing to the store
	cellarIDs := make([]string, 0, len(c.Ids))
	cellarIDs = append(cellarIDs, c.Ids...)
	sort.Strings(cellarIDs)
	c.Ids = cellarIDs
	ctx.KVStore(k.storeKey).Set(types.MakeCellarIDsKey(chainID), bz)
}

func (k Keeper) GetCellarIDs(ctx sdk.Context, chainID uint64) (cellars []common.Address) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MakeCellarIDsKey(chainID))

	var cids types.CellarIDSet
	k.cdc.MustUnmarshal(bz, &cids)

	for _, cid := range cids.Ids {
		cellars = append(cellars, common.HexToAddress(cid))
	}

	return cellars
}

func (k Keeper) HasCellarID(ctx sdk.Context, chainID uint64, address common.Address) (found bool) {
	found = false
	for _, id := range k.GetCellarIDs(ctx, chainID) {
		if id == address {
			found = true
			break
		}
	}

	return found
}

/////////////////////////////////
// Axelar Contract Call Nonces //
/////////////////////////////////

// SetAxelarContractCallNonce sets the nonce for the given chainID and address
func (k Keeper) SetAxelarContractCallNonce(ctx sdk.Context, chainID uint64, address string, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAxelarContractCallNonceKey(chainID, common.HexToAddress(address)), sdk.Uint64ToBigEndian(nonce))
}

// GetAxelarContractCallNonce returns the nonce for the given chainID and address, returning a zero if not found
func (k Keeper) GetAxelarContractCallNonce(ctx sdk.Context, chainID uint64, address string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAxelarContractCallNonceKey(chainID, common.HexToAddress(address)))
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// IncrementAxelarContractCallNonce increments the nonce for the given chainID and address
func (k Keeper) IncrementAxelarContractCallNonce(ctx sdk.Context, chainID uint64, address string) uint64 {
	nonce := k.GetAxelarContractCallNonce(ctx, chainID, address) + 1
	k.SetAxelarContractCallNonce(ctx, chainID, address, nonce)

	return nonce
}

// IterateAxelarContractCallNonces iterates over all axelar contract call nonces
func (k Keeper) IterateAxelarContractCallNonces(ctx sdk.Context, cb func(chainID uint64, address common.Address, nonce uint64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.AxelarContractCallNoncePrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		chainID := sdk.BigEndianToUint64(keyPair.Next(8))
		address := common.BytesToAddress(keyPair.Next(20))
		nonce := sdk.BigEndianToUint64(iter.Value())
		if cb(chainID, address, nonce) {
			break
		}
	}
}

/////////////////////////
// Axelar Upgrade Data //
/////////////////////////

// SetAxelarProxyUpgradeData sets the upgrade data for the given chainID
func (k Keeper) SetAxelarProxyUpgradeData(ctx sdk.Context, chainID uint64, upgradeData types.AxelarUpgradeData) {
	ud := k.cdc.MustMarshal(&upgradeData)
	ctx.KVStore(k.storeKey).Set(types.GetAxelarProxyUpgradeDataKey(chainID), ud)
}

// GetAxelarProxyUpgradeData returns the upgrade data for the given chainID, returning an empty payload if not found
func (k Keeper) GetAxelarProxyUpgradeData(ctx sdk.Context, chainID uint64) (types.AxelarUpgradeData, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetAxelarProxyUpgradeDataKey(chainID))
	if len(bz) == 0 {
		return types.AxelarUpgradeData{}, false
	}

	upgradeData := types.AxelarUpgradeData{}
	k.cdc.MustUnmarshal(bz, &upgradeData)

	return upgradeData, true
}

// DeleteAxelarProxyUpgradeData deletes the upgrade data for the given chainID
func (k Keeper) DeleteAxelarProxyUpgradeData(ctx sdk.Context, chainID uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetAxelarProxyUpgradeDataKey(chainID))
}

// IterateAxelarProxyUpgradeData iterates over all axelar proxy upgrade data
func (k Keeper) IterateAxelarProxyUpgradeData(ctx sdk.Context, cb func(chainID uint64, upgradeData types.AxelarUpgradeData) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.AxelarProxyUpgradeDataPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		chainID := sdk.BigEndianToUint64(keyPair.Next(8))
		upgradeData := types.AxelarUpgradeData{}
		k.cdc.MustUnmarshal(iter.Value(), &upgradeData)
		if cb(chainID, upgradeData) {
			break
		}
	}
}

///////////////////////////
// Validator Cork counts //
///////////////////////////

/////////////////////
// Module Accounts //
/////////////////////

func (k Keeper) GetSenderAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

///////////////////////////
// ICS4Wrapper functions //
///////////////////////////

func (k Keeper) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (sequence uint64, err error) {
	if err := k.ValidateAxelarPacket(ctx, sourceChannel, data); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("ICS20 packet send was denied: %s", err.Error()))
		// based on the default implementation of SendPacket in ibc-go, we return 0 for the sequence on error conditions
		return 0, err
	}
	return k.Ics4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

func (k Keeper) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return k.Ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (k Keeper) GetAppVersion(ctx sdk.Context, portID string, channelID string) (string, bool) {
	return k.Ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}
