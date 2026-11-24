package app

import (
	"bytes"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
	"github.com/stretchr/testify/require"
)

func newSafeModeSecurityApp(t *testing.T) (*SommelierApp, sdk.Context) {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{tmtypes.NewValidator(pubKey, 1)})
	accountKey := secp256k1.GenPrivKeyFromSecret([]byte("safe-mode-security-account"))
	account := authtypes.NewBaseAccount(
		accountKey.PubKey().Address().Bytes(), accountKey.PubKey(), 0, 0,
	)
	balance := banktypes.Balance{
		Address: account.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)),
	}

	app := NewSommelierApp(
		log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{},
		t.TempDir(), 0, MakeEncodingConfig(), EmptyAppOptions{},
	)
	genesis := genesisStateWithValSet(
		t, app, NewDefaultGenesisState(), valSet,
		[]authtypes.GenesisAccount{account}, balance,
	)
	stateBytes, err := tmjson.Marshal(genesis)
	require.NoError(t, err)
	app.InitChain(abci.RequestInitChain{
		ConsensusParams: DefaultConsensusParams,
		AppStateBytes:   stateBytes,
	})
	ctx := app.NewContext(false, tmproto.Header{Height: 1})
	app.PoaKeeper.SetSafeMode(ctx, true)
	return app, ctx
}

func TestEndBlockOrderMakesPoaPowerVisibleToGravity(t *testing.T) {
	app, _ := newSafeModeSecurityApp(t)

	index := func(moduleName string) int {
		for i, name := range app.mm.OrderEndBlockers {
			if name == moduleName {
				return i
			}
		}
		return -1
	}

	gravityIndex := index(gravitytypes.ModuleName)
	poaIndex := index(poatypes.ModuleName)
	require.NotEqual(t, -1, gravityIndex)
	require.NotEqual(t, -1, poaIndex)
	require.Less(t, poaIndex, gravityIndex,
		"gravity must observe PoA's current multiplier snapshot and validator powers")
}

type toggledSafeMode struct{ active bool }

func (s *toggledSafeMode) SafeModeActive(sdk.Context) bool { return s.active }

func TestGravitySafeModeModuleKeepsThawBlockFrozen(t *testing.T) {
	poa := &toggledSafeMode{active: true}
	module := &gravitySafeModeModule{poa: poa}
	ctx := sdk.Context{}

	// BeginBlock observes the freeze. PoA clears its persisted flag later in
	// the same block, but gravity EndBlock must still no-op because validators
	// had no unfrozen DeliverTx phase in which to submit confirmations.
	module.BeginBlock(ctx, abci.RequestBeginBlock{})
	poa.active = false
	require.NotPanics(t, func() {
		module.EndBlock(ctx, abci.RequestEndBlock{})
	})
}

func TestSafeModeMessageRouterRejectsCommunityPoolSpend(t *testing.T) {
	app, ctx := newSafeModeSecurityApp(t)
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	msg := &distrtypes.MsgCommunityPoolSpend{
		Authority: govAddr.String(),
		Recipient: sdk.AccAddress(bytes.Repeat([]byte{0x41}, 20)).String(),
		Amount:    sdk.NewCoins(),
	}

	handler := app.MsgServiceRouter().Handler(msg)
	require.NotNil(t, handler)
	_, err := handler(ctx, msg)
	require.ErrorContains(t, err, "safe mode active")

	app.PoaKeeper.SetSafeMode(ctx, false)
	_, err = handler(ctx, msg)
	require.NoError(t, err, "community-pool spends must resume after safe mode")
}

func TestSafeModeMessageRouterRejectsStakingParams(t *testing.T) {
	app, ctx := newSafeModeSecurityApp(t)
	msg := &stakingtypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    app.StakingKeeper.GetParams(ctx),
	}

	handler := app.MsgServiceRouter().Handler(msg)
	require.NotNil(t, handler)
	_, err := handler(ctx, msg)
	require.ErrorContains(t, err, "safe mode active")

	app.PoaKeeper.SetSafeMode(ctx, false)
	_, err = handler(ctx, msg)
	require.NoError(t, err, "staking parameter updates must resume after safe mode")
}

func TestSafeModeMessageRouterRejectsAuthzGrant(t *testing.T) {
	app, ctx := newSafeModeSecurityApp(t)
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	grantee := sdk.AccAddress(bytes.Repeat([]byte{0x42}, 20))
	msg, err := authz.NewMsgGrant(
		govAddr,
		grantee,
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&poatypes.MsgUpdateAuthoritySet{})),
		nil,
	)
	require.NoError(t, err)

	handler := app.MsgServiceRouter().Handler(msg)
	require.NotNil(t, handler)
	_, err = handler(ctx, msg)
	require.ErrorContains(t, err, "safe mode active")
	stored, _ := app.AuthzKeeper.GetAuthorization(
		ctx, grantee, govAddr, sdk.MsgTypeURL(&poatypes.MsgUpdateAuthoritySet{}),
	)
	require.Nil(t, stored)
}
