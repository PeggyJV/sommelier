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

	"github.com/peggyjv/sommelier/v9/x/poa/keeper"
	"github.com/peggyjv/sommelier/v9/x/poa/types"
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
	k := keeper.NewKeeper(cdc, storeKey, subspace, fake, nil, "cosmos1zkmrn5j2t9k3s7n9z2c5w0n9c0w8e8mq8w8mq8")
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

func TestWrapper_Slash_NoSnapshotSkips(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	cons := sdk.ConsAddress([]byte("auth-cons-aaaaaaaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.mapCons(cons, auth)

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	// NOTE: no snapshot at the infraction height.

	burned := w.Slash(ctx, cons, 999, 500, sdk.MustNewDecFromStr("0.05"))
	require.False(t, fake.slashCalled, "Slash should be refused when no snapshot exists")
	require.True(t, burned.IsZero())
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
