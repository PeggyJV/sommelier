package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/x/reinvest/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper
	gravityKeeper types.GravityKeeper
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper, gravityKeeper types.GravityKeeper,
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
		gravityKeeper: gravityKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

///////////////////////
// MsgSubmitReinvest //
///////////////////////

// SetReinvestment sets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetReinvestment(ctx sdk.Context, val sdk.ValAddress, reinvestment types.Reinvestment) {
	bz := k.cdc.MustMarshal(&reinvestment)
	ctx.KVStore(k.storeKey).Set(types.GetReinvestmentForValidatorAddressKey(val, common.HexToAddress(reinvestment.Address)), bz)
}

// GetReinvestment gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetReinvestment(ctx sdk.Context, val sdk.ValAddress, cel common.Address) (types.Reinvestment, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetReinvestmentForValidatorAddressKey(val, cel))
	if len(bz) == 0 {
		return types.Reinvestment{}, false
	}

	var reinvestment types.Reinvestment
	k.cdc.MustUnmarshal(bz, &reinvestment)
	return reinvestment, true
}

// DeleteReinvestment deletes the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteReinvestment(ctx sdk.Context, val sdk.ValAddress, cel common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetReinvestmentForValidatorAddressKey(val, cel))
}

// HasReinvestmentForContract gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasReinvestmentForContract(ctx sdk.Context, val sdk.ValAddress, contract common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetReinvestmentForValidatorAddressKey(val, contract))
}

// HasReinvestment gets the existence of any commit for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasReinvestment(ctx sdk.Context, val sdk.ValAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetReinvestmentValidatorKeyPrefix(val))
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	return iter.Valid()
}

// IterateReinvestments iterates over all votes in the store
func (k Keeper) IterateReinvestments(ctx sdk.Context, handler func(val sdk.ValAddress, cel common.Address, reinvestment types.Reinvestment) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.ReinvestmentForAddressKeyPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.ReinvestmentForAddressKeyPrefix}))
		val := sdk.ValAddress(keyPair.Next(20))
		cel := common.BytesToAddress(keyPair.Bytes())

		var reinvestment types.Reinvestment
		k.cdc.MustUnmarshal(iter.Value(), &reinvestment)
		if handler(val, cel, reinvestment) {
			break
		}
	}
}

// IterateReinvestmentAddresses iterates over all addresses who have committed reinvests
func (k Keeper) IterateReinvestmentAddresses(ctx sdk.Context, handler func(addr common.Address) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.ReinvestmentForAddressKeyPrefix})
	defer iter.Close()

	seenAddresses := mapset.NewThreadUnsafeSet()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.ReinvestmentForAddressKeyPrefix}))
		keyPair.Next(20)
		address := common.BytesToAddress(keyPair.Bytes())

		// add seen address to set. if already in set, don't return to consumer
		if !seenAddresses.Add(address) {
			continue
		}

		if handler(address) {
			break
		}
	}
}

// IterateAddressReinvestments iterates over all reinvestments for an address
func (k Keeper) IterateAddressReinvestments(ctx sdk.Context, handler func(addr common.Address) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.ReinvestmentForAddressKeyPrefix})
	defer iter.Close()

	seenAddresses := mapset.NewThreadUnsafeSet()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.ReinvestmentForAddressKeyPrefix}))
		keyPair.Next(20)
		address := common.BytesToAddress(keyPair.Bytes())

		// add seen address to set. if already in set, don't return to consumer
		if !seenAddresses.Add(address) {
			continue
		}

		if handler(address) {
			break
		}
	}
}

//////////////////
// CommitPeriod //
//////////////////

// SetCommitPeriodStart sets the current vote period start height
func (k Keeper) SetCommitPeriodStart(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.CommitPeriodStartKey}, sdk.Uint64ToBigEndian(uint64(height)))
}

// GetCommitPeriodStart returns the vote period start height
func (k Keeper) GetCommitPeriodStart(ctx sdk.Context) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.CommitPeriodStartKey})
	if len(bz) == 0 {
		return 0, false
	}

	return int64(sdk.BigEndianToUint64(bz)), true
}

// HasCommitPeriodStart returns true if the vote period start has been set
func (k Keeper) HasCommitPeriodStart(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has([]byte{types.CommitPeriodStartKey})
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

// setParams sets the parameters in the store
func (k Keeper) setParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}


/////////////////////////
// Invalidation Nonces //
/////////////////////////

func (k Keeper) GetLatestInvalidationNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.LatestInvalidationNonceKey})
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLatestInvalidationNonce(ctx sdk.Context, invalidationNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LatestInvalidationNonceKey}, sdk.Uint64ToBigEndian(invalidationNonce))
}

func (k Keeper) IncrementInvalidationNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextNonce := k.GetLatestInvalidationNonce(ctx) + 1
	store.Set([]byte{types.LatestInvalidationNonceKey}, sdk.Uint64ToBigEndian(nextNonce))
	return nextNonce
}


///////////
// Votes //
///////////

func (k Keeper) GetWinningVotes(ctx sdk.Context, threshold sdk.Dec) (winningVotes []types.Reinvestment) {

	var reinvestments []types.Reinvestment
	var reinvestmentPowers []int64

	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)

	k.IterateReinvestments(ctx, func(val sdk.ValAddress, addr common.Address, reinvestment types.Reinvestment) (stop bool) {
		validator := k.stakingKeeper.Validator(ctx, val)
		validatorPower := validator.GetConsensusPower(k.stakingKeeper.PowerReduction(ctx))

		found := false
		for i, rv := range reinvestments {
			if rv.Equals(reinvestment) {
				reinvestmentPowers[i] += validatorPower

				found = true
				break
			}
		}

		if !found {
			reinvestments = append(reinvestments, reinvestment)
			reinvestmentPowers = append(reinvestmentPowers, validatorPower)
		}

		k.DeleteReinvestment(ctx, val, addr)

		return false
	})

	var winningReinvestments []types.Reinvestment

	for i, power := range reinvestmentPowers {
		quorumReached := sdk.NewDec(power).Quo(totalPower.ToDec()).GT(threshold)
		if quorumReached {
			winningReinvestments = append(winningReinvestments, reinvestments[i])
		}
	}

	return winningReinvestments
}
