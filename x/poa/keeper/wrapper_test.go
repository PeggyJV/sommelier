package keeper_test

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v10/x/poa/keeper"
	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

// newWrapperTestKeeper builds a real PoA Keeper backed by an in-memory store
// and a hand-rolled fake staking keeper. Returns the keeper, context, fake,
// and the wrapper for direct use in tests.
func newWrapperTestKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *fakeStakingKeeper, keeper.WrappedStakingKeeper) {
	t.Helper()
	db := tmdb.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	cms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	paramsKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	cms.MountStoreWithDB(paramsKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTKey, storetypes.StoreTypeTransient, nil)

	require.NoError(t, cms.LoadLatestVersion())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	pk := paramskeeper.NewKeeper(cdc, codec.NewLegacyAmino(), paramsKey, paramsTKey)
	subspace := pk.Subspace(types.ModuleName).WithKeyTable(types.ParamKeyTable())

	fake := newFakeStaking()
	k := keeper.NewKeeper(cdc, storeKey, subspace, fake, nil, testGovAuthority)
	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())
	k.SetParams(ctx, types.DefaultParams())

	w := k.WrappedStakingKeeper()
	return k, ctx, fake, w
}

func TestWrapper_Validator_BoostsAuthorityTokens(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.addValidator(com, sdk.NewInt(3_000_000))
	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: ctx.BlockHeight(),
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "5.0"},
		},
	})

	vAuth := w.Validator(ctx, auth)
	require.Equal(t, sdk.NewInt(5_000_000), vAuth.GetTokens())

	vCom := w.Validator(ctx, com)
	require.Equal(t, sdk.NewInt(3_000_000), vCom.GetTokens())
}

func TestWrapper_IterateBondedValidatorsByPower_Boosted(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.addValidator(com, sdk.NewInt(3_000_000))
	fake.bondedOrder = []sdk.ValAddress{com, auth}

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: ctx.BlockHeight(),
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "5.0"},
		},
	})

	var seenTokens []string
	w.IterateBondedValidatorsByPower(ctx, func(_ int64, v stakingtypes.ValidatorI) bool {
		seenTokens = append(seenTokens, v.GetTokens().String())
		return false
	})
	require.Equal(t, []string{"3000000", "5000000"}, seenTokens)
}

func TestWrapper_GetBondedValidatorsByPower_Boosted(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.addValidator(com, sdk.NewInt(3_000_000))
	fake.bondedOrder = []sdk.ValAddress{com, auth}

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: ctx.BlockHeight(),
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "5.0"},
		},
	})

	out := w.GetBondedValidatorsByPower(ctx)
	require.Len(t, out, 2)
	// post-boost: auth is 5M, com is 3M, so auth must be first.
	require.Equal(t, auth.String(), out[0].OperatorAddress)
	require.Equal(t, sdk.NewInt(5_000_000), out[0].Tokens)
	require.Equal(t, sdk.NewInt(3_000_000), out[1].Tokens)
}

func TestWrapper_Slash_AuthorityNormalises(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	cons := sdk.ConsAddress([]byte("auth-cons-aaaaaaaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.mapCons(cons, auth)

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: 50,
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "5.0"},
		},
	})

	// Caller passes boosted power of 500. Wrapper must slash on raw 100.
	w.Slash(ctx, cons, 50, 500, sdk.MustNewDecFromStr("0.05"))
	require.True(t, fake.slashCalled)
	require.Equal(t, int64(100), fake.lastSlashPower)
}

func TestWrapper_Slash_CommunityPassesThrough(t *testing.T) {
	_, ctx, fake, w := newWrapperTestKeeper(t)

	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	cons := sdk.ConsAddress([]byte("com-cons-aaaaaaaaaaa"))
	fake.addValidator(com, sdk.NewInt(3_000_000))
	fake.mapCons(cons, com)

	w.Slash(ctx, cons, 100, 300, sdk.MustNewDecFromStr("0.05"))
	require.True(t, fake.slashCalled)
	require.Equal(t, int64(300), fake.lastSlashPower)
}

func TestWrapper_Slash_NoSnapshotPassesThrough(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	cons := sdk.ConsAddress([]byte("auth-cons-aaaaaaaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.mapCons(cons, auth)

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	// NOTE: no snapshot and no activation height set. With activation unset the
	// module is treated as not-yet-active, so the slash passes through against
	// raw power unchanged.
	w.Slash(ctx, cons, 999, 500, sdk.MustNewDecFromStr("0.05"))
	require.True(t, fake.slashCalled)
	require.Equal(t, int64(500), fake.lastSlashPower)
}

// Below the activation height, a missing snapshot is benign (pre-PoA, no boost
// was ever applied) — the slash passes through against raw power.
func TestWrapper_Slash_PreActivationNoSnapshotPassesThrough(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	cons := sdk.ConsAddress([]byte("auth-cons-aaaaaaaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.mapCons(cons, auth)

	k.SetActivationHeight(ctx, 100)
	// Infraction at height 50 < activation 100, no snapshot.
	w.Slash(ctx, cons, 50, 500, sdk.MustNewDecFromStr("0.05"))
	require.True(t, fake.slashCalled)
	require.Equal(t, int64(500), fake.lastSlashPower)
}

// At or above the activation height an empty snapshot is written every block,
// so a missing snapshot means corruption — the slash must be refused rather
// than risk over-slashing on (possibly boosted) power.
func TestWrapper_Slash_PostActivationMissingSnapshotRefused(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	cons := sdk.ConsAddress([]byte("auth-cons-aaaaaaaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.mapCons(cons, auth)

	k.SetActivationHeight(ctx, 100)
	// Infraction at height 200 >= activation 100, but no snapshot exists.
	burned := w.Slash(ctx, cons, 200, 500, sdk.MustNewDecFromStr("0.05"))
	require.False(t, fake.slashCalled, "post-activation slash with missing snapshot must be refused")
	require.True(t, burned.IsZero())
}

// TestWrapper_Slash_NormalisesByInfractionHeightSnapshot exercises the core
// fix: a validator boosted at the infraction height must be normalised even if
// it has since been removed from the current authority set.
func TestWrapper_Slash_NormalisesByInfractionHeightSnapshot(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	cons := sdk.ConsAddress([]byte("auth-cons-aaaaaaaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.mapCons(cons, auth)

	// Snapshot records a 5x boost at height 50, but the validator is NOT in the
	// current authority set (removed after the infraction). Normalisation must
	// still divide out the snapshot multiplier.
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: 50,
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "5.0"},
		},
	})

	w.Slash(ctx, cons, 50, 500, sdk.MustNewDecFromStr("0.05"))
	require.True(t, fake.slashCalled)
	require.Equal(t, int64(100), fake.lastSlashPower)
}

func TestWrapper_Slash_AuthorityNoBoostPassesThrough(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	cons := sdk.ConsAddress([]byte("auth-cons-aaaaaaaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.mapCons(cons, auth)

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	// Snapshot exists but multiplier is 1 (e.g., authority already above floor).
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: 50,
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "1.0"},
		},
	})

	w.Slash(ctx, cons, 50, 100, sdk.MustNewDecFromStr("0.05"))
	require.True(t, fake.slashCalled)
	require.Equal(t, int64(100), fake.lastSlashPower)
}

// Boosted reads must use the latest committed snapshot, not the current block's
// (not-yet-written) snapshot. During block N, the snapshot for N is written
// only in N's EndBlocker, so mid-block reads must fall back to N-1's snapshot —
// matching the LastValidatorPower the staking store still holds from N-1.
func TestWrapper_BoostUsesLatestCommittedSnapshot(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})

	// Snapshot committed at height 9; no snapshot yet for the current block 10.
	ctx = ctx.WithBlockHeight(10)
	k.SetMultiplierSnapshot(ctx.WithBlockHeight(9), types.MultiplierSnapshot{
		Height: 9,
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "5.0"},
		},
	})

	// Reading at height 10 must still see the 5x boost from height 9.
	vAuth := w.Validator(ctx, auth)
	require.Equal(t, sdk.NewInt(5_000_000), vAuth.GetTokens())
}
